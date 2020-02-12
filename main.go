package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func main() {
	s := Settings{}

	flag.StringVar(&s.target, "target", "", "URL of an individual target you wish to scan.")
	flag.StringVar(&s.file, "filename", "", "File name of a list of targets. One per line.")
	flag.IntVar(&s.threads, "threads", 10, "Number of threads to use.")
	flag.BoolVar(&s.quiet, "quiet", false, "Option to suppress status messages.")
	flag.Parse()

	if s.quiet {
		log.SetOutput(ioutil.Discard)
	}

	if s.target == "" && s.file == "" {
		fmt.Printf("Either -target or -filename is required. See -h")
		os.Exit(0)
	}

	if s.file != "" {
		f, err := os.Open(s.file)
		if err != nil {
			log.Fatal("[main] Could not open the supplied filename for reading")
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			s.targetList = append(s.targetList, scanner.Text())
		}
	} else {
		s.targetList = append(s.targetList, s.target)
	}

	SeeMS(s)
}

func SeeMS(set Settings) {
	var targetQueue = make(chan Target)
	var targetSync sync.WaitGroup

	// We want to set test runners ready to pull from the channel ahead of pushing results into the channel
	for i := 0; i < set.threads; i++ {
		go testWorker(targetQueue, &targetSync)
	}

	// Push targets into the target queue. Channel is pulled off by test runners.
	for _, t := range set.targetList {
		targetSync.Add(1)
		targetQueue <- addTarget(t)
	}

	// We can close the queue now as all targets have been placed onto it.
	close(targetQueue)

	targetSync.Wait()
}

func addTarget(t string) Target {
	// this function will set up a unique UUID for each target.
	// each target has a series of tests associated with it.
	return Target{
		uuid.New(),
		t,
		GenerateTests(),
	}
}

func testWorker(targetQueue chan Target, targetSync *sync.WaitGroup) {
	// There are n testWorker functions, where n = set.threads.
	// this function will pull of a target and action it. the channel will close once all targets
	// have been read. Number of test runners is determined by set.threads.
	for target := range targetQueue {
		var webComms = make(chan WebResponse, len(target.Tests))

		for _, test := range target.Tests {
			// We can do all web requests independently and use a channel to read back when they are complete.
			go webWorker(fmt.Sprintf("%s%s", target.Hostname, test.Url), webComms, test.Id)
		}

		for i := 0; i < cap(webComms); i++ {
			r := <-webComms
			target.Tests[r.TestId].Response = r
		}

		close(webComms)

		scoreMarkup(target, targetSync)
	}
}

func scoreMarkup(target Target, targetSync *sync.WaitGroup) {
	// we now have a target with fully populated tests / associated errors.
	var scoreComms = make(chan Score, len(target.Tests))

	for _, test := range target.Tests {
		switch action := test.Action; action {
		case "substring":
			go substringWorker(scoreComms, test)
		case "regex":
			go regexWorker(scoreComms, test)
		case "hash":
			go hashWorker(scoreComms, test)
		case "header":
			go headerWorker(scoreComms, test)
		}
	}

	for i := 0; i < cap(scoreComms); i++ {
		s := <-scoreComms
		target.Tests[s.TestId].Score = s
	}

	close(scoreComms)

	scoreEvaluation(target, targetSync)
}

func scoreEvaluation(target Target, targetSync *sync.WaitGroup) {
	resMap := make(map[string]int)
	maxScore := 0
	result := ""

	// At this point we have a full target with scores and results on appropriate tests
	for i := range target.Tests {
		resMap[target.Tests[i].Cms] += target.Tests[i].Score.Value
	}

	for k, v := range resMap {
		if v > maxScore {
			maxScore = v
			result = k
		}
	}

	if maxScore > 0 {
		switch result {
		case "wordpress":
			ver := getWordpressVersion(target.Tests)
			fmt.Printf("%s (%s) @ %s\n", strings.Title(result), ver, target.Hostname)
			targetSync.Done()
		case "drupal":
			ver := getDrupalVersion(target.Tests)
			fmt.Printf("%s (%s) @ %s\n", strings.Title(result), ver, target.Hostname)
			targetSync.Done()
		case "joomla":
			ver := getJoomlaVersion(target.Tests)
			fmt.Printf("%s (%s) @ %s\n", strings.Title(result), ver, target.Hostname)
			targetSync.Done()
		case "sharepoint":
			ver := getSharepointVersion(target.Tests)
			fmt.Printf("%s (%s) @ %s\n", strings.Title(result), ver, target.Hostname)
			targetSync.Done()
		case "moodle":
			ver := "Unknown"
			fmt.Printf("%s (%s) @ %s\n", strings.Title(result), ver, target.Hostname)
		}
	} else {
		fmt.Printf("No recognised CMS @ %s\n", target.Hostname)
		targetSync.Done()
	}
}

func getSharepointVersion(sd []Test) string {
	spVersion := ""

	for _, s := range sd {
		if s.Id == 18 {
			spVersion = s.Response.Headers.Get("microsoftsharepointteamservices")
		}
	}

	if spVersion != "" {
		return spVersion
	} else {
		return "Unknown"
	}
}

func getWordpressVersion(sd []Test) string {
	wpFeedVersion := ""
	wpUpgradeVersion := ""

	for _, s := range sd {
		if s.Id == 1 {
			wpFeedVersion = s.Score.Regex
		} else if s.Id == 2 {
			wpUpgradeVersion = s.Score.Regex
		}
	}

	if wpFeedVersion != "" {
		return wpFeedVersion
	} else if wpUpgradeVersion != "" {
		return wpUpgradeVersion
	} else {
		return "Unknown"
	}
}

func getJoomlaVersion(sd []Test) string {
	joomlaVersion := ""

	for _, s := range sd {
		if s.Id == 3 {
			joomlaVersion = s.Score.Regex
		}
	}

	if joomlaVersion != "" {
		return joomlaVersion
	} else {
		return "Unknown"
	}
}

func getDrupalVersion(sd []Test) string {
	var verMin float64
	var verMax float64

	for _, s := range sd {
		versionSlice := []float64{}

		for _, k := range drupalVersion[s.Url] {
			if s.Score.Hash == k.Md5 {
				v, _ := strconv.ParseFloat(k.Id, 64)
				versionSlice = append(versionSlice, v)
			}
		}

		if len(versionSlice) > 0 {
			min := versionSlice[0]
			max := versionSlice[0]
			for _, v := range versionSlice {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
			}

			if min > verMin {
				verMin = min
			}

			if max > verMax {
				verMax = max
			}
		}
	}

	if verMin != 0 && verMax != 0 {
		return fmt.Sprintf("%v - %v", verMin, verMax)
	} else {
		return "Unknown"
	}
}

package main

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
)

func substringWorker(scoreComms chan Score, test Test) {
	var sd Score
	sd.Value = 0
	sd.TestId = test.Id

	if test.Response.Err == nil {
		if strings.Contains(test.Response.Body, test.Match) {
			sd.Value = 5
		}
	}

	scoreComms <- sd
}

func regexWorker(scoreComms chan Score, test Test) {
	var sd Score
	sd.Value = 0
	sd.TestId = test.Id

	if test.Response.Err == nil {
		r, _ := regexp.Compile(test.Match)
		matched := r.FindString(test.Response.Body)
		if matched != "" {
			matched = strings.TrimPrefix(matched, strings.Split(test.Match, "(.*)")[0])
			matched = strings.TrimSuffix(matched, strings.Split(test.Match, "(.*)")[1])
			sd.Regex = matched
			sd.Value = 5
		}
	}

	scoreComms <- sd
}

func hashWorker(scoreComms chan Score, test Test) {
	var sd Score
	sd.Value = 0
	sd.TestId = test.Id

	if test.Response.Err == nil {
		h := md5.Sum([]byte(test.Response.Body))
		sd.Hash = hex.EncodeToString(h[:])
	}

	scoreComms <- sd
}

func headerWorker(scoreComms chan Score, test Test) {
	var sd Score
	sd.Value = 0
	sd.TestId = test.Id

	if test.Response.Err == nil {
		if test.Response.Headers.Get(test.Match) != "" {
			sd.Value = 5
		}
	}

	scoreComms <- sd
}

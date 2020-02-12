package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

func webWorker(target string, webComms chan WebResponse, testId int) {
	// This function should make a web request on behalf of a target / test id.
	// It will either return the body result of that request, or the error associated with
	// attempting to retrieve it.

	w := WebResponse{
		testId,
		nil,
		"",
		nil,
	}

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(target)
	if err != nil {
		w.Err = err
		webComms <- w
		return
	}

	w.Headers = resp.Header

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Err = err
		webComms <- w
		return
	}

	w.Body = string(body)

	resp.Body.Close()

	webComms <- w
	return
}
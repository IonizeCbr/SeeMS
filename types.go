package main

import (
	"net/http"
)

type Settings struct {
	targetList []string
	target     string
	threads    int
	file       string
	quiet      bool
}

type Target struct {
	Hostname string
	Tests    []Test
}

type Test struct {
	Score    Score
	Id       int
	Cms      string
	Url      string
	Action   string
	Match    string
	Response WebResponse
}

type WebResponse struct {
	TestId  int
	Headers http.Header
	Body    string
	Err     error
}

type Version struct {
	Md5 string
	Id  string
}

type Score struct {
	TestId int
	Value  int
	Regex  string
	Hash   string
}

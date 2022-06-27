package models

import "net/http"

type RequestType struct {
	Id      int32       `json:"id"`
	Url     string      `json:"url"`
	Method  string      `json:"method"`
	Body    string      `json:"body"`
	Headers http.Header `json:"headers"`
}

type ResponseType struct {
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}

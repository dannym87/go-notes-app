package main

import (
	"github.com/google/jsonapi"
	"strconv"
)

const (
	InternalServerError = "Internal Server Error"
	MalformedJson       = "Malformed JSON"
	NotFound            = "Not Found"
	ValidationError     = "Validation Error"
)

type ErrorObject struct {
	Title  string
	Detail string
	Status int
}

func (e *ErrorObject) ToJSONAPIErrorObject() *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  e.Title,
		Detail: e.Detail,
		Status: strconv.Itoa(e.Status),
	}
}

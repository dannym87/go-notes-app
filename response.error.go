package main

const (
	InternalServerError = "Internal Server Error"
	MalformedJson       = "Malformed JSON"
	NotFound            = "Not Found"
	ValidationError     = "Validation Error"
	AuthenticationError = "Authentication Error"
	Unauthorised        = "Unauthorised"
)

type ErrorObject struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

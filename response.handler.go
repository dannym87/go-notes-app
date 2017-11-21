package main

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"github.com/google/jsonapi"
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"fmt"
)

const fieldErrMsg = "Field validation for '%s' failed on the '%s' tag [Key: '%s']"

type ResponseHandler interface {
	JSON(c *gin.Context, status int, model interface{})
	Errors(c *gin.Context, status int, errorObjects []*ErrorObject)
	Error(c *gin.Context, title string, status int, detail string)
	ValidationErrors(c *gin.Context, status int, errors error)
	InternalServerError(c *gin.Context)
	NotFound(c *gin.Context)
	MalformedJSON(c *gin.Context)
	NoRoute(c *gin.Context)
}

type JSONAPIResponse struct{}

func NewResponseHandler() ResponseHandler {
	return &JSONAPIResponse{}
}

func (*JSONAPIResponse) JSON(c *gin.Context, status int, model interface{}) {
	response := new(bytes.Buffer)
	jsonapi.MarshalPayload(response, model)
	c.Data(status, jsonapi.MediaType, response.Bytes())
}

func (*JSONAPIResponse) Errors(c *gin.Context, status int, errorObjects []*ErrorObject) {
	var errs []*jsonapi.ErrorObject

	for _, errorObject := range errorObjects {
		errs = append(errs, errorObject.ToJSONAPIErrorObject())
	}

	response := new(bytes.Buffer)
	jsonapi.MarshalErrors(response, errs)
	c.Data(status, jsonapi.MediaType, response.Bytes())
}

func (r *JSONAPIResponse) Error(c *gin.Context, title string, status int, detail string) {
	r.Errors(c, status, []*ErrorObject{
		&ErrorObject{title, detail, status},
	})
}

func (r *JSONAPIResponse) InternalServerError(c *gin.Context) {
	r.Error(c, InternalServerError, http.StatusInternalServerError, "Something went wrong")
}

func (r *JSONAPIResponse) NotFound(c *gin.Context) {
	r.Error(c, NotFound, http.StatusNotFound, "Resource does not exist")
}

func (r *JSONAPIResponse) MalformedJSON(c *gin.Context) {
	r.Error(c, MalformedJson, http.StatusBadRequest, "Request contains invalid JSON")
}

func (r *JSONAPIResponse) NoRoute(c *gin.Context) {
	r.Error(c, NotFound, http.StatusNotFound, "No route found")
}

func (r *JSONAPIResponse) ValidationErrors(c *gin.Context, status int, err error) {
	var errors []*ErrorObject

	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, &ErrorObject{
			ValidationError,
			fmt.Sprintf(fieldErrMsg, err.Field(), err.Tag(), err.Namespace()),
			status,
		})
	}

	r.Errors(c, status, errors)
}

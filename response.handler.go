package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"fmt"
)

const fieldErrMsg = "Field validation for '%s' failed on the '%s' tag [Key: '%s']"

type ResponseHandler interface {
	JSON(c *gin.Context, status int, model interface{})
	Errors(c *gin.Context, status int, errorObjects []*ErrorObject)
	Error(c *gin.Context, title string, status int, detail string)
	ValidationErrors(c *gin.Context, errors error)
	InternalServerError(c *gin.Context)
	NotFound(c *gin.Context)
	MalformedJSON(c *gin.Context)
	NoRoute(c *gin.Context)
	Unauthorised(c *gin.Context)
}

type APIResponseHandler struct{}

func NewResponseHandler() ResponseHandler {
	return &APIResponseHandler{}
}

func (*APIResponseHandler) JSON(c *gin.Context, status int, model interface{}) {
	c.JSON(status, gin.H{
		"data": model,
	})
}

func (*APIResponseHandler) Errors(c *gin.Context, status int, errorObjects []*ErrorObject) {
	c.JSON(status, gin.H{
		"errors": errorObjects,
	})
}

func (r *APIResponseHandler) Error(c *gin.Context, title string, status int, detail string) {
	r.Errors(c, status, []*ErrorObject{
		&ErrorObject{title, detail, status},
	})
}

func (r *APIResponseHandler) InternalServerError(c *gin.Context) {
	r.Error(c, InternalServerError, http.StatusInternalServerError, "Something went wrong")
}

func (r *APIResponseHandler) NotFound(c *gin.Context) {
	r.Error(c, NotFound, http.StatusNotFound, "Resource does not exist")
}

func (r *APIResponseHandler) MalformedJSON(c *gin.Context) {
	r.Error(c, MalformedJson, http.StatusBadRequest, "Request contains invalid JSON")
}

func (r *APIResponseHandler) NoRoute(c *gin.Context) {
	r.Error(c, NotFound, http.StatusNotFound, "No route found")
}

func (r *APIResponseHandler) ValidationErrors(c *gin.Context, err error) {
	var errors []*ErrorObject

	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, &ErrorObject{
			ValidationError,
			fmt.Sprintf(fieldErrMsg, err.Field(), err.Tag(), err.Namespace()),
			http.StatusUnprocessableEntity,
		})
	}

	r.Errors(c, http.StatusUnprocessableEntity, errors)
}

func (r *APIResponseHandler) Unauthorised(c *gin.Context) {
	r.Error(c, Unauthorised, http.StatusUnauthorized, "You don't have permission for this resource")
}

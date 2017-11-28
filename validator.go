package main

import (
	"gopkg.in/go-playground/validator.v9"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

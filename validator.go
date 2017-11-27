package main

import (
	"gopkg.in/go-playground/validator.v9"
)

func InitValidator() *validator.Validate {
	return validator.New()
}

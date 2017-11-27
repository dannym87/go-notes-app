package main

import (
	"gopkg.in/go-playground/validator.v9"
	"context"
	"fmt"
	"strings"
	"github.com/jinzhu/gorm"
)

func InitValidator() *validator.Validate {
	validator := validator.New()
	validator.RegisterValidationCtx("dbunique", DbUnique)

	return validator
}

func DbUnique(ctx context.Context, fl validator.FieldLevel) bool {
	db := ctx.Value("db").(*gorm.DB)
	parts := strings.Split(fl.Param(), ".")
	// todo check these exist before assignment
	table := parts[0]
	column := parts[1]

	var count int
	db.Table(table).Where(fmt.Sprintf("%s = ?", column), fl.Field().String()).Count(&count)

	if count > 0 {
		return false
	}

	return true
}

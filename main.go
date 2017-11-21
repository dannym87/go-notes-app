package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/go-playground/validator.v9"
)

var (
	validate        *validator.Validate
	responseHandler ResponseHandler
	db              *gorm.DB
)

func main() {
	var err error

	db, err = gorm.Open("sqlite3", "./notes.db")
	db.LogMode(true)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Note{}, &Tag{})

	validate = validator.New()
	responseHandler = NewResponseHandler()
	r := gin.Default()
	r.NoRoute(responseHandler.NoRoute)

	InitNotesHandler(r)

	r.Run()
}

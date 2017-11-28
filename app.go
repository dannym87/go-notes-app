package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"gopkg.in/go-playground/validator.v9"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type App struct {
	engine          *gin.Engine
	db              *gorm.DB
	responseHandler ResponseHandler
	validator       *validator.Validate
}

func InitApp() *App {
	db, err := gorm.Open("sqlite3", "./notes.db")
	db.LogMode(true)
	db.SingularTable(true)

	if err != nil {
		log.Fatal("Could not connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Note{}, &Tag{})

	validator := NewValidator()
	responseHandler := NewResponseHandler()
	r := gin.Default()
	r.NoRoute(responseHandler.NoRoute)

	app := &App{r, db, responseHandler, validator}

	InitHandlers(app)

	return app
}

func (app *App) Run() {
	app.Engine().Run()
}

func (app *App) Engine() *gin.Engine {
	return app.engine
}

func (app *App) Db() *gorm.DB {
	return app.db
}

func (app *App) ResponseHandler() ResponseHandler {
	return app.responseHandler
}

func (app *App) Validator() *validator.Validate {
	return app.validator
}

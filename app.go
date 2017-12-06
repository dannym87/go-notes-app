package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"gopkg.in/go-playground/validator.v9"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/gin-contrib/cors"
	"github.com/RangelReale/osin"
)

type App struct {
	engine          *gin.Engine
	db              *gorm.DB
	responseHandler ResponseHandler
	validator       *validator.Validate
	oauth2Server    *osin.Server
}

func InitApp() *App {
	db, err := gorm.Open("sqlite3", "./notes.db")
	db.LogMode(true)
	db.SingularTable(true)

	if err != nil {
		log.Fatal("Could not connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&Note{},
		&Tag{},
		&User{},
		&OAuth2Client{},
		&OAuth2RefreshToken{},
		&OAuth2AccessToken{},
	)

	validator := NewValidator()
	responseHandler := NewResponseHandler()
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(responseHandler.NoRoute)

	oauth2 := NewOAuth2Server(db)

	app := &App{r, db, responseHandler, validator, oauth2}

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

func (app *App) OAuth2Server() *osin.Server {
	return app.oauth2Server
}

package main

import "net/http"

func InitAPIDocHandler(app *App) {
	app.Engine().StaticFS("/api-doc", http.Dir("api-doc/dist"))
}

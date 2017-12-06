package main

func InitHandlers(app *App) {
	InitNotesHandler(app)
	InitTagsHandler(app)
	InitAuthHandler(app)
}

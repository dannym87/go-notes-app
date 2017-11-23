package main

func InitHandlers(app *App) {
	InitAPIDocHandler(app)
	InitNotesHandler(app)
	InitTagsHandler(app)
}

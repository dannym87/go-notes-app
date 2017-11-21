package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	app := InitApp()
	app.Run()
}

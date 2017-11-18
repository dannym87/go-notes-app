package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	db, err := gorm.Open("sqlite3", "./notes.db")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Note{})

	r := gin.Default()
	InitNotesHandler(r, db)

	r.Run()
}

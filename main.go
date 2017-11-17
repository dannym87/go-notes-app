package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

func main() {
	db, err := gorm.Open("sqlite3", "./notes.db")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Note{})

	n1 := &Note{1, "Note 1", time.Now()}
	n2 := &Note{2, "Note 2", time.Now()}

	db.Create(n1)
	db.Create(n2)

	r := gin.Default()
	InitNotesHandler(r, db)

	r.Run()
}

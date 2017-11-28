package main

import (
	"testing"
	"github.com/gin-gonic/gin"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/jinzhu/gorm"
	"fmt"
)

var app *App

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	initTestApp()
	os.Exit(m.Run())
}

func initTestApp() {
	db, err := gorm.Open("sqlite3", "./notes-test.db")
	db.SingularTable(true)

	if err != nil {
		panic("Cannot connect to test database")
	}

	r := gin.Default()
	populateDB(db)

	app = &App{r, db, NewResponseHandler(), NewValidator()}

	InitHandlers(app)
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func populateDB(db *gorm.DB) {
	db.DropTable(&Note{}, &Tag{})
	db.AutoMigrate(&Note{}, &Tag{})
	createTags(db, 11)
}

func createTags(db *gorm.DB, count int) {
	for i := 1; i < count + 1; i++ {
		db.Create(&Tag{Name: fmt.Sprintf("Tag %d", i)})
	}
}

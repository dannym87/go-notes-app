package main

import (
	"testing"
	"github.com/gin-gonic/gin"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/jinzhu/gorm"
	"fmt"
	"time"
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

	app = &App{
		r,
		db,
		NewResponseHandler(),
		NewRequestHandler(),
		NewValidator(),
		NewOAuth2Server(db),
	}

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
	dropSchema(db)
	createSchema(db)
	createTags(db, 11)
	createNotes(db, 11)
	createUsers(db)
	createOAuth2Clients(db)
	createOAuthAccessTokens(db)
	createOAuthRefreshTokens(db)
}

func dropSchema(db *gorm.DB) {
	db.DropTable(
		&Note{},
		&Tag{},
		&OAuth2Client{},
		&OAuth2AccessToken{},
		&OAuth2RefreshToken{},
		&User{},
		// many to many relationships
		"note_tags",
	)
}

func createSchema(db *gorm.DB) {
	db.AutoMigrate(
		&Note{},
		&Tag{},
		&OAuth2Client{},
		&OAuth2AccessToken{},
		&OAuth2RefreshToken{},
		&User{},
	)
}

func createTags(db *gorm.DB, count int) {
	for i := 1; i < count+1; i++ {
		db.Create(&Tag{Name: fmt.Sprintf("Tag %d", i)})
	}
}

func createNotes(db *gorm.DB, count int) {
	for i := 1; i < count+1; i++ {
		db.Create(&Note{
			Title:       fmt.Sprintf("Note %d", i),
			Text:        fmt.Sprintf("Note %d text...", i),
			CreatedById: 1,
		})
	}

	db.Create(&Note{Title: "User 2 Note 1", CreatedById: 2})
}

func createUsers(db *gorm.DB) {
	db.Create(&User{Email: "test@go-notes.com", Firstname: "Go", Lastname: "Notes", Password: "$2a$12$RFgkr30MuLQmPU5LNrVNZ.gev80MwIZRwTcTUfZBmf19vegxQq9CS"})
	db.Create(&User{Email: "test2@go-notes.com", Firstname: "Go2", Lastname: "Notes2", Password: "$2a$12$RFgkr30MuLQmPU5LNrVNZ.gev80MwIZRwTcTUfZBmf19vegxQq9CS"})
	db.Create(&User{Email: "test3@go-notes.com", Firstname: "Go3", Lastname: "Notes3", Password: "$2a$12$RFgkr30MuLQmPU5LNrVNZ.gev80MwIZRwTcTUfZBmf19vegxQq9CS"})
}

func createOAuth2Clients(db *gorm.DB) {
	notesClient := new(OAuth2Client)
	notesClient.RedirectURI = "http://www.go-notes.com/callback"
	notesClient.Secret = "$2a$12$UIvK0nN/7fvwT0PV/zaSc.vf.b7b0ItknYSjjNILapftiCbhxTDGm"
	notesClient.Extra = "User data..."

	db.Create(notesClient)
}

func createOAuthAccessTokens(db *gorm.DB) {
	// +1 hour
	expiry := time.Now().Local().Add(time.Hour)
	db.Create(&OAuth2AccessToken{AccessToken: "access-token", Scope: "email", Expires: expiry, ClientId: 1, UserId: 1})
	db.Create(&OAuth2AccessToken{AccessToken: "OWYzYjI3NDctY2ZmNy00ZjExLWExM2YtOTBlMmJhMWM2MDc1", Scope: "email", Expires: expiry, ClientId: 1, UserId: 3})
}

func createOAuthRefreshTokens(db *gorm.DB) {
	// +31 days
	expiry := time.Now().Local().Add(time.Hour * 24 * 31)
	db.Create(&OAuth2RefreshToken{RefreshToken: "refresh-token", Scope: "email", Expires: expiry, UserId: 1, ClientId: 1, AccessTokenId: 1})
	db.Create(&OAuth2RefreshToken{RefreshToken: "N2U3MTdmYjgtMzJhNi00MTE4LThjODMtYzQzM2RlZTBjZGFm", Scope: "email", Expires: expiry, UserId: 3, ClientId: 1, AccessTokenId: 2})
}

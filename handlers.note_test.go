package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"fmt"
	"bytes"
)

func TestNotesHandler_GetSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/notes/1", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestNotesHandler_GetNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/notes/0", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNotFound
	})
}

func TestNotesHandler_ListSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/notes", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got '%d'", w.Code)
			return false
		}

		data := struct {
			Notes []*Note `json:"data"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Error("Failed to unmarshal json")
			return false
		}

		if len(data.Notes) != 10 {
			t.Errorf("Expected 10 notes, got '%d'", len(data.Notes))
			return false
		}

		return true
	})
}

func TestNotesHandler_ListSuccessPage2(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/notes?page=2", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got '%d'", w.Code)
			return false
		}

		data := struct {
			Notes []*Note `json:"data"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Error("Failed to unmarshal json")
			return false
		}

		if len(data.Notes) != 1 {
			t.Errorf("Expected 1 note, got '%d'", len(data.Notes))
			return false
		}

		return true
	})
}

func TestNotesHandler_DeleteSuccess(t *testing.T) {
	note := &Note{Title: "Note X", Text: "Note X text..."}
	app.Db().Create(&note)
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/notes/%d", note.ID), nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNoContent
	})
}

func TestNotesHandler_DeleteNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "/v1/notes/0", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNotFound
	})
}

func TestNotesHandler_CreateSuccess(t *testing.T) {
	data, _ := json.Marshal(Note{Title: "Note X", Text: "Note X text..."})
	req, _ := http.NewRequest(http.MethodPost, "/v1/notes", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code 201, got '%d'", w.Code)
			return false
		}

		app.Db().Where("title = ?", "Note X").Delete(Note{})

		return true
	})
}

func TestNotesHandler_CreateSuccessWithTags(t *testing.T) {
	data, _ := json.Marshal(Note{
		Title: "Note X",
		Text:  "Note X text...",
		Tags: []*Tag{
			&Tag{Name: "Tag 1"},
			&Tag{Name: "Tag 2"},
			&Tag{Name: "New Tag"},
		},
	})
	req, _ := http.NewRequest(http.MethodPost, "/v1/notes", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code 201, got '%d'", w.Code)
			return false
		}

		data := struct {
			Note *Note `json:"data"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Error("Failed to unmarshal json")
			return false
		}

		if data.Note.Title != "Note X" {
			t.Errorf("Expected note title 'Note X', got '%s'", data.Note.Title)
			return false
		}

		if data.Note.Text != "Note X text..." {
			t.Errorf("Expected note text 'Note X text...', got '%s'", data.Note.Text)
			return false
		}

		if len(data.Note.Tags) != 3 {
			t.Errorf("Expected 3 tags, got '%d'", len(data.Note.Tags))
			return false
		}

		if data.Note.Tags[0].ID != 1 && data.Note.Tags[0].Name != "Tag 1" {
			t.Errorf("Exepected tag[1, Tag 1], got tag[%d, %s]", data.Note.Tags[0].ID, data.Note.Tags[0].Name)
			return false
		}

		if data.Note.Tags[1].ID != 1 && data.Note.Tags[1].Name != "Tag 2" {
			t.Errorf("Exepected tag[2, Tag 2], got tag[%d, %s]", data.Note.Tags[1].ID, data.Note.Tags[1].Name)
			return false
		}

		if data.Note.Tags[2].Name != "New Tag" {
			t.Errorf("Exepected tag[New Tag], got tag[%s]", data.Note.Tags[2].Name)
			return false
		}

		app.Db().Delete(data.Note)
		app.Db().Where("name = ?", "New Tag").Delete(Tag{})

		return true
	})
}

func TestNotesHandler_CreateValidationErrors(t *testing.T) {
	data, _ := json.Marshal(Note{})
	req, _ := http.NewRequest(http.MethodPost, "/v1/notes", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code 422, got '%d'", w.Code)
			return false
		}

		data := struct {
			Errors []*ErrorObject `json:"errors"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Error("Failed to unmarshal json")
			return false
		}

		if len(data.Errors) != 1 {
			t.Errorf("Expected 1 error, got '%d'", len(data.Errors))
			return false
		}

		expectedDetail := "Field validation for 'Title' failed on the 'required' tag [Key: 'Note.Title']"
		if data.Errors[0].Detail != expectedDetail {
			t.Errorf("Expected '%s', got '%s'", expectedDetail, data.Errors[0].Detail)
			return false
		}

		return true
	})
}

func TestNotesHandler_UpdateSuccess(t *testing.T) {
	data, _ := json.Marshal(Note{Title: "Go"})
	req, _ := http.NewRequest(http.MethodPatch, "/v1/notes/1", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got '%d'", w.Code)
			return false
		}

		// reset note to default state
		note := new(Note)
		app.Db().First(note, 1)
		app.Db().Model(note).UpdateColumns(Note{Title: "Note 1"})

		return true
	})
}

func TestNotesHandler_UpdateValidationErrors(t *testing.T) {
	data, _ := json.Marshal(Note{})
	req, _ := http.NewRequest(http.MethodPatch, "/v1/notes/1", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code 422, got '%d'", w.Code)
			return false
		}

		data := struct {
			Errors []*ErrorObject `json:"errors"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Error("Failed to unmarshal json")
			return false
		}

		if len(data.Errors) != 1 {
			t.Errorf("Expected 1 error, got '%d'", len(data.Errors))
			return false
		}

		expectedDetail := "Field validation for 'Title' failed on the 'required' tag [Key: 'Note.Title']"
		if data.Errors[0].Detail != expectedDetail {
			t.Errorf("Expected '%s', got '%s'", expectedDetail, data.Errors[0].Detail)
			return false
		}

		return true
	})
}

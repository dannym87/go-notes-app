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

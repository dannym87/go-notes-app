package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"fmt"
	"bytes"
)

func TestTagsHandler_GetSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/tags/1", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestTagsHandler_GetNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/tags/0", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNotFound
	})
}

func TestTagsHandler_ListSuccess(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/tags", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got '%d'", w.Code)
			return false
		}

		data := struct {
			Tags []*Tag `json:"data"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Error("Failed to unmarshal json")
			return false
		}

		if len(data.Tags) != 10 {
			t.Errorf("Expected 10 tags, got '%d'", len(data.Tags))
			return false
		}

		return true
	})
}

func TestTagsHandler_ListSuccessPage2(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/v1/tags?page=2", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got '%d'", w.Code)
			return false
		}

		data := struct {
			Tags []*Tag `json:"data"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Error("Failed to unmarshal json")
			return false
		}

		if len(data.Tags) != 1 {
			t.Errorf("Expected 1 tag, got '%d'", len(data.Tags))
			return false
		}

		return true
	})
}

func TestTagsHandler_DeleteSuccess(t *testing.T) {
	tag := &Tag{Name: "Go"}
	app.Db().Create(&tag)
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/tags/%d", tag.ID), nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNoContent
	})
}

func TestTagsHandler_DeleteNotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "/v1/tags/0", nil)

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNotFound
	})
}

func TestTagsHandler_CreateSuccess(t *testing.T) {
	data, _ := json.Marshal(Tag{Name: "Go"})
	req, _ := http.NewRequest(http.MethodPost, "/v1/tags", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code 201, got '%d'", w.Code)
			return false
		}

		app.Db().Where("name = ?", "Go").Delete(Tag{})

		return true
	})
}

func TestTagsHandler_CreateTagAlreadyExists(t *testing.T) {
	data, _ := json.Marshal(Tag{Name: "Tag 1"})
	req, _ := http.NewRequest(http.MethodPost, "/v1/tags", bytes.NewBuffer(data))

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

		expectedDetail := "Tag 'Tag 1' already exists"
		if data.Errors[0].Detail != expectedDetail {
			t.Errorf("Expected '%s', got '%s'", expectedDetail, data.Errors[0].Detail)
			return false
		}

		return true
	})
}

func TestTagsHandler_CreateValidationErrors(t *testing.T) {
	data, _ := json.Marshal(Tag{})
	req, _ := http.NewRequest(http.MethodPost, "/v1/tags", bytes.NewBuffer(data))

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

		expectedDetail := "Field validation for 'Name' failed on the 'required' tag [Key: 'Tag.Name']"
		if data.Errors[0].Detail != expectedDetail {
			t.Errorf("Expected '%s', got '%s'", expectedDetail, data.Errors[0].Detail)
			return false
		}

		return true
	})
}

func TestTagsHandler_UpdateSuccess(t *testing.T) {
	data, _ := json.Marshal(Tag{Name: "Go"})
	req, _ := http.NewRequest(http.MethodPatch, "/v1/tags/1", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got '%d'", w.Code)
			return false
		}

		// reset tag to default state
		tag := new(Tag)
		app.Db().First(tag, 1)
		app.Db().Model(tag).UpdateColumns(Tag{Name: "Tag 1"})

		return true
	})
}

func TestTagsHandler_UpdateValidationErrors(t *testing.T) {
	data, _ := json.Marshal(Tag{})
	req, _ := http.NewRequest(http.MethodPatch, "/v1/tags/1", bytes.NewBuffer(data))

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

		expectedDetail := "Field validation for 'Name' failed on the 'required' tag [Key: 'Tag.Name']"
		if data.Errors[0].Detail != expectedDetail {
			t.Errorf("Expected '%s', got '%s'", expectedDetail, data.Errors[0].Detail)
			return false
		}

		return true
	})
}

func TestTagsHandler_UpdateTagAlreadyExists(t *testing.T) {
	data, _ := json.Marshal(Tag{Name: "Tag 2"})
	req, _ := http.NewRequest(http.MethodPatch, "/v1/tags/1", bytes.NewBuffer(data))

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

		expectedDetail := "Tag 'Tag 2' already exists"
		if data.Errors[0].Detail != expectedDetail {
			t.Errorf("Expected '%s', got '%s'", expectedDetail, data.Errors[0].Detail)
			return false
		}

		return true
	})
}

func TestTagsHandler_UpdateNameCanRemainTheSame(t *testing.T) {
	data, _ := json.Marshal(Tag{Name: "Tag 1"})
	req, _ := http.NewRequest(http.MethodPatch, "/v1/tags/1", bytes.NewBuffer(data))

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got '%d'", w.Code)
			return false
		}

		return true
	})
}

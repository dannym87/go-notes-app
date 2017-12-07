package main

import (
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	"bytes"
	"github.com/RangelReale/osin"
	"net/url"
)

func TestAuthHandler_TokenPasswordSuccess(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("username", "test2@go-notes.com")
	params.Add("password", "password")
	params.Add("client_id", "1")
	params.Add("client_secret", "secret")
	params.Add("scope", "email")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code '201', got '%d'", w.Code)
			return false
		}

		token := struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			Scope        string `json:"scope"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &token); err != nil {
			t.Errorf("Could not unmarshal json")
			return false
		}

		if token.AccessToken == "" {
			t.Errorf("Access token is empty")
			return false
		}

		if token.RefreshToken == "" {
			t.Errorf("Refresh token is empty")
			return false
		}

		if token.Scope != "email" {
			t.Errorf("Expected scope 'email', got '%s'", token.Scope)
			return false
		}

		if token.ExpiresIn != 3600 {
			t.Errorf("Expected token to expire in '3600', got '%d'", token.ExpiresIn)
			return false
		}

		if token.TokenType != "Bearer" {
			t.Errorf("Expected token type 'Bearer', got '%s'", token.TokenType)
			return false
		}

		return true
	})
}

func TestAuthHandler_TokenPasswordInvalidGrantType(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "invalid_grant")
	params.Add("username", "test2@go-notes.com")
	params.Add("password", "password")
	params.Add("client_id", "1")
	params.Add("client_secret", "secret")
	params.Add("scope", "email")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code '400', got '%d'", w.Code)
			return false
		}

		data := struct {
			Errors []*ErrorObject `json:"errors"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Errorf("Unable to unmarshal json: '%s'", err.Error())
			return false
		}

		if len(data.Errors) != 1 {
			t.Errorf("Expected '1' error, got '%d'", len(data.Errors))
			return false
		}

		if status := data.Errors[0].Status; status != http.StatusBadRequest {
			t.Errorf("Expected status '400', got '%d'", status)
			return false
		}

		if title := data.Errors[0].Title; title != osin.E_UNSUPPORTED_GRANT_TYPE {
			t.Errorf("Expected status '%s', got '%s'", osin.E_UNSUPPORTED_GRANT_TYPE, title)
			return false
		}

		return true
	})
}

func TestAuthHandler_TokenPasswordInvalidUserUsername(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("username", "invalid_username")
	params.Add("password", "password")
	params.Add("client_id", "1")
	params.Add("client_secret", "secret")
	params.Add("scope", "email")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code '400', got '%d'", w.Code)
			return false
		}

		data := struct {
			Errors []*ErrorObject `json:"errors"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Errorf("Unable to unmarshal json: '%s'", err.Error())
			return false
		}

		if len(data.Errors) != 1 {
			t.Errorf("Expected '1' error, got '%d'", len(data.Errors))
			return false
		}

		if status := data.Errors[0].Status; status != http.StatusBadRequest {
			t.Errorf("Expected status '400', got '%d'", status)
			return false
		}

		if title := data.Errors[0].Title; title != AuthenticationError {
			t.Errorf("Expected status '%s', got '%s'", AuthenticationError, title)
			return false
		}

		return true
	})
}

func TestAuthHandler_TokenPasswordInvalidUserPassword(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("username", "test2@go-notes.com")
	params.Add("password", "invalid_password")
	params.Add("client_id", "1")
	params.Add("client_secret", "secret")
	params.Add("scope", "email")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code '400', got '%d'", w.Code)
			return false
		}

		data := struct {
			Errors []*ErrorObject `json:"errors"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Errorf("Unable to unmarshal json: '%s'", err.Error())
			return false
		}

		if len(data.Errors) != 1 {
			t.Errorf("Expected '1' error, got '%d'", len(data.Errors))
			return false
		}

		if status := data.Errors[0].Status; status != http.StatusBadRequest {
			t.Errorf("Expected status '400', got '%d'", status)
			return false
		}

		if title := data.Errors[0].Title; title != AuthenticationError {
			t.Errorf("Expected status '%s', got '%s'", AuthenticationError, title)
			return false
		}

		return true
	})
}

func TestAuthHandler_TokenPasswordInvalidClientCredentials(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("username", "test2@go-notes.com")
	params.Add("password", "invalid_password")
	params.Add("client_id", "1")
	params.Add("client_secret", "invalid_secret")
	params.Add("scope", "email")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code '400', got '%d'", w.Code)
			return false
		}

		data := struct {
			Errors []*ErrorObject `json:"errors"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Errorf("Unable to unmarshal json: '%s'", err.Error())
			return false
		}

		if len(data.Errors) != 1 {
			t.Errorf("Expected '1' error, got '%d'", len(data.Errors))
			return false
		}

		if status := data.Errors[0].Status; status != http.StatusBadRequest {
			t.Errorf("Expected status '400', got '%d'", status)
			return false
		}

		if title := data.Errors[0].Title; title != osin.E_UNAUTHORIZED_CLIENT {
			t.Errorf("Expected status '%s', got '%s'", osin.E_UNAUTHORIZED_CLIENT, title)
			return false
		}

		return true
	})
}

func TestAuthHandler_TokenPasswordMalformedClientCredentials(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("username", "test2@go-notes.com")
	params.Add("password", "invalid_password")
	params.Add("client_id", "")
	params.Add("client_secret", "")
	params.Add("scope", "email")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code '400', got '%d'", w.Code)
			return false
		}

		data := struct {
			Errors []*ErrorObject `json:"errors"`
		}{}

		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			t.Errorf("Unable to unmarshal json: '%s'", err.Error())
			return false
		}

		if len(data.Errors) != 1 {
			t.Errorf("Expected '1' error, got '%d'", len(data.Errors))
			return false
		}

		if status := data.Errors[0].Status; status != http.StatusBadRequest {
			t.Errorf("Expected status '400', got '%d'", status)
			return false
		}

		if title := data.Errors[0].Title; title != osin.E_INVALID_REQUEST {
			t.Errorf("Expected status '%s', got '%s'", osin.E_INVALID_REQUEST, title)
			return false
		}

		return true
	})
}

func TestAuthHandler_TokenRefreshTokenSuccess(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", "N2U3MTdmYjgtMzJhNi00MTE4LThjODMtYzQzM2RlZTBjZGFm")
	params.Add("client_id", "1")
	params.Add("client_secret", "secret")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code '201', got '%d'", w.Code)
			return false
		}

		return true
	})
}

func TestAuthHandler_TokenInvalidRefreshToken(t *testing.T) {
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", "invalid_token")
	params.Add("client_id", "1")
	params.Add("client_secret", "secret")
	req, _ := http.NewRequest(http.MethodPost, "/token", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	testHTTPResponse(t, app.Engine(), req, func(w *httptest.ResponseRecorder) bool {
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code '400', got '%d'", w.Code)
			return false
		}

		return true
	})
}

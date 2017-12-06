package main

import (
	"testing"
	"github.com/RangelReale/osin"
	"time"
	"github.com/satori/go.uuid"
)

func TestGORMStorage_GetClient(t *testing.T) {
	s := app.OAuth2Server().Storage
	c, err := s.GetClient("1")

	if err != nil {
		t.Errorf("Unable to get oauth c: '%s'", err.Error())
		return
	}

	client, clientOk := c.(*OAuth2Client)
	if !clientOk {
		t.Errorf("Could not assert type *OAuth2Client")
		return
	}

	if client.GetId() != "1" {
		t.Errorf("Expected c ID '1', got '%s'", client.GetId())
		return
	}

	if !client.ClientSecretMatches("secret") {
		t.Errorf("Expected c secret 'Secret', got '%s'", client.GetSecret())
		return
	}

	if client.GetRedirectUri() != "http://www.go-notes.com/callback" {
		t.Errorf("Expected c redirect uri 'http://www.go-notes.com/callback', got '%s'", client.GetRedirectUri())
		return
	}

	if client.GetUserData() != "User data..." {
		t.Errorf("Expected c user data 'User data...', got '%s'", client.GetUserData())
		return
	}
}

func TestGORMStorage_SaveAccess(t *testing.T) {
	u := new(User)
	app.Db().First(u, 1)
	s := app.OAuth2Server().Storage
	ad := &osin.AccessData{
		Client: &osin.DefaultClient{
			Id: "1",
		},
		RedirectUri:  "http://www.go-notes.com/callback",
		UserData:     u,
		AccessToken:  uuid.NewV4().String(),
		ExpiresIn:    3600,
		Scope:        "email",
		CreatedAt:    time.Now(),
		RefreshToken: uuid.NewV4().String(),
	}

	if err := s.SaveAccess(ad); err != nil {
		t.Errorf("Unable to save access token: '%s'", err.Error())
		return
	}
}

func TestGORMStorage_LoadAccess(t *testing.T) {
	s := app.OAuth2Server().Storage
	ad, err := s.LoadAccess("access-token")

	if err != nil {
		t.Errorf("Could not load access token: '%s'", err.Error())
		return
	}

	if ad.Client.GetId() != "1" {
		t.Errorf("Expected client id '1', got '%s'", ad.Client.GetId())
		return
	}

	if ad.AccessToken != "access-token" {
		t.Errorf("Expected access token 'access-token', got '%s'", ad.AccessToken)
		return
	}

	if ad.Scope != "email" {
		t.Errorf("Expeted scope 'email', got '%s'", ad.Scope)
		return
	}

	if ad.RedirectUri != "http://www.go-notes.com/callback" {
		t.Errorf("Expected redirect uri 'http://www.go-notes.com/callback', got '%s'", ad.RedirectUri)
		return
	}

	if _, ok := ad.UserData.(*User); !ok {
		t.Error("User is not an instance of 'User'")
		return
	}

	if ad.UserData.(*User).ID != 1 {
		t.Errorf("Expected user id '1', got '%d'", ad.UserData.(*User).ID)
		return
	}
}

func TestGORMStorage_RemoveAccess(t *testing.T) {
	token := &OAuth2AccessToken{
		AccessToken: uuid.NewV4().String(),
		UserId:      1,
		ClientId:    1,
		Scope:       "email",
		Expires:     time.Now(),
	}

	if err := app.Db().Create(token).Error; err != nil {
		t.Errorf("Could not create access token: '%s'", err.Error)
		return
	}

	s := app.OAuth2Server().Storage
	if err := s.RemoveAccess(token.AccessToken); err != nil {
		t.Errorf("Could not delete access token: '%s'", err.Error())
		return
	}
}

func TestGORMStorage_RemoveRefresh(t *testing.T) {
	token := &OAuth2RefreshToken{
		RefreshToken: uuid.NewV4().String(),
		Expires:      time.Now(),
		ClientId:     1,
		UserId:       1,
		Scope:        "email",
	}

	if err := app.Db().Create(token).Error; err != nil {
		t.Errorf("Could not create refresh token: '%s'", err.Error())
		return
	}

	s := app.OAuth2Server().Storage
	if err := s.RemoveRefresh(token.RefreshToken); err != nil {
		t.Errorf("Could not delete refresh token: '%s'", err.Error())
		return
	}
}

func TestGORMStorage_LoadRefresh(t *testing.T) {
	s := app.OAuth2Server().Storage
	ad, err := s.LoadRefresh("refresh-token")

	if err != nil {
		t.Errorf("Could not load refresh token: '%s'", err.Error())
		return
	}

	if ad.Client.GetId() != "1" {
		t.Errorf("Expected client id '1', got '%s'", ad.Client.GetId())
		return
	}

	if ad.AccessToken != "access-token" {
		t.Errorf("Expected access token 'access-token', got '%s'", ad.AccessToken)
		return
	}

	if ad.RefreshToken != "refresh-token" {
		t.Errorf("Expected refresh token 'refresh-token', got '%s'", err.Error())
		return
	}

	if ad.Scope != "email" {
		t.Errorf("Expeted scope 'email', got '%s'", ad.Scope)
		return
	}

	if _, ok := ad.UserData.(*User); !ok {
		t.Error("User is not an instance of 'User'")
		return
	}

	if ad.UserData.(*User).ID != 1 {
		t.Errorf("Expected user id '1', got '%d'", ad.UserData.(*User).ID)
		return
	}
}

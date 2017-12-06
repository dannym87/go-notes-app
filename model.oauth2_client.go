package main

import (
	"strconv"
	"golang.org/x/crypto/bcrypt"
)

type OAuth2Client struct {
	BaseModel
	Secret      string `json:"-"`
	Extra       string `json:"extra"`
	RedirectURI string `json:"redirect_uri"`
}

func (c *OAuth2Client) GetId() string {
	return strconv.Itoa(int(c.ID))
}

func (c *OAuth2Client) GetSecret() string {
	return c.Secret
}

func (c *OAuth2Client) GetRedirectUri() string {
	return c.RedirectURI
}

func (c *OAuth2Client) GetUserData() interface{} {
	return c.Extra
}

func (*OAuth2Client) TableName() string {
	return "oauth2_client"
}

func (c *OAuth2Client) ClientSecretMatches(secret string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(c.Secret), []byte(secret)); err != nil {
		return false
	}

	return true
}

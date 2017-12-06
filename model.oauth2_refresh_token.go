package main

import "time"

type OAuth2RefreshToken struct {
	BaseModel
	RefreshToken  string             `json:"refresh_token" gorm:"unique_index"`
	AccessToken   *OAuth2AccessToken `json:"access_token" gorm:"ForeignKey:AccessTokenId"`
	AccessTokenId uint               `json:"-"`
	Client        *OAuth2Client      `json:"client" gorm:"ForeignKey:ClientId"`
	ClientId      uint               `json:"-"`
	User          *User              `json:"user" gorm:"ForeignKey:UserId"`
	UserId        uint               `json:"-"`
	Expires       time.Time          `json:"expires"`
	Scope         string             `json:"scope"`
}

func (*OAuth2RefreshToken) TableName() string {
	return "oauth2_refresh_token"
}

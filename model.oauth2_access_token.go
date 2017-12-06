package main

import "time"

type OAuth2AccessToken struct {
	BaseModel
	AccessToken string        `json:"access_token" gorm:"unique_index"`
	Client      *OAuth2Client `json:"client" gorm:"ForeignKey:ClientId"`
	ClientId    uint          `json:"-"`
	User        *User         `json:"user" gorm:"ForeignKey:UserId"`
	UserId      uint          `json:"-"`
	Expires     time.Time     `json:"expires"`
	Scope       string        `json:"scope"`
}

func (*OAuth2AccessToken) TableName() string {
	return "oauth2_access_token"
}

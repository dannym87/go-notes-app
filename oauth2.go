package main

import (
	"github.com/RangelReale/osin"
	"github.com/jinzhu/gorm"
	"strconv"
	"errors"
	"time"
	"net/http"
)

type GORMStorage struct {
	db *gorm.DB
}

func NewOAuth2Server(db *gorm.DB) *osin.Server {
	conf := osin.NewServerConfig()
	conf.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD}
	conf.AllowGetAccessRequest = true
	conf.ErrorStatusCode = http.StatusBadRequest
	conf.AccessExpiration = 3600

	return osin.NewServer(conf, &GORMStorage{db})
}

func (s *GORMStorage) Clone() osin.Storage {
	return s
}

func (s *GORMStorage) Close() {

}

func (s *GORMStorage) GetClient(id string) (osin.Client, error) {
	client := new(OAuth2Client)

	if err := s.db.First(client, id).Error; err != nil {
		return nil, osin.ErrNotFound
	}

	return client, nil
}

func (s *GORMStorage) SaveAuthorize(*osin.AuthorizeData) error {
	return errors.New("Not implemented")
}

func (s *GORMStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	return nil, errors.New("Not implemented")
}

func (s *GORMStorage) RemoveAuthorize(code string) error {
	return errors.New("Not implemented")
}

func (s *GORMStorage) SaveAccess(t *osin.AccessData) error {
	c, err := s.GetClient(t.Client.GetId())
	if err != nil {
		return err
	}

	client, clientOk := c.(*OAuth2Client)
	if !clientOk {
		return errors.New("Could not assert type *OAuth2Client")
	}

	user, userOk := t.UserData.(*User)
	if !userOk {
		return errors.New("Could not assert type User")
	}

	prev := new(OAuth2AccessToken)
	s.db.Where("client_id = ? AND user_id = ?", client.GetId(), user.ID).Find(prev)

	if prev.AccessToken != "" {
		t.AccessData, _ = s.LoadAccess(prev.AccessToken)
	}

	token := &OAuth2AccessToken{
		AccessToken: t.AccessToken,
		Client:      client,
		ClientId:    client.ID,
		Expires:     t.ExpireAt(),
		Scope:       t.Scope,
		User:        user,
		UserId:      user.ID,
	}

	tx := s.db.Begin()

	if err := tx.Set("gorm:save_associations", false).Create(token).Error; err != nil {
		tx.Rollback()
		return err
	}

	refreshToken := &OAuth2RefreshToken{
		AccessTokenId: token.ID,
		RefreshToken:  t.RefreshToken,
		Client:        client,
		ClientId:      client.ID,
		Expires:       time.Now().AddDate(0, 0, 31),
		Scope:         t.Scope,
		User:          t.UserData.(*User),
		UserId:        t.UserData.(*User).ID,
	}

	if err := tx.Set("gorm:save_associations", false).Create(refreshToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *GORMStorage) LoadAccess(token string) (*osin.AccessData, error) {
	accessToken := new(OAuth2AccessToken)
	if err := s.db.Where("access_token = ?", token).Preload("Client").Preload("User").Find(accessToken).Error; err != nil {
		return nil, osin.ErrNotFound
	}

	refreshToken := new(OAuth2RefreshToken)
	if err := s.db.Where("access_token_id = ?", accessToken.ID).Find(refreshToken).Error; err != nil {
		return nil, err
	}

	t := &osin.AccessData{
		Client: &osin.DefaultClient{
			Id:          strconv.Itoa(int(accessToken.Client.ID)),
			RedirectUri: accessToken.Client.RedirectURI,
			Secret:      accessToken.Client.Secret,
			UserData:    accessToken.User.ID,
		},
		UserData:     accessToken.User,
		RedirectUri:  accessToken.Client.RedirectURI,
		CreatedAt:    accessToken.CreatedAt,
		ExpiresIn:    int32(accessToken.Expires.Sub(time.Now()).Seconds()),
		AccessToken:  accessToken.AccessToken,
		Scope:        accessToken.Scope,
		RefreshToken: refreshToken.RefreshToken,
	}

	return t, nil
}

func (s *GORMStorage) RemoveAccess(token string) error {
	if err := s.db.Where("access_token = ?", token).Delete(&OAuth2AccessToken{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *GORMStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	refreshToken := new(OAuth2RefreshToken)
	if err := s.db.Where("refresh_token = ?", token).Preload("AccessToken").Find(refreshToken).Error; err != nil {
		return nil, osin.ErrNotFound
	}

	t, err := s.LoadAccess(refreshToken.AccessToken.AccessToken)
	if err != nil {
		return nil, osin.ErrNotFound
	}

	return t, nil
}

func (s *GORMStorage) RemoveRefresh(token string) error {
	if err := s.db.Where("refresh_token = ?", token).Delete(&OAuth2RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

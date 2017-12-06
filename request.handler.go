package main

import (
	"github.com/gin-gonic/gin"
	"github.com/RangelReale/osin"
	"errors"
)

type RequestHandler interface {
	GetUser(c *gin.Context) (*User, error)
}

type APIRequestHandler struct{}

func NewRequestHandler() RequestHandler {
	return &APIRequestHandler{}
}

func (h *APIRequestHandler) GetUser(c *gin.Context) (*User, error) {
	var isToken bool
	var token *osin.AccessData

	t, exists := c.Get("token")
	if !exists {
		return nil, errors.New("Token does not exist")
	}

	if token, isToken = t.(*osin.AccessData); !isToken {
		return nil, errors.New("Could not assert token is *osin.AccessData")
	}

	var isUser bool
	var user *User
	if user, isUser = token.UserData.(*User); !isUser {
		return nil, errors.New("Could not assert user is *User")
	}

	return user, nil
}

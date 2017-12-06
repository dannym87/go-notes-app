package main

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/RangelReale/osin"
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db              *gorm.DB
	oauth2Server    *osin.Server
	responseHandler ResponseHandler
	validator       *validator.Validate
}

func InitAuthHandler(app *App) *AuthHandler {
	h := &AuthHandler{
		app.Db(),
		app.OAuth2Server(),
		app.ResponseHandler(),
		app.Validator(),
	}

	app.Engine().POST("/token", h.Token)

	return h
}

func (h *AuthHandler) Token(c *gin.Context) {
	resp := h.oauth2Server.NewResponse()
	defer resp.Close()

	if ar := h.oauth2Server.HandleAccessRequest(resp, c.Request); ar != nil {
		switch ar.Type {
		case osin.PASSWORD:
			data := struct {
				GrantType    string `form:"grant_type" validate:"required"`
				Username     string `form:"username" validate:"required"`
				Password     string `form:"password" validate:"required"`
				Scope        string `form:"scope" validate:"omitempty"`
				ClientId     string `form:"client_id" validate:"required"`
				ClientSecret string `form:"client_secret" validate:"required"`
			}{}

			if err := c.ShouldBind(&data); err != nil {
				// todo not malformed json. fix error
				h.responseHandler.MalformedJSON(c)
				return
			}

			if err := h.validator.Struct(data); err != nil {
				h.responseHandler.ValidationErrors(c, err)
				return
			}

			user := new(User)

			if err := h.db.Where("email = ?", data.Username).Find(user).Error; err != nil {
				h.responseHandler.Error(c, AuthenticationError, http.StatusBadRequest, "Username or Password is incorrect")
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
				h.responseHandler.Error(c, AuthenticationError, http.StatusBadRequest, "Username or Password is incorrect")
				return
			}

			ar.UserData = user
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		}

		h.oauth2Server.FinishAccessRequest(resp, c.Request, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		h.responseHandler.Error(c, resp.ErrorId, resp.ErrorStatusCode, resp.InternalError.Error())
		return
	}

	if resp.IsError {
		h.responseHandler.Error(c, resp.ErrorId, resp.ErrorStatusCode, resp.StatusText)
		return
	}

	if !resp.IsError {
		resp.StatusCode = http.StatusCreated
	}

	osin.OutputJSON(resp, c.Writer, c.Request)
}

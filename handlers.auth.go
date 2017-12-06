package main

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/RangelReale/osin"
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"golang.org/x/crypto/bcrypt"
	"fmt"
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

	app.Engine().POST("/token", h.LoginPassword)
	app.engine.GET("/info", h.Info)

	return h
}

func (h *AuthHandler) Info(c *gin.Context) {
	resp := h.oauth2Server.NewResponse()
	defer resp.Close()

	if ir := h.oauth2Server.HandleInfoRequest(resp, c.Request); ir != nil {
		fmt.Println(ir.AccessData.UserData)
		h.oauth2Server.FinishInfoRequest(resp, c.Request, ir)
	}

	osin.OutputJSON(resp, c.Writer, c.Request)
}

func (h *AuthHandler) LoginPassword(c *gin.Context) {
	resp := h.oauth2Server.NewResponse()
	defer resp.Close()

	data := struct {
		Username string `form:"username" validate:"required"`
		Password string `form:"password" validate:"required"`
		Scope    string `form:"scope" validate:"omitempty"`
	}{}

	if err := c.BindQuery(&data); err != nil {
		// todo not malformed json. fix error
		h.responseHandler.MalformedJSON(c)
		return
	}

	if err := h.validator.Struct(data); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	if ar := h.oauth2Server.HandleAccessRequest(resp, c.Request); ar != nil {
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

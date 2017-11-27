package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"gopkg.in/go-playground/validator.v9"
	"fmt"
)

type TagsHandler struct {
	db              *gorm.DB
	tagRepository   TagRepository
	responseHandler ResponseHandler
	validator       *validator.Validate
}

func InitTagsHandler(app *App) *TagsHandler {
	h := &TagsHandler{
		app.Db(),
		NewTagRepository(app.Db()),
		app.ResponseHandler(),
		app.Validator(),
	}

	v1 := app.engine.Group("/v1")
	{
		v1.GET("/tags", h.List)
		v1.GET("/tags/:id", h.Get)
		v1.POST("/tags", h.Create)
		v1.DELETE("/tags/:id", h.Delete)
		v1.PATCH("/tags/:id", h.Update)
	}

	return h
}

func (h *TagsHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10
	offset := (page * limit) - limit

	tags, err := h.tagRepository.FindAll(limit, offset)
	if err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, tags)
}

func (h *TagsHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := h.tagRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, tag)
}

func (h *TagsHandler) Create(c *gin.Context) {
	t := new(Tag)

	if err := c.BindJSON(t); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	if err := h.validator.Struct(t); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	if _, err := h.tagRepository.FindByName(t.Name); err == nil {
		h.responseHandler.Error(c, ValidationError, http.StatusUnprocessableEntity, fmt.Sprintf("Tag '%s' already exists", t.Name))
		return
	}

	tag, err := h.tagRepository.Create(t)
	if err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, tag)
}

func (h *TagsHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := h.tagRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	if err := h.tagRepository.Delete(tag); err != nil {
		h.responseHandler.InternalServerError(c)
	}

	h.responseHandler.JSON(c, http.StatusNoContent, "")
}

func (h *TagsHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := h.tagRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	if err := c.BindJSON(t); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	if err := h.validator.Struct(t); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	tagExists, err := h.tagRepository.FindByName(t.Name)
	if err == nil && tagExists.ID != uint(id) {
		h.responseHandler.Error(c, ValidationError, http.StatusUnprocessableEntity, fmt.Sprintf("Tag '%s' already exists", t.Name))
		return
	}

	tag, err := h.tagRepository.Update(id, t)
	if err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, tag)
}

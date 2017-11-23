package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"gopkg.in/go-playground/validator.v9"
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
	tag := new(Tag)

	if err := c.BindJSON(tag); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	tag.Id = 0

	if err := h.validator.Struct(tag); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	if err := h.db.Create(tag).Error; err != nil {
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

	if err := h.db.Delete(tag).Error; err != nil {
		h.responseHandler.InternalServerError(c)
	}

	c.Data(http.StatusNoContent, jsonapi.MediaType, []byte(""))
}

func (h *TagsHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := h.tagRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	if err := c.BindJSON(tag); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	tag.Id = uint(id)

	if err := h.validator.Struct(tag); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	if err := h.db.Save(tag).Error; err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, tag)
}

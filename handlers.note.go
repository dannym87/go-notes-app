package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"gopkg.in/go-playground/validator.v9"
)

type NotesHandler struct {
	db              *gorm.DB
	noteRepository  NoteRepository
	tagRepository   TagRepository
	responseHandler ResponseHandler
	validator       *validator.Validate
}

func InitNotesHandler(app *App) *NotesHandler {
	h := &NotesHandler{
		app.Db(),
		NewNoteRepository(app.Db()),
		NewTagRepository(app.Db()),
		app.ResponseHandler(),
		app.Validator(),
	}

	v1 := app.engine.Group("/v1")
	{
		v1.GET("/notes", h.List)
		v1.GET("/notes/:id", h.Get)
		v1.POST("/notes", h.Create)
		v1.DELETE("/notes/:id", h.Delete)
		v1.PATCH("/notes/:id", h.Update)
	}

	return h
}

func (h *NotesHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10
	offset := (page * limit) - limit

	notes, err := h.noteRepository.FindAll(limit, offset)

	if err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, notes)
}

func (h *NotesHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	note, err := h.noteRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, note)
}

func (h *NotesHandler) Create(c *gin.Context) {
	note := new(Note)

	if err := c.BindJSON(note); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	note.Id = 0
	h.mapTags(note)

	if err := h.validator.Struct(note); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	if err := h.db.Create(note).Error; err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, note)
}

func (h *NotesHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	note, err := h.noteRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	if err := h.db.Delete(note).Error; err != nil {
		h.responseHandler.InternalServerError(c)
	}

	c.Data(http.StatusNoContent, jsonapi.MediaType, []byte(""))
}

func (h *NotesHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	note, err := h.noteRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	if err := c.BindJSON(note); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	note.Id = uint(id)
	h.mapTags(note)

	if err := h.validator.Struct(note); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	if err := h.db.Save(note).Error; err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, note)
}

func (h *NotesHandler) mapTags(note *Note) {
	for _, tag := range note.Tags {
		t, err := h.tagRepository.FindById(int(tag.Id))
		if err == nil {
			tag.ExchangeData(t)
			continue
		}

		t, err = h.tagRepository.FindByName(tag.Name)
		if err == nil {
			tag.ExchangeData(t)
			continue
		}
	}

	h.db.Model(note).Association("Tags").Replace(note.Tags)
}

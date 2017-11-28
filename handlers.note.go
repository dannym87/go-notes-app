package main

import (
	"github.com/gin-gonic/gin"
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
	n := new(Note)

	if err := c.BindJSON(n); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	if err := h.validator.Struct(n); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	note, err := h.noteRepository.Create(n)

	if err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusCreated, note)
}

func (h *NotesHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	note, err := h.noteRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	if err := h.noteRepository.Delete(note); err != nil {
		h.responseHandler.InternalServerError(c)
	}

	h.responseHandler.JSON(c, http.StatusNoContent, "")
}

func (h *NotesHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	n, err := h.noteRepository.FindById(id)

	if err != nil {
		h.responseHandler.NotFound(c)
		return
	}

	if err := c.BindJSON(n); err != nil {
		h.responseHandler.MalformedJSON(c)
		return
	}

	if err := h.validator.Struct(n); err != nil {
		h.responseHandler.ValidationErrors(c, err)
		return
	}

	note, err := h.noteRepository.Update(id, n)

	if err != nil {
		h.responseHandler.InternalServerError(c)
		return
	}

	h.responseHandler.JSON(c, http.StatusOK, note)
}

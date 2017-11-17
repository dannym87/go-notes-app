package main

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"github.com/google/jsonapi"
	"net/http"
	"github.com/jinzhu/gorm"
	"time"
)

type NotesHandler struct {
	engine *gin.Engine
	db     *gorm.DB
}

func InitNotesHandler(e *gin.Engine, db *gorm.DB) *NotesHandler {
	nh := NotesHandler{e, db}
	v1 := e.Group("/v1")
	{
		v1.GET("/notes", nh.ListNotesHandler)
		v1.GET("/notes/:id", nh.GetNoteHandler)
		v1.POST("/notes", nh.CreateNoteHandler)
		v1.DELETE("/notes/:id", nh.DeleteNoteHandler)
		v1.PATCH("/notes/:id", nh.UpdateNoteHandler)
	}

	return &nh
}

func (nh *NotesHandler) ListNotesHandler(c *gin.Context) {
	var notes []*Note

	if err := nh.db.Find(&notes).Error; err != nil {
		c.Data(400, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	response := new(bytes.Buffer)
	jsonapi.MarshalPayload(response, notes)
	c.Data(http.StatusOK, jsonapi.MediaType, response.Bytes())
}

func (nh *NotesHandler) GetNoteHandler(c *gin.Context) {
	var note Note

	if err := nh.db.Where("id = ?", c.Param("id")).Find(&note).Error; err != nil {
		c.Data(400, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	response := new(bytes.Buffer)
	jsonapi.MarshalPayload(response, &note)
	c.Data(http.StatusOK, jsonapi.MediaType, response.Bytes())
}

func (nh *NotesHandler) CreateNoteHandler(c *gin.Context) {
	var note Note
	payload, _ := c.GetRawData()

	request := bytes.NewBuffer(payload)
	if err := jsonapi.UnmarshalPayload(request, &note); err != nil {
		c.Data(http.StatusBadRequest, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	note.Created = time.Now()

	if err := nh.db.Create(&note).Error; err != nil {
		c.Data(http.StatusBadRequest, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	response := new(bytes.Buffer)
	jsonapi.MarshalPayload(response, &note)
	c.Data(http.StatusCreated, jsonapi.MediaType, response.Bytes())
}

func (nh *NotesHandler) DeleteNoteHandler(c *gin.Context) {
	var note Note

	if err := nh.db.Where("id = ?", c.Param("id")).Find(&note).Error; err != nil {
		c.Data(http.StatusBadRequest, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	if err := nh.db.Delete(&note).Error; err != nil {
		c.Data(http.StatusBadRequest, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	c.Data(http.StatusNoContent, jsonapi.MediaType, []byte(""))
}

func (nh *NotesHandler) UpdateNoteHandler(c *gin.Context) {
	var note Note

	if err := nh.db.Where("id = ?", c.Param("id")).Find(&note).Error; err != nil {
		c.Data(http.StatusBadRequest, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	payload, _ := c.GetRawData()
	request := bytes.NewBuffer(payload)

	if err := jsonapi.UnmarshalPayload(request, &note); err != nil {
		c.Data(http.StatusBadRequest, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	if err := nh.db.Save(&note).Error; err != nil {
		c.Data(http.StatusBadRequest, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	response := new(bytes.Buffer)
	jsonapi.MarshalPayload(response, &note)
	c.Data(http.StatusOK, jsonapi.MediaType, response.Bytes())
}

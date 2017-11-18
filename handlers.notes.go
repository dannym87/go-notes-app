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
		v1.GET("/notes", nh.List)
		v1.GET("/notes/:id", nh.Get)
		v1.POST("/notes", nh.Create)
		v1.DELETE("/notes/:id", nh.Delete)
		v1.PATCH("/notes/:id", nh.Update)
	}

	return &nh
}

func (nh *NotesHandler) List(c *gin.Context) {
	var notes []*Note

	if err := nh.db.Find(&notes).Error; err != nil {
		c.Data(400, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	response := new(bytes.Buffer)
	jsonapi.MarshalPayload(response, notes)
	c.Data(http.StatusOK, jsonapi.MediaType, response.Bytes())
}

func (nh *NotesHandler) Get(c *gin.Context) {
	var note Note

	if err := nh.db.Where("id = ?", c.Param("id")).Find(&note).Error; err != nil {
		c.Data(400, jsonapi.MediaType, []byte(err.Error()))
		return
	}

	response := new(bytes.Buffer)
	jsonapi.MarshalPayload(response, &note)
	c.Data(http.StatusOK, jsonapi.MediaType, response.Bytes())
}

func (nh *NotesHandler) Create(c *gin.Context) {
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

func (nh *NotesHandler) Delete(c *gin.Context) {
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

func (nh *NotesHandler) Update(c *gin.Context) {
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

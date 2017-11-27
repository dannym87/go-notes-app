package main

type Tag struct {
	BaseModel
	Name string `json:"name" validate:"required,dbunique=tag.name"`
}

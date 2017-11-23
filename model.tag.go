package main

import "time"

type Tag struct {
	Id        uint      `json:"tags"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tag) ExchangeData(tag *Tag) *Tag {
	t.Id = tag.Id
	t.Name = tag.Name
	t.CreatedAt = tag.CreatedAt

	return t
}

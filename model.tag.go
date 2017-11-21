package main

import "time"

type Tag struct {
	Id        uint      `jsonapi:"primary,tags"`
	Name      string    `jsonapi:"attr,name" validate:"required"`
	CreatedAt time.Time `jsonapi:"attr,created_at,iso8601"`
}

func (t *Tag) ExchangeData(tag *Tag) *Tag {
	t.Id = tag.Id
	t.Name = tag.Name
	t.CreatedAt = tag.CreatedAt

	return t
}

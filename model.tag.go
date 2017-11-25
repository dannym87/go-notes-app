package main

type Tag struct {
	BaseModel
	Name string `json:"name" validate:"required"`
}

func (t *Tag) ExchangeData(tag *Tag) *Tag {
	t.ID = tag.ID
	t.Name = tag.Name
	t.CreatedAt = tag.CreatedAt

	return t
}

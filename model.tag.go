package main

import "github.com/jinzhu/gorm"

type Tag struct {
	BaseModel
	Name string `json:"name" validate:"required"`
}

func (t *Tag) BeforeDelete(tx *gorm.DB) {
	tx.Exec("DELETE FROM note_tags WHERE tag_id = ?", t.ID)
}

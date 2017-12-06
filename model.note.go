package main

import (
	"github.com/jinzhu/gorm"
)

type Note struct {
	BaseModel
	Title       string `json:"title" validate:"required"`
	Text        string `json:"text" validate:"omitempty"`
	Tags        []*Tag `json:"tags,omitempty" gorm:"many2many:note_tags;" validate:"omitempty,dive,required"`
	CreatedBy   *User  `json:"created_by" gorm:"ForeignKey:CreatedById"`
	CreatedById uint   `json:"-" gorm:"column:created_by"`
}

/*
 * GORM Event Callbacks
 */
func (n *Note) BeforeDelete(tx *gorm.DB) {
	tx.Model(n).Association("Tags").Clear()
}

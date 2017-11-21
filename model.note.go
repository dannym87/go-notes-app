package main

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Note struct {
	Id        uint      `jsonapi:"primary,notes"`
	Title     string    `jsonapi:"attr,title" validate:"required"`
	Text      string    `jsonapi:"attr,text" validate:"omitempty"`
	Tags      []*Tag    `jsonapi:"relation,tags,omitempty" gorm:"many2many:note_tags;" validate:"omitempty,dive,required"`
	CreatedAt time.Time `jsonapi:"attr,created_at,iso8601"`
}

/*
 * GORM Event Callbacks
 */
func (n *Note) BeforeDelete(tx *gorm.DB) {
	tx.Model(n).Association("Tags").Clear()
}

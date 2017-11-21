package main

import (
	"github.com/jinzhu/gorm"
)

type NoteRepository interface {
	FindById(id int) (*Note, error)
	FindAll(limit int, offset int) ([]*Note, error)
}

type ORMNoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &ORMNoteRepository{db}
}

func (r ORMNoteRepository) FindById(id int) (*Note, error) {
	note := new(Note)

	if err := r.db.Preload("Tags").First(note, id).Error; err != nil {
		return nil, err
	}

	return note, nil
}

func (r ORMNoteRepository) FindAll(limit int, offset int) ([]*Note, error) {
	var notes []*Note

	if err := r.db.Preload("Tags").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

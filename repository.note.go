package main

import (
	"github.com/jinzhu/gorm"
)

type NoteRepository interface {
	FindById(id int) (*Note, error)
	FindAll(limit int, offset int) ([]*Note, error)
	Create(n *Note) (*Note, error)
	Update(id int, n *Note) (*Note, error)
}

type ORMNoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &ORMNoteRepository{db}
}

func (r *ORMNoteRepository) FindById(id int) (*Note, error) {
	note := new(Note)

	if err := r.db.Preload("Tags").First(note, id).Error; err != nil {
		return nil, err
	}

	return note, nil
}

func (r *ORMNoteRepository) FindAll(limit int, offset int) ([]*Note, error) {
	var notes []*Note

	if err := r.db.Preload("Tags").Limit(limit).Offset(offset).Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *ORMNoteRepository) Create(n *Note) (*Note, error) {
	note := &Note{
		Title: n.Title,
		Text:  n.Text,
	}

	if err := r.db.Create(note).Error; err != nil {
		return n, err
	}

	note.Tags = append(note.Tags, n.Tags...)
	r.SaveTags(note)

	return note, nil
}

func (r *ORMNoteRepository) Update(id int, n *Note) (*Note, error) {
	note, _ := r.FindById(id)

	if err := r.db.Model(note).UpdateColumns(&Note{Title: n.Title, Text: n.Text}).Error; err != nil {
		return n, err
	}

	note.Tags = nil
	note.Tags = append(note.Tags, n.Tags...)
	r.SaveTags(note)

	return note, nil
}

func (r *ORMNoteRepository) SaveTags(n *Note) {
	var tags []*Tag

	for _, tag := range n.Tags {
		t := new(Tag)

		if err := r.db.Where("name = ?", tag.Name).Find(t).Error; err == nil {
			tags = append(tags, t)
			continue
		}

		t.Name = tag.Name
		tags = append(tags, t)
	}

	r.db.Model(n).Association("Tags").Replace(&tags)
}

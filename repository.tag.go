package main

import "github.com/jinzhu/gorm"

type TagRepository interface {
	FindById(id int) (*Tag, error)
	FindByName(name string) (*Tag, error)
	FindAll() ([]*Tag, error)
}

type ORMTagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &ORMTagRepository{db}
}

func (r *ORMTagRepository) FindById(id int) (*Tag, error) {
	tag := new(Tag)

	if err := r.db.Where("id = ?", id).Find(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *ORMTagRepository) FindByName(name string) (*Tag, error) {
	tag := new(Tag)

	if err := r.db.Where("name = ?", name).Find(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *ORMTagRepository) FindAll() ([]*Tag, error) {
	var tags []*Tag

	if err := r.db.Find(tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

package main

import "github.com/jinzhu/gorm"

type TagRepository interface {
	FindById(id int) (*Tag, error)
	FindByName(name string) (*Tag, error)
	FindAll(limit int, offset int) ([]*Tag, error)
}

type ORMTagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &ORMTagRepository{db}
}

func (r *ORMTagRepository) FindById(id int) (*Tag, error) {
	tag := new(Tag)

	if err := r.db.First(tag, id).Error; err != nil {
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

func (r *ORMTagRepository) FindAll(limit int, offset int) ([]*Tag, error) {
	var tags []*Tag

	if err := r.db.Limit(limit).Offset(offset).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

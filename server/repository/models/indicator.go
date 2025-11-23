package models

import "gorm.io/gorm"

type Indicator struct {
	gorm.Model
	Name string `gorm:"unique"`
	Type string
}

type Repository struct {
	db *gorm.DB
}

func NewIndicator(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindOne(name string) (*Indicator, error) {
	var i Indicator
	tx := r.db.First(&i, "name = ?", name)

	return &i, tx.Error
}

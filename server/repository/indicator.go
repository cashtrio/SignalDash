package repository

import (
	"github.com/developerasun/SignalDash/server/models"
	"gorm.io/gorm"
)

type IndicatorRepository struct {
	db *gorm.DB
}

func NewIndicator(db *gorm.DB) *IndicatorRepository {
	return &IndicatorRepository{db: db}
}

func (r *IndicatorRepository) FindOne(name string) (*models.Indicator, error) {
	var i models.Indicator
	tx := r.db.First(&i, "name = ?", name)

	return &i, tx.Error
}

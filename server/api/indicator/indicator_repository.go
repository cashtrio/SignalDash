package indicator

import (
	"github.com/developerasun/SignalDash/server/models"
	"gorm.io/gorm"
)

type IndicatorRepo struct {
	db *gorm.DB
}

func (ir *IndicatorRepo) FindOne(ticker string) (*models.Indicator, error) {
	var i models.Indicator
	tx := ir.db.First(&i, "ticker = ?", ticker)
	return &i, tx.Error
}

func (ir *IndicatorRepo) Create(i *models.Indicator) error {
	if tx := ir.db.Create(i); tx.Error != nil {
		return tx.Error
	}
	return nil
}

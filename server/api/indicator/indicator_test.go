package indicator

import (
	"testing"

	"github.com/developerasun/SignalDash/server/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIndicatorInsert(t *testing.T) {
	assert := assert.New(t)
	testdb, oErr := gorm.Open(sqlite.Open("../../sd_app.test.db"), &gorm.Config{})

	if oErr != nil {
		t.Logf("failed to open test db: %s", oErr.Error())
		t.Fail()
	}
	testdb.AutoMigrate(&models.Indicator{})
	indicatorRepo := IndicatorRepo{db: testdb}
	cErr := indicatorRepo.Create(&models.Indicator{
		Name:   "apple",
		Ticker: "aapl",
		Value:  1.00,
		Domain: "apple.com",
	})

	assert.Nil(cErr)
}

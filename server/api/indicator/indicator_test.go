package indicator

import (
	"testing"
	"time"

	"github.com/developerasun/SignalDash/server/config"
	"github.com/developerasun/SignalDash/server/models"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIndicatorInsert(t *testing.T) {
	t.Skip()
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

func TestCronExpression(t *testing.T) {
	assert := assert.New(t)
	c := cron.New(cron.WithSeconds())

	v := config.NewEnvironment("../../config", "options").Instance
	expression := v.GetString("server.cron.expression.every1sec")
	require.NotEqual(t, 0, len(expression))

	count := 0
	c.AddFunc(expression, func() {
		count++
	})
	c.Start()
	<-time.After(time.Second * 2)
	c.Stop()

	expected := 2
	assert.Equal(expected, count)
}

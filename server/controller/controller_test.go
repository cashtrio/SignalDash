package controller

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/developerasun/SignalDash/server/config"
	"github.com/developerasun/SignalDash/server/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup(t *testing.T) *gin.Engine {
	absPath, err := filepath.Abs("../config")
	require.NoError(t, err)
	v := config.NewEnvironment(absPath, "options").Instance
	testDB := v.GetString("server.database.test")

	dbPath, err := filepath.Abs("../" + testDB)
	require.NoError(t, err)
	db, oErr := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, oErr)

	db.AutoMigrate(&models.Indicator{})

	r := gin.Default()
	r.GET("/api/health", Health)
	r.GET("/api/indicator", func(ctx *gin.Context) {
		ScrapeDollarIndex(ctx, db)
	})
	r.GET("/api/exchange-rate", func(ctx *gin.Context) {
		CreateExchangeRateDiff(ctx, db)
	})

	return r
}

func TestApis(t *testing.T) {
	r := setup(t)

	t.Run("health", func(t *testing.T) {
		t.Skip()
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		r.ServeHTTP(w, req)

		require.True(t, w.Code == http.StatusOK)
	})

	t.Run("dollar index", func(t *testing.T) {
		t.Skip()
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/indicator", nil)
		r.ServeHTTP(w, req)

		require.True(t, w.Code == http.StatusOK)
	})

	t.Run("krw usdt diff", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/exchange-rate", nil)
		r.ServeHTTP(w, req)

		require.True(t, w.Code == http.StatusOK)
	})
}

package controller

import (
	"fmt"
	"net/http"

	"github.com/developerasun/SignalDash/server/dto"
	"github.com/developerasun/SignalDash/server/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ScrapeDollarIndex godoc
// @Summary crawl financial indicators
// @Description crawl tradingview for dxy and insert to db
// @Tags api
// @Produce json
// @Success 200 {object} dto.ScrapeDollarIndexResponse
// @Router /api/indicator [get]
func ScrapeDollarIndex(ctx *gin.Context, db *gorm.DB) {
	indicator := service.NewIndicator([]string{
		"www.tradingview.com", "tradingview.com",
	}, "Mozilla/5.0 (compatible; DeveloperAsunBot/1.0)")

	dxy, sErr := indicator.ScrapeDollarIndex()
	if sErr != nil {
		ctx.Error(sErr)
		return
	}

	cErr := service.CreateDollarIndex(db, dxy)
	if cErr != nil {
		ctx.Error(cErr)
		return
	}

	ctx.JSON(http.StatusOK, dto.ScrapeDollarIndexResponse{
		DollarIndex: dxy,
	})
}

// CreateExchangeRateDiff godoc
// @Summary calculate exchange rate premium
// @Description fetch KRW and USDT price and get price diff
// @Tags api
// @Produce json
// @Success 200 {object} dto.CreateExchangeRateDiffResponse
// @Router /api/exchange-rate [get]
func CreateExchangeRateDiff(ctx *gin.Context, db *gorm.DB) {
	won, tether, err := service.CreateExchangeRateDiff()
	if err != nil {
		ctx.Error(err)
		return
	}

	var strong string
	diff := won - tether

	if diff > 0 {
		strong = "KRW"
	} else if diff == 0 {
		strong = "same"
	} else {
		strong = "USDT"
	}

	ctx.JSON(http.StatusOK, dto.CreateExchangeRateDiffResponse{
		ExchangeRateDiff: fmt.Sprintf("%.2f", diff),
		Strong:           strong,
	})
}

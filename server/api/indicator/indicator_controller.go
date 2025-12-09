package indicator

import (
	"net/http"

	"github.com/developerasun/SignalDash/server/dto"
	"github.com/developerasun/SignalDash/server/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetIndicator godoc
// @Summary crawl financial indicators
// @Description crawl tradingview for dxy and insert to db
// @Tags api
// @Produce json
// @Success 200 {object} dto.OkResponse
// @Router /api/indicator [get]
func ScrapeDollarIndex(ctx *gin.Context, db *gorm.DB) {
	indicator := service.NewIndicator([]string{
		"www.tradingview.com", "tradingview.com",
	}, "Mozilla/5.0 (compatible; DeveloperAsunBot/1.0)")

	dxy, sErr := indicator.ScrapeDollarIndex()
	if sErr != nil {
		ctx.Error(sErr)
	}

	ctx.JSON(http.StatusOK, dto.ScrapeDollarIndexResponse{
		DollarIndex: dxy,
	})
}

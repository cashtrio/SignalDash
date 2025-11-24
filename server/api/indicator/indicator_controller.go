package indicator

import (
	"net/http"

	"github.com/developerasun/SignalDash/server/dto"
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
func GetIndicator(ctx *gin.Context, db *gorm.DB) {
	repo := IndicatorRepo{db: db}
	service := IndicatorService{repo: &repo}
	ciErr := service.CrawlAndInsert()

	if ciErr != nil {
		ctx.Error(ciErr)
	}

	ctx.JSON(http.StatusOK, dto.OkResponse{
		Message: "ok",
	})
}

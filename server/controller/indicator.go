package controller

import (
	"net/http"

	"github.com/developerasun/SignalDash/server/service"
	"github.com/gin-gonic/gin"
)

type IndicatorController struct {
	service *service.IndicatorService
}

func NewIndicatorController(s *service.IndicatorService) *IndicatorController {
	return &IndicatorController{
		service: s,
	}
}

// Health godoc
// @Summary dummy test
// @Description dummy test
// @Tags api
// @Produce json
// @Success 200 {object} any
// @Router /api/indicator [get]
func (c *IndicatorController) GetIndicator(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

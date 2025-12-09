package instance

import (
	"fmt"
	"net/http"

	"github.com/developerasun/SignalDash/server/api/health"
	"github.com/developerasun/SignalDash/server/api/indicator"
	"github.com/developerasun/SignalDash/server/docs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type apiServer struct {
	engine *gin.Engine
}

func NewApiServer(restApi *gin.Engine, db *gorm.DB) *apiServer {
	restApi.SetTrustedProxies(nil)

	docs.SwaggerInfo.BasePath = ""
	restApi.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	restApi.Use(ErrorHandler())
	restApi.Use(gin.Recovery())

	// TODO refactor grouping with controller package
	api := restApi.Group("/api")
	api.GET("/health", health.Health)
	api.GET("/indicator", func(ctx *gin.Context) {
		indicator.ScrapeDollarIndex(ctx, db)
	})

	return &apiServer{
		engine: restApi,
	}
}

func (a *apiServer) Run(addr ...string) (err error) {
	return a.engine.Run(addr...)
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error!",
			})
		}
	}
}

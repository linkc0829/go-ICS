package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/linkc0829/go-icsharing/internal/handlers"
	"github.com/linkc0829/go-icsharing/pkg/utils"
)

func FreeTrial(cfg *utils.ServerConfig, r *gin.Engine, sqlite *gorm.DB) {
	r.GET("/", handlers.FreeTrialHandler())
	log.Println("Free Trial Page @ " + cfg.ListenEndpoint())

	//free trial API
	trial := r.Group(cfg.VersioningEndpoint("/trial"))
	trial.GET("", handlers.GetPortfolioHandlers(cfg, sqlite))
	trial.POST("", handlers.CreatePortfolioHandlers(cfg, sqlite))
	trial.PATCH("/:id", handlers.UpdatePortfolioHandlers(cfg, sqlite))
	trial.DELETE("/:id", handlers.DeletePortfolioHandlers(cfg, sqlite))

}

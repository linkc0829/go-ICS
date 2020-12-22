package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handlers"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func FreeTrial(cfg *utils.ServerConfig, r *gin.Engine) {
	r.GET("/", handlers.FreeTrialHandler())
	log.Println("Free Trial Page @ " + cfg.ListenEndpoint())
}

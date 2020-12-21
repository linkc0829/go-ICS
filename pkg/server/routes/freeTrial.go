package routes

import(
	"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handlers"

)

func FreeTrail(cfg *utils.ServerConfig, r *gin.Engine){
	r.GET("/", handlers.FreeTrialHandler())
	log.Println("Free Trial Page @ " + cfg.ListenEndPoint())
}
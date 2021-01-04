package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handlers"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func UserProfolio(cfg *utils.ServerConfig, r *gin.Engine) {
	g := r.Group(cfg.VersioningEndpoint("/profolio"))
	g.GET("/:id", handlers.UserProfileHandler(cfg))
}

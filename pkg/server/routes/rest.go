package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handlers/client"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func RestAPI(cfg *utils.ServerConfig, r *gin.Engine) {

	g := r.Group(cfg.VersioningEndpoint("/user"))
	g.GET("/", client.GetUser())
	g.POST("/", client.CreateUser())
	g.POST("/:id", client.UpdateUser())
	g.DELETE("/:id", client.DeleteUSer())
	g.POST("/:id/friend", client.AddFriend())

	g.GET("/income", client.GetUserIncome(cfg))
	g.POST("/income", cient.CreateIncome())
	g.POST("/income/:id", client.UpdateIncome())
	g.DELETE("/income/:id", client.DeleteIncome())
	g.POST("/income/:id/vote", client.VoteIncome())

}

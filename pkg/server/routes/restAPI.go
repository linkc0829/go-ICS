package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handlers/gqlclient"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func RestAPI(cfg *utils.ServerConfig, r *gin.Engine) {

	g := r.Group(cfg.VersioningEndpoint("/user"))

	g.GET("/:id", gqlclient.GetUser(cfg))
	g.POST("", gqlclient.CreateUser(cfg))
	g.PATCH("/:id", gqlclient.UpdateUser(cfg))
	g.DELETE("/:id", gqlclient.DeleteUser(cfg))
	g.PUT("addfriend/:id", gqlclient.AddFriend(cfg))
	g.GET(":id/income", gqlclient.GetUserIncome(cfg))
	g.GET(":id/cost", gqlclient.GetUserCost(cfg))
	g.GET(":id/income/history", gqlclient.GetUserIncomeHistory(cfg))
	g.GET(":id/cost/history", gqlclient.GetUserCostHistory(cfg))

	in := r.Group(cfg.VersioningEndpoint("/income"))
	in.POST("", gqlclient.CreateIncome(cfg))
	in.PATCH(":id", gqlclient.UpdateIncome(cfg))
	in.DELETE(":id", gqlclient.DeleteIncome(cfg))
	in.PUT("vote/:id", gqlclient.VoteIncome(cfg))

	co := r.Group(cfg.VersioningEndpoint("/cost"))
	co.POST("", gqlclient.CreateCost(cfg))
	co.PATCH(":id", gqlclient.UpdateCost(cfg))
	co.DELETE(":id", gqlclient.DeleteCost(cfg))
	co.PUT("vote/:id", gqlclient.VoteCost(cfg))

}

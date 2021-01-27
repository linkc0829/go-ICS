package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-icsharing/internal/handlers"
	"github.com/linkc0829/go-icsharing/internal/handlers/auth"
	"github.com/linkc0829/go-icsharing/internal/handlers/secret"
	"github.com/linkc0829/go-icsharing/pkg/utils"
	"github.com/linkc0829/go-icsharing/pkg/utils/datasource"
)

func Auth(cfg *utils.ServerConfig, r *gin.Engine, db *datasource.DB) {

	// OAuth handlers
	g := r.Group(cfg.VersioningEndpoint("/auth"))
	g.GET("/:provider", auth.Begin())
	g.GET("/:provider/callback", auth.CallBack(cfg, db))

	// ics secrets handler
	g.POST("/:provider/signup", secret.SignupHandler(cfg, db))
	g.POST("/:provider/login", secret.LoginHandler(cfg, db))
	g.GET("/:provider/refresh_token", secret.RefreshTokenHandler(cfg, db))

	g.GET("/:provider/logout", handlers.LogoutHandler())

}

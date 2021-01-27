package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-icsharing/internal/db/mongodb"
	"github.com/linkc0829/go-icsharing/internal/handlers"
	auth "github.com/linkc0829/go-icsharing/internal/handlers/auth/middleware"
	"github.com/linkc0829/go-icsharing/pkg/utils"
)

//Graph starts the graphql server routes
func Graph(cfg *utils.ServerConfig, r *gin.Engine, db *mongodb.MongoDB) {

	gqlPath := cfg.VersioningEndpoint(cfg.GraphQL.Path)
	gqlPGPath := cfg.GraphQL.PlaygroundPath
	g := r.Group(gqlPath)

	//Graph Handler
	g.POST("", auth.Middleware(g.BasePath(), cfg, db), handlers.GraphqlHandler(db))
	//g.POST("", handlers.GraphqlHandler(db))
	log.Println("Graphql Server @ " + g.BasePath())

	if cfg.GraphQL.IsPlaygroundEnabled {
		g.GET(gqlPGPath, handlers.PlaygroundHandler(g.BasePath()))
		log.Println("GraphQL Playground @ " + g.BasePath() + gqlPGPath)
	}
}

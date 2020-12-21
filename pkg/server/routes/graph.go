package routes

import(
	"github.com/linkc0829/go-ics/internal/handlers"
	auth "github.com/linkc0829/go-ics/internal/handlers/auth/middleware"
	"github.com/linkc0829/go-ics/internal/mongodb"
	"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/gin-gonic/gin"
)


func Graph(cfg *utils.ServerConfig, r *gin.Engine, db *mongodb.MongoDB){

	gqlPath := cfg.VersioningEndPoint(cfg.GraphQL.Path)
	gqlPGPath := cfg.GraphQL.PlaygroundPath
	g := r.Group(gqlPath)

	//Graph Handler
	g.POST("", auth.Middleware(g.BasePath(), cfg, db), handlers.GraphqlHandler(db))
	log.Println("Graphql Server @ " + g.BasePath())

	if cfg.GraphQL.IsPlaygroundEnabled{
		g.GET(gqlPgPath, handlers.PlaygroundHandler(g.BasePath()))
		log.Println("GraphQL Playground @ " + g.BasePath() + gqlPath)
	}
}
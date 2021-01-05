package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handlers"
	"github.com/linkc0829/go-ics/internal/mongodb"
	"github.com/linkc0829/go-ics/pkg/server/routes"
	"github.com/linkc0829/go-ics/pkg/utils"
)

var host, port, gqlPath, gqlPgPath string
var isPgEnabled bool

func RegisterRoutes(cfg *utils.ServerConfig, r *gin.Engine, db *mongodb.MongoDB) {
	routes.Auth(cfg, r, db)
	routes.Graph(cfg, r, db)
	routes.FreeTrial(cfg, r)
	routes.RestAPI(cfg, r)
	r.POST("/profile/:id", handlers.UserProfileHandler(cfg))
	r.GET("/profile/:id", handlers.UserProfileHandler(cfg))
	r.GET("/history/:id", handlers.UserHistoryHandler(cfg))

}

//Run will steup the routes and start the server
func Run(serverconf *utils.ServerConfig, db *mongodb.MongoDB) {

	r := gin.Default()

	InitalizeAuthProviders(serverconf)
	RegisterRoutes(serverconf, r, db)

	// Inform the user where the server is listening
	r.LoadHTMLGlob("views/*")
	r.Static(serverconf.StaticPath, "./public")

	log.Println("Running @ " + serverconf.ListenEndpoint())

	// Run the server
	// Print out and exit(1) to the OS if the server cannot run
	log.Fatal(r.Run(serverconf.ListenEndpoint()))

}

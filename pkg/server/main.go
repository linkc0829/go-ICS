package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/linkc0829/go-ics/internal/db/mongodb"
	"github.com/linkc0829/go-ics/pkg/server/routes"
	"github.com/linkc0829/go-ics/pkg/utils"
)

var host, port, gqlPath, gqlPgPath string
var isPgEnabled bool

func RegisterRoutes(cfg *utils.ServerConfig, r *gin.Engine, mongoDB *mongodb.MongoDB, sqlite *gorm.DB) {
	routes.Auth(cfg, r, mongoDB)
	routes.Graph(cfg, r, mongoDB)
	routes.FreeTrial(cfg, r, sqlite)
	routes.RestAPI(cfg, r, mongoDB)

}

//Run will steup the routes and start the server
func Run(serverconf *utils.ServerConfig, mongoDB *mongodb.MongoDB, sqlite *gorm.DB) {

	r := gin.Default()

	InitalizeAuthProviders(serverconf)
	RegisterRoutes(serverconf, r, mongoDB, sqlite)

	// Inform the user where the server is listening
	r.LoadHTMLGlob("views/*")
	r.Static(serverconf.StaticPath, "./public")

	log.Println("Running @ " + serverconf.ListenEndpoint())

	// Run the server
	// Print out and exit(1) to the OS if the server cannot run
	log.Fatal(r.Run(serverconf.ListenEndpoint()))

}

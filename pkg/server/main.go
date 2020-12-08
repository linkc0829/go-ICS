package server

import(
	"log"
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handler"
	"github.com/linkc0829/go-ics/pkg/utils"
)

var host, port string

func init(){
	host = utils.MustGet("GQL_SERVER_HOST")
	port = utils.MustGet("GQL_SERVER_PORT")
}

func Run(){
	
	pathGQL := "/graphql"

	r := gin.Default()
	r.GET("/", handler.FreeTrialHandler())


	log.Println("Connect to Graphql Server @ http://" + host + ":" + port + pathGQL)

	log.Fatalln(r.Run(host + ":" + port))

}
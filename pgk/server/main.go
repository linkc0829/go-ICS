package server

import(
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/linkc0829/go-ics/internal/handler"
)

var Host, Port string

func init(){
	HOST = "localhost"
	PORT = "8080"
}

func Run(){
	
	//pathGQL := "/graphql"

	r := gin.Default()
	r.GET("/", handler.IndexHandler)


	log.Println("Connect to Graphql Server @ http://" + host + ":" + port + pathGQL)

}
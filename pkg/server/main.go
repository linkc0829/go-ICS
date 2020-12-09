package server

import(
	"log"
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/handlers"
	"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/linkc0829/go-ics/internal/mongodb"
)

var host, port, gqlPath, gqlPgPath string
var isPgEnabled bool

func init(){
	host = utils.MustGet("GQL_SERVER_HOST")
	port = utils.MustGet("GQL_SERVER_PORT")
	gqlPath = utils.MustGet("GQL_SERVER_GRAPHQL_PATH")
    gqlPgPath = utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH")
    isPgEnabled = utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED")
}

func Run(db mongodb.MongoDB){

	endpoint := "http://" + host + ":" + port


	r := gin.Default()
	r.GET("/", handlers.FreeTrialHandler())
	log.Println("Free Trial Page @ " + endpoint)

	r.POST(gqlPath, handlers.GraphqlHandler(db))
	log.Println("Graphql Server @ " + endpoint + gqlPath)

	if isPgEnabled{
		r.GET(gqlPgPath, handlers.PlaygroundHandler(gqlPath))
		log.Println("GraphQL Playground @ " + endpoint + gqlPgPath)
	}

	log.Fatalln(r.Run(host + ":" + port))

}
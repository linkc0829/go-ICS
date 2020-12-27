package client

import (
	//"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"

	"golang.org/x/oauth2"

	//"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func GetUserIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		client := newClient(c, cfg)
		ID := c.Query("id")
		log.Println(client)
		var query struct {
			GetUserIncome []struct {
				Id    graphql.ID
				Owner struct {
					Id graphql.ID
				}
				Amount      graphql.Int
				Category    graphql.String
				OccurDate   graphql.String
				Description graphql.String
				Vote        []struct {
					Id graphql.ID
				}
			} `graphql:"getUserIncome(id: $ID)"`
		}
		variables := map[string]interface{}{
			"ID": graphql.ID(ID),
		}
		err := client.Query(c, &query, variables)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		log.Println(query)
		c.JSON(http.StatusOK, query)

	}
}

func newClient(c *gin.Context, cfg *utils.ServerConfig) *graphql.Client {
	accessToken := c.GetHeader("Authorization")

	//accessToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhvbWVVc2VyQGdtYWlsLmNvbSIsImV4cCI6MTYwOTAwNDk3NSwianRpIjoiaG9tZVVzZXIiLCJpYXQiOjE2MDkwMDEzNzUsImlzcyI6ImljcyIsIm5iZiI6MTYwOTAwMTM3NX0.PO5Vp0T-dL4PdEkcHfnbP5widVClwLK-xWYidNOl2MAvpeiPFTBAs9F3Spv84pX3Cf2Rz_Qh9NPMA1tUsqs7Bg"
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(c, src)

	gqlServerPath := cfg.SchemaVersioningEndpoint(cfg.GraphQL.Path)
	return graphql.NewClient(gqlServerPath, httpClient)
}

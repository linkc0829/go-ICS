package gqlclient

import (
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

func newClient(c *gin.Context, cfg *utils.ServerConfig) *graphql.Client {
	accessToken := c.GetHeader("Authorization")
	//log.Println("accessToken: " + accessToken)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(c, src)

	gqlServerPath := cfg.SchemaVersioningEndpoint(cfg.GraphQL.Path)
	return graphql.NewClient(gqlServerPath, httpClient)
}

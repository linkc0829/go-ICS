package gqlclient

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-icsharing/pkg/utils"
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

	//Skip verify SSL certificate if using SSL connection
	if cfg.URISchema == "https://" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		sslcli := &http.Client{Transport: tr}
		ctx := context.WithValue(c, oauth2.HTTPClient, sslcli)

		httpClient = oauth2.NewClient(ctx, src)
		gqlServerPath = cfg.RealSchemaVersioningEndpoint(cfg.GraphQL.Path)
	}

	return graphql.NewClient(gqlServerPath, httpClient)
}

//ErrorWriter set redirect header to index and show error message
func ErrorWriter(c *gin.Context, code int, err error) {
	//c.Writer.Header().Set("Location", "/")
	err = errors.New("[gql client] error: " + err.Error())
	c.Error(err)
	json := gin.H{
		"Title":        http.StatusText(code),
		"ErrorMessage": err.Error(),
	}
	c.JSON(code, json)
}

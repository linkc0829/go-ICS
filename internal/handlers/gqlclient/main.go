package gqlclient

import (
	"context"
	"errors"
	"net/http"

	"crypto/tls"

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
	//Skip verify SSL certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.WithValue(c, oauth2.HTTPClient, sslcli)

	httpClient := oauth2.NewClient(ctx, src)

	gqlServerPath := cfg.SchemaVersioningEndpoint(cfg.GraphQL.Path)
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

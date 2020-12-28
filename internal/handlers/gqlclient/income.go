package gqlclient

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/pkg/utils"
)

//"encoding/json"
//"github.com/linkc0829/go-ics/internal/graph/models"

//CreateIncome handle request POST /income
func CreateIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//UpdateIncome handle request POST /income/:id
func UpdateIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//DeleteIncome handle request POST /income
func DeleteIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//VoteIncome handle request PUT /income/vote/:id
func VoteIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

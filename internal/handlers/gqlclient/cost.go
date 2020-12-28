package gqlclient

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/pkg/utils"
)

//"encoding/json"
//"github.com/linkc0829/go-ics/internal/graph/models"

//CreateCost handle request POST /cost
func CreateCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//UpdateCost handle request POST /cost/:id
func UpdateCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//DeleteCost handle request POST /cost
func DeleteCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//VoteCost handle request PUT /cost/vote/:id
func VoteCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

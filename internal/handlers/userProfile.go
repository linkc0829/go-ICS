package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func UserProfileHandler(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		data := struct {
			Title string
		}{
			Title: "User Profile | ICS",
		}

		c.HTML(http.StatusOK, "profile", data)

	}
}

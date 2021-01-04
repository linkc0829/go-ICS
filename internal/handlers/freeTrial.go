package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FreeTrialHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		data := struct {
			Title        string
			ErrorMessage string
		}{
			Title:        "Free try Income & Cost SFunction here. Signup to share with friends.",
			ErrorMessage: "",
		}

		c.HTML(http.StatusOK, "layout", data)
	}

}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"

	"github.com/linkc0829/go-icsharing/internal/handlers/auth"
	"github.com/linkc0829/go-icsharing/pkg/utils"
)

//LogoutHandler will set header token expire and redirect to free trial page
func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param(string(utils.ProjectContextKeys.ProviderCtxKey))
		if provider == "ics" {
			c.Header("token_expiry", "-1")
		} else {
			c.Request = auth.AddProviderToContext(c, c.Param(string(utils.ProjectContextKeys.ProviderCtxKey)))
			gothic.Logout(c.Writer, c.Request)
		}
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
		//c.Writer.Header().Set("Location", "/")
		//c.Writer.WriteHeader(http.StatusTemporaryRedirect)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

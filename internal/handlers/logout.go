package handlers

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"

	"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/linkc0829/go-ics/internal/handlers/auth"
)

//LogoutHandler will set header token expire and redirect to free trial page
func LogoutHandler() gin.HandlerFunc{
	return func(c *gin.Context){
		provider := c.Param(string(utils.ProjectContextKeys.ProviderCtxKey))
		if provider == "ics" {
			c.Header("token_expiry", "-1")	
		} else {
			c.Request = auth.AddProviderToContext(c, c.Param(string(utils.ProjectContextKeys.ProviderCtxKey)))
		    gothic.Logout(c.Writer, c.Request)
		}
		c.Writer.Header().Set("Location", "/")
		c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}
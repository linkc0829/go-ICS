package handlers

import(
	"net/http"
    "github.com/gin-gonic/gin"
)


func FreeTrialHandler() gin.HandlerFunc{

	return func(c *gin.Context){
		c.JSON(http.StatusOK, "Free try Income & Cost Function here. Signup to share with friends.")
	}

}
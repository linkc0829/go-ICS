package auth

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/linkc0829/go-ics/internal/handlers/secret"

	"github.com/dgrijalva/jwt-go"
	"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/linkc0829/go-ics/pkg/utils/datasource"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// Claims JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Begin login with the auth provider
func Begin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = AddProviderToContext(c, c.Param("provider"))
		provider := c.Request.Context().Value("provider").(string)
		log.Println("Add Provider To Context: " + provider)
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		} else {
			log.Printf("user: %#v", gothUser)
		}
	}
}

// CallBack callback to complete auth provider flow
func CallBack(cfg *utils.ServerConfig, db *datasource.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = AddProviderToContext(c, c.Param("provider"))
		log.Println("CallBack adds provider to context")

		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		log.Println("CallBack CompleteUserAuth")

		u, err := db.Mongo.FindUserByJWT(user.Email, user.Provider, user.UserID)
		// logger.Infof("gothUser: %#v", user)
		if err != nil {
			if u, err = db.Mongo.CreateUserFromGoth(&user); err != nil {
				log.Println("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: " + err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		log.Println("[Auth.CallBack.UserLoggedIn]: ", u.ID)
		//generate new token pair
		accToken, tokenExpiry, refreshToken, err := secret.CreateTokenPair(cfg, u.ID.Hex(), user.Provider)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		//update redis
		refExp, _ := strconv.Atoi(cfg.JWT.RefreshTokenExpire[:len(cfg.JWT.RefreshTokenExpire)-1])
		db.Redis.Do("Set", u.ID.Hex(), refreshToken, "EX", refExp*3600)

		data := struct {
			TokenType   string
			Token       string
			TokenExpiry time.Time
			Redirect    string
		}{
			Token:       accToken,
			TokenType:   "Bearer",
			TokenExpiry: tokenExpiry,
			Redirect:    "/profile/" + u.ID.Hex(),
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.HTML(http.StatusOK, "callback", data)
		//c.Writer.Header().Set("Location", "/"+u.ID.Hex())
		//c.JSON(http.StatusPermanentRedirect, json)

	}
}

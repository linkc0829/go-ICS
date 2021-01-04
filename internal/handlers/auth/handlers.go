package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/linkc0829/go-ics/internal/handlers/secret"
	"github.com/linkc0829/go-ics/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/linkc0829/go-ics/pkg/utils"

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
func CallBack(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
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

		u, err := db.FindUserByJWT(user.Email, user.Provider, user.UserID)
		// logger.Infof("gothUser: %#v", user)
		if err != nil {
			if u, err = db.CreateUserFromGoth(&user); err != nil {
				log.Println("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: " + err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}

		log.Println("[Auth.CallBack.UserLoggedIn]: ", u.ID)

		//generate new token pair
		accToken, tokenExpiry, refreshToken, err := secret.CreateTokenPair(cfg, u)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		//update DB
		q := bson.M{"_id": u.ID}
		update := bson.M{"$set": bson.M{"refreshToken": refreshToken}}
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Users.UpdateOne(ctx, q, update)

		//set token
		json := gin.H{
			"type":         "Bearer",
			"token":        accToken,
			"token_expiry": tokenExpiry,
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.Writer.Header().Set("Location", "/"+u.ID.Hex())
		c.JSON(http.StatusPermanentRedirect, json)

	}
}

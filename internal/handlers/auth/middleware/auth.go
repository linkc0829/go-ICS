package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/linkc0829/go-ics/internal/db/mongodb"
	"github.com/linkc0829/go-ics/internal/db/mongodb/models"
	"github.com/linkc0829/go-ics/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

func authError(c *gin.Context, err error) {

	err = errors.New("[Auth] error: " + err.Error())
	c.Error(err)
	data := struct {
		Title        string
		ErrorMessage string
	}{
		Title:        http.StatusText(http.StatusUnauthorized),
		ErrorMessage: err.Error(),
	}
	c.HTML(http.StatusUnauthorized, "layout", data)
}

// Middleware wraps the request with auth middleware
func Middleware(path string, cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	log.Println("[Auth.Middleware] Applied to path: " + path)
	return gin.HandlerFunc(func(c *gin.Context) {

		//parse apiKey first
		if a, err := ParseAPIKey(c, cfg); err == nil {
			user, err := db.FindUserByAPIKey(a)
			if err != nil {
				authError(c, ErrForbidden)
			}
			log.Println("User: " + user.UserID)
			c.Next()

		} else {
			//if apiKey is empty, check jwt-token
			if err != ErrEmptyAPIKeyHeader {
				authError(c, err)
			} else {
				t, err := ParseToken(c, cfg)
				if err != nil {
					authError(c, err)
				} else {
					if t.Valid {
						if claims, ok := t.Claims.(jwt.MapClaims); ok {
							if claims["exp"] != nil {
								issuer := claims["iss"].(string)
								id := claims["id"].(string)
								primID, _ := primitive.ObjectIDFromHex(id)

								user := models.UserModel{}
								q := bson.M{"_id": primID, "provider": issuer}
								ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
								if err := db.Users.FindOne(ctx, q).Decode(&user); err != nil {
									authError(c, ErrForbidden)
								}
								c.Request = addToContext(c, utils.ProjectContextKeys.UserCtxKey, &user)
								c.Next()

							} else {
								authError(c, ErrMissingExpField)
							}

						} else {
							authError(c, ErrInvalidAccessToken)
						}
					}
				}
			}
		}
	})
}

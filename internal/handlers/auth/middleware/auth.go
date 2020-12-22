package middleware

import (
	"net/http"
	"log"

    "github.com/linkc0829/go-ics/internal/mongodb"
    "github.com/linkc0829/go-ics/pkg/utils"
    "github.com/dgrijalva/jwt-go"

    "github.com/gin-gonic/gin"
)

func authError(c *gin.Context, err error){
	e := gin.H{"message": "[Auth] error" + err.Error()}
	c.AbortWithStatusJSON(http.StatusUnauthorized, e)
}

// Middleware wraps the request with auth middleware
func Middleware(path string, cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	log.Println("[Auth.Middleware] Applied to path: " + path)
	return gin.HandlerFunc(func(c *gin.Context){

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
                                userid := claims["jti"].(string)
                                email := claims["email"].(string)
                                user, err := db.FindUserByJWT(email, issuer, userid)
                                if err != nil {
                                    authError(c, ErrForbidden)
                                }
                                c.Request = addToContext(c, utils.ProjectContextKeys.UserCtxKey, user)
                                c.Next()

                            } else {
                                authError(c, ErrMissingExpField)
                            }
                            
                        } else{
                            authError(c, ErrInvalidAccessToken)
                        }
                    }
            	}
            }
		}
	})
}
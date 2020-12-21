package auth

import(
	"net/http"
    "time"

    "github.com/linkc0829/go-ics/internal/mongodb"

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

	return func(c *gin.Context){

		c.Request = addProviderToContext(c, c.Param("provider"))
		// try to get the user without re-authenticating
        if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
            gothic.BeginAuthHandler(c.Writer, c.Request)
        } else {
            log.Printf("user: %#v", gothUser)
        }
	}
}

// Callback callback to complete auth provider flow
func CallBack(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context){
		// You have to add value context with provider name to get provider name in GetProviderName method
        c.Request = addProviderToContext(c, c.Param("provider"))
        user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }

        u, err := db.FindUserByJWT(user.Email, user.Provider, user.UserID)
        // logger.Infof("gothUser: %#v", user)
        if err != nil {
            if u, err = orm.CreateUserFromGoth(&user); err != nil {
                log.Println("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: " + err)
                c.AbortWithError(http.StatusInternalServerError, err)
            }
        }

        log.Println("[Auth.CallBack.UserLoggedIn]: ", u.ID)

        jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.JWT.Algorithm), Claims{
            Email: user.Email,
            StandardClaims: jwt.StandardClaims{
                Id:        user.ID,
                Issuer:    user.Provider,
                IssuedAt:  time.Now().UTC().Unix(),
                NotBefore: time.Now().UTC().Unix(),
                ExpiresAt: user.ExpiresAt.UTC().Unix(),
            },
        })

        token, err := jwtToken.SignedString([]byte(cfg.JWT.Secret))
        if err != nil {
            log.Println("[Auth.Callback.JWT] error: " + err)
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }
        log.Println("token: " + token)
        json := gin.H{
            "type":          "Bearer",
            "token":         token,
            "refresh_token": user.RefreshToken,
        }
        c.JSON(http.StatusOK, json)

	}
}


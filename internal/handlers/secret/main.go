package secret

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-ics/internal/mongodb"
	"github.com/linkc0829/go-ics/internal/mongodb/models"
	"github.com/linkc0829/go-ics/pkg/utils"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claims JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type RefClaims struct {
	ID string `json:"_id"`
	jwt.StandardClaims
}

func SignupHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//check if user exists
		email := c.Request.FormValue("email")
		userID := c.Request.FormValue("userID")
		nickname := c.Request.FormValue("nickname")
		provider := "ics"
		password := c.Request.FormValue("password")

		_, err := db.FindUserByJWT(email, provider, userID)
		if err == nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("ics signup: user exists"))
		}
		//encript password
		password, err = encriptPassword(password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("ics signup: password encripted failed, server error."))
		}

		newUser := &models.UserModel{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Password:  password,
			Email:     email,
			NickName:  nickname,
			CreatedAt: time.Now(),
			LastQuery: time.Now(),
			Provider:  provider,
		}

		//create access token and refresh token
		token, tokenExpiry, refreshToken, err := CreateTokenPair(cfg, newUser)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		//add to db
		newUser.RefreshToken = refreshToken

		//insert to db
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Users.InsertOne(ctx, newUser)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		//set token
		json := gin.H{
			"type":         "Bearer",
			"token":        token,
			"token_expiry": tokenExpiry,
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, json)
	}

}

func LoginHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Request.FormValue("email")
		userID := c.Request.FormValue("userID")
		provider := "ics"
		password := c.Request.FormValue("password")

		user, err := db.FindUserByJWT(email, provider, userID)

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		log.Println(user)

		if !checkPassword(user.Password, password) {
			c.AbortWithError(http.StatusUnauthorized, errors.New("password incorrect"))
		}

		//create access token and refresh token
		token, tokenExpiry, refreshToken, err := CreateTokenPair(cfg, user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		//update DB
		user.RefreshToken = refreshToken
		q := bson.M{"_id": user.ID}
		update := bson.M{"$set": user}
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Users.UpdateOne(ctx, q, update)

		//set token
		json := gin.H{
			"type":         "Bearer",
			"token":        token,
			"token_expiry": tokenExpiry,
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, json)

	}
}

//RefreshTokenHandler will verify refresh_token is valid or not, then issue new Tokens if valid
func RefreshTokenHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring, err := c.Cookie("refresh_token")
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("cannot find refresh token string in cookie"))
		}

		key := []byte(cfg.JWT.Secret)
		token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod(cfg.JWT.Algorithm) != t.Method {
				return nil, errors.New("invalid signing algorithm")
			}
			return key, nil
		})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}

		//check refresh token is expired or not
		if token.Valid == false {
			c.AbortWithError(http.StatusUnauthorized, errors.New("token invalid"))
		}
		claims := token.Claims.(jwt.MapClaims)
		ID, err := primitive.ObjectIDFromHex(claims["_id"].(string))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("invalid object id"))
		}

		//check if user exists
		var result models.UserModel
		q := bson.M{"_id": ID}
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		if err = db.Users.FindOne(ctx, q).Decode(&result); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		//check token match or not
		if result.RefreshToken != tokenstring {
			c.AbortWithError(http.StatusInternalServerError, errors.New("tokens doesn't match"))
		}

		//generate new token pair
		accToken, tokenExpiry, refreshToken, err := CreateTokenPair(cfg, &result)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		//update DB
		result.RefreshToken = refreshToken
		update := bson.M{"$set": result}
		ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Users.UpdateOne(ctx, q, update)

		//set token
		json := gin.H{
			"type":         "Bearer",
			"token":        accToken,
			"token_expiry": tokenExpiry,
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, json)

	}
}

func encriptPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(bytes), err
}

func checkPassword(pwdHash string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}

func CreateTokenPair(conf *utils.ServerConfig, user *models.UserModel) (string, time.Time, string, error) {
	accExp, _ := time.ParseDuration(conf.JWT.AccessTokenExpire)
	accExpireAt := time.Now().Add(accExp).UTC()

	claims := Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			Id:        user.UserID,
			Issuer:    user.Provider,
			IssuedAt:  time.Now().UTC().Unix(),
			NotBefore: time.Now().UTC().Unix(),
			ExpiresAt: accExpireAt.Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(conf.JWT.Algorithm), claims)

	accToken, err := jwtToken.SignedString([]byte(conf.JWT.Secret))
	if err != nil {
		log.Println("ICS Auth error: " + err.Error())
		return "", time.Now(), "", err
	}

	//RefreshToken, https://bit.ly/3r7753B
	refExp, _ := time.ParseDuration(conf.JWT.RefreshTokenExpire)
	rToken := jwt.NewWithClaims(jwt.GetSigningMethod(conf.JWT.Algorithm), RefClaims{
		ID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refExp).UTC().Unix(),
		},
	})

	refToken, err := rToken.SignedString([]byte(conf.JWT.Secret))
	if err != nil {
		log.Println("ICS Auth error: " + err.Error())
		return "", time.Now(), "", err
	}
	return accToken, accExpireAt, refToken, nil
}

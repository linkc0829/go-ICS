package secret

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"
	"github.com/linkc0829/go-icsharing/pkg/utils/datasource"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claims JWT claims
type Claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type RefClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

//ErrorWriter set redirect header to index and show error message
func ErrorWriter(c *gin.Context, code int, err error) {
	//c.Writer.Header().Set("Location", "/")
	err = errors.New("[ics secret] error: " + err.Error())
	c.Error(err)
	data := struct {
		Title        string
		ErrorMessage string
	}{
		Title:        http.StatusText(code),
		ErrorMessage: err.Error(),
	}
	c.HTML(code, "layout", data)
}

func SignupHandler(cfg *utils.ServerConfig, db *datasource.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//check if user exists
		email := c.Request.FormValue("email")
		userID := c.Request.FormValue("userID")
		nickname := c.Request.FormValue("nickname")
		provider := "ics"
		password := c.Request.FormValue("password")

		_, err := db.Mongo.FindUserByJWT(email, provider, userID)
		if err == nil {
			ErrorWriter(c, http.StatusBadRequest, errors.New("ics signup: user exists"))
			return
		}
		//encript password
		password, err = encriptPassword(password)
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, errors.New("ics signup: password encripted failed, server error."))
			return
		}

		newUser := &models.UserModel{
			ID:              primitive.NewObjectID(),
			UserID:          userID,
			Password:        &password,
			Email:           email,
			NickName:        &nickname,
			CreatedAt:       time.Now(),
			LastIncomeQuery: time.Now(),
			LastCostQuery:   time.Now(),
			Provider:        provider,
			Role:            models.USER,
		}

		//create access token and refresh token
		token, tokenExpiry, refreshToken, err := CreateTokenPair(cfg, newUser.ID.Hex(), provider)
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}

		//add to redis
		refExp, _ := strconv.Atoi(cfg.JWT.RefreshTokenExpire[:len(cfg.JWT.RefreshTokenExpire)-1])
		db.Redis.Do("Set", newUser.ID.Hex(), refreshToken, "EX", refExp*3600)

		//insert to db
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Mongo.Users.InsertOne(ctx, newUser)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		//set token
		data := struct {
			TokenType   string
			Token       string
			TokenExpiry time.Time
			Redirect    string
		}{
			Token:       token,
			TokenType:   "Bearer",
			TokenExpiry: tokenExpiry,
			Redirect:    "/profile/" + newUser.ID.Hex(),
		}
		security := false
		if cfg.URISchema == "https://" {
			security = true
		}
		domain := cfg.CookiesDomain
		c.SetCookie("refresh_token", refreshToken, 0, "/", domain, security, true)
		c.HTML(http.StatusOK, "callback", data)
	}

}

func LoginHandler(cfg *utils.ServerConfig, db *datasource.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Request.FormValue("email")
		userID := c.Request.FormValue("userID")
		provider := "ics"
		password := c.Request.FormValue("password")

		user, err := db.Mongo.FindUserByJWT(email, provider, userID)

		if err != nil {
			err = errors.New("Login failed: " + err.Error())
			ErrorWriter(c, http.StatusUnauthorized, err)
			return
		}

		if !checkPassword(*user.Password, password) {
			err = errors.New("Login failed: password incorrect")
			ErrorWriter(c, http.StatusUnauthorized, err)
			return
		}

		//create access token and refresh token
		token, tokenExpiry, refreshToken, err := CreateTokenPair(cfg, user.ID.Hex(), provider)
		if err != nil {
			err = errors.New("Login failed: " + err.Error())
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}

		//update redis
		refExp, _ := strconv.Atoi(cfg.JWT.RefreshTokenExpire[:len(cfg.JWT.RefreshTokenExpire)-1])
		db.Redis.Do("Set", user.ID.Hex(), refreshToken, "EX", refExp*3600)

		//set token
		data := struct {
			TokenType   string
			Token       string
			TokenExpiry time.Time
			Redirect    string
		}{
			Token:       token,
			TokenType:   "Bearer",
			TokenExpiry: tokenExpiry,
			Redirect:    "/profile/" + user.ID.Hex(),
		}
		security := false
		if cfg.URISchema == "https://" {
			security = true
		}
		domain := cfg.CookiesDomain
		c.SetCookie("refresh_token", refreshToken, 0, "/", domain, security, true)
		c.HTML(http.StatusOK, "callback", data)

	}
}

//RefreshTokenHandler will verify refresh_token is valid or not, then issue new Tokens if valid
func RefreshTokenHandler(cfg *utils.ServerConfig, db *datasource.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring, err := c.Cookie("refresh_token")
		if err != nil {
			ErrorWriter(c, http.StatusUnauthorized, errors.New("cannot find refresh token string in cookie"))
			return
		}

		key := []byte(cfg.JWT.Secret)
		token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod(cfg.JWT.Algorithm) != t.Method {
				return nil, errors.New("invalid signing algorithm")
			}
			return key, nil
		})

		if err != nil {
			ErrorWriter(c, http.StatusUnauthorized, err)
			return
		}
		//check refresh token is expired or not
		if token.Valid == false {
			ErrorWriter(c, http.StatusUnauthorized, errors.New("token invalid"))
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		ID := claims["id"].(string)
		issuer := claims["iss"].(string)
		if err != nil {
			ErrorWriter(c, http.StatusUnauthorized, errors.New("invalid object id"))
			return
		}
		//check redis for token
		tokenFromRedis, err := redis.String(db.Redis.Do("GET", ID))
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}

		//check token match or not
		if tokenFromRedis != tokenstring {
			ErrorWriter(c, http.StatusInternalServerError, errors.New("tokens doesn't match"))
			return
		}

		//generate new token pair
		accToken, tokenExpiry, refreshToken, err := CreateTokenPair(cfg, ID, issuer)
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}

		//soft refresh
		refExp, _ := strconv.Atoi(cfg.JWT.RefreshTokenExpire[:len(cfg.JWT.RefreshTokenExpire)-1])
		db.Redis.Do("Set", ID, refreshToken, "EX", refExp*3600)

		//set token
		json := gin.H{
			"type":         "Bearer",
			"token":        accToken,
			"token_expiry": tokenExpiry.Local(),
		}
		security := false
		if cfg.URISchema == "https://" {
			security = true
		}
		domain := cfg.CookiesDomain
		c.SetCookie("refresh_token", refreshToken, 0, "/", domain, security, true)
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

func CreateTokenPair(conf *utils.ServerConfig, id string, issuer string) (string, time.Time, string, error) {
	accExp, _ := time.ParseDuration(conf.JWT.AccessTokenExpire)
	accExpireAt := time.Now().Add(accExp).UTC()

	claims := Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			Id:        primitive.NewObjectID().Hex(),
			Issuer:    issuer,
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
		ID: id,
		StandardClaims: jwt.StandardClaims{
			Id:        primitive.NewObjectID().Hex(),
			Issuer:    issuer,
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

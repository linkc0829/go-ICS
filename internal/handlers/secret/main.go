package secret

import(
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/linkc0829/go-ics/internal/mongodb"
	"github.com/linkc0829/go-ics/internal/mongodb/models"
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


func SignupHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc{
	return func(c *gin.Context){
		//check if user exists
		email := c.Request.FormValue("email")
		userId := c.Request.FormValue("userId")
		nickname := c.Request.FormValue("nickname")
		provider := "ics"
		password := c.Request.FormValue("password")

		if u, err := db.FindUserByJWT(email, provider, userId); err == nil {
			c.AbortWithError(http.StatusBadRequest, err.New("ics signup: user exists"))
		}
		//encript password
		password, err = encriptPassword(password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err.New("ics signup: password encripted failed, server error."))
		}

		//create access token and refresh token
		token, token_expiry, refreshToken, err := createTokenPair(cfg.JWT, result)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err.Error())
		}

		//add to db
		newUser := &models.UserModel{
			ID:				primitive.NewObjectID(),
			UserId: 		userId,
			Password:		password,
			Email:			email,
			NickName:		nickname
			CreatedAt:		time.Now(),
			LastQuery:		time.Now(),
			Provider:		provider,
			RefreshToken:	refreshToken,
		}

		//insert to db
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Users.InsertOne(ctx, newUser)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err.Error())
		}

		//set token
		json := gin.H{
			"type":          	"Bearer",
			"token":         	token,
			"token_expiry": 	token_expiry,
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, json)
	}
	
}

func LoginHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc{
	return func(c *gin.Context){
		email := c.Request.FormValue("email")
		userId := c.Request.FormValue("userId")
		provider := "ics"
		password := c.Request.FormValue("password")

		user, err := db.users.FindUserByJWT(email, provider, userId)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err.Error())
		}

		if checkPassword(user.Password, password) == false {
			c.AbortWithError(http.StatusUnauthorized, error.New("password incorrect"))
		}

		//create access token and refresh token
		token, token_expiry, refreshToken, err := createTokenPair(cfg.JWT, result)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err.Error())
		}

		//update DB
		user.RefreshToken = refreshToken
		q := bson.M{"_id", user.ID}
		update := bson.M{"$set", user}
		ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Users.UpdateOne(ctx, q, update)

		//set token
		json := gin.H{
			"type":          	"Bearer",
			"token":         	token,
			"token_expiry": 	token_expiry,
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, json)

	}
}



//RefreshTokenHandler will verify refresh_token is valid or not, then issue new Tokens if valid
func RefreshTokenHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc{
	return func(c *gin.Context){
		tokenstring := c.Query("refresh_token")

		token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error){
	    	if jwt.GetSigningMethod(cfg.JWT.Algorithm) != t.Method{
	    		return nil, errors.New("invalid signing algorithm")
	    	}
	    	return key, nil
	    })

	    if err != nil {
	    	c.AbortWithError(http.StatusUnauthorized, err.Error())
	    }

	    //check refresh token is expired or not
	    if token.Valid == false {
	    	c.AbortWithError(http.StatusUnauthorized, err.Error())
	    }
	    claims := token.Claims.(jwt.MapClaims)
	    ID := claims["_id"]

	    //check if user exists
	    var result models.UserModel
	    q := bson.M{"_id": ID}
	    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		if err = db.Users.FindOne(ctx, q).Decode(&result); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err.Error())
		}

		//check token match or not
		if result.RefreshToken != tokenstring {
			c.AbortWithError(http.StatusInternalServerError, err.New("tokens doesn't match"))
		}

		//generate new token pair
		token, token_expiry, refreshToken, err := createTokenPair(cfg.JWT, result)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err.Error())
		}

		//update DB
		result.RefreshToken = refreshToken
		update := bson.M{"$set", result}
		ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
		_, err = db.Users.UpdateOne(ctx, q, update)

		//set token
		json := gin.H{
			"type":          	"Bearer",
			"token":         	token,
			"token_expiry": 	token_expiry,
		}
		c.SetCookie("refresh_token", refreshToken, 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, json)


	}
}

func encriptPassword(pwd string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(bytes), err
}

func checkPassword(pwd_hash string, pwd string) (bool, error){
	err := bcrypt.CompareHashAndPassword(pwd_hash, pwd)
	if err != nil{
		return err
	}
	return nil
}

func createTokenPair(JWT *config.JWT, user *models.UserModel) (string, time, string, error){
		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(JWT.Algorithm), Claims{
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				Id:        user.UserId,
				Issuer:    user.Provider,
				IssuedAt:  time.Now().UTC().Unix(),
				NotBefore: time.Now().UTC().Unix(),
				ExpiresAt: time.Now().Add(JWT.AccessTokenExpire).UTC().Unix(),
			},
		})
		accToken, err := jwtToken.SignedString([]byte(JWT.Secret))
		if err != nil {
			log.Println("ICS Auth error: " + err)
			return "", err
		}
		accExpire := jwtToken.StandardClaims.ExpiresAt

		//RefreshToken, https://bit.ly/3r7753B
		rToken := jwt.NewWithClaims(jwt.GetSigningMethod(JWT.Algorithm), RefClaims{
			ID: user.ID,
			StandardClaims: jwt.StandardClaims{
				ID:			user.ID
				ExpiresAt: time.Now().Add(JWT.AccessTokenExpire).UTC().Unix(),
			},
		})

		refToken, err := rToken.SignedString([]byte(JWT.Secret))
		if err != nil {
			log.Println("ICS Auth error: " + err)
			return "", err
		}
		return accToken, accExpire, refToken, nil
}
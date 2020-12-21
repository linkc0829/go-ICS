package mongodb

import(
	"context"
	"log"

	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/linkc0829/go-ics/internal/mongodb/models"

)

//FindUserByAPIKey will find user related to the APIKey
func (db *MongoDB) FindUserByAPIKey(apiKey String) (*models.UserModel, error){

	var results UserModel

	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}

	q := bson.M{"apiKey": apiKey}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err := db.Users.FindOne(ctx, q).Decode(&results); err != nil {
		return nil, err
	}
	return &results, nil
}

//FindUserByJWTToken will find user related to the JWT Token
func (db *MongoDB) FindUserByJWT(email string, provider string, userId string) (*models.UserModel, error){
	var results UserModel

	q := bson.M{"email": email, "provider": provider, "userId": userId}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err := db.Users.FindOne(ctx, q).Decode(&results); err != nil {
		return nil, err
	}
	return &results, nil
}

//CreateUserFromGoth using data from goth to create user
func (db *MongoDB) CreateUserFromGoth(user *goth.User) (*model.UserModel, error){

	if user.UserID == "" || user.Email == "" {
		return nil, error.New("Goth user field incomplete, UserID and Email are necessary.")
	}

	newUser := &UserModel{
		ID:			primitive.NewObjectID(),
		UserId: 	user.UserID,
		Email:		user.Email,
		NickName:	user.NickName,
		CreatedAt:	time.Now(),
		LastQuery:	time.Now(),
		Provider:	user.Provider,
		AvatarURL:	user.AvatarURL,
	}

	//insert to db
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err = db.Users.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/markbates/goth"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	//"github.com/linkc0829/go-ics/pkg/utils"
	"github.com/linkc0829/go-ics/internal/mongodb/models"
)

//FindUserByAPIKey will find user related to the APIKey
func (db *MongoDB) FindUserByAPIKey(apiKey string) (*models.UserModel, error) {

	var results models.UserModel

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

//FindUserByJWT will find user related to the JWT Token
func (db *MongoDB) FindUserByJWT(email string, provider string, userID string) (*models.UserModel, error) {
	var results models.UserModel

	q := bson.M{"email": email, "provider": provider, "userid": userID}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if err := db.Users.FindOne(ctx, q).Decode(&results); err != nil {
		return nil, err
	}
	return &results, nil
}

//CreateUserFromGoth using data from goth to create user
func (db *MongoDB) CreateUserFromGoth(user *goth.User) (*models.UserModel, error) {

	if user.UserID == "" || user.Email == "" {
		return nil, errors.New("Goth user field incomplete, UserID and Email are necessary.")
	}

	newUser := &models.UserModel{
		ID:              primitive.NewObjectID(),
		UserID:          user.UserID,
		Email:           user.Email,
		NickName:        &user.NickName,
		CreatedAt:       time.Now(),
		LastIncomeQuery: time.Now(),
		LastCostQuery:   time.Now(),
		Provider:        user.Provider,
		AvatarURL:       user.AvatarURL,
		Role:            models.USER,
	}

	//insert to db
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err := db.Users.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

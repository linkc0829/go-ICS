package resolvers

import (
	"context"
	"time"
	"fmt"

	"github.com/linkc0829/go-ics/internal/graph/models"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
	"github.com/linkc0829/go-ics/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if result, err := getUserById(ctx, r.DB, userID); err != nil{
		return nil, err
	}else{
		return result, nil
	}
}

func (r *queryResolver) MyFriends(ctx context.Context) (*models.Users, error) {
	panic("not implemented")
}

func (r *queryResolver) MyFollowers(ctx context.Context) (*models.Users, error) {
	panic("not implemented")
}


func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	
	//check if userId exists
	result, err := getUserById(ctx, r.DB, *input.UserID)
	if err == nil{
		return nil, fmt.Errorf("Create user failed. UserID already exists.")
	}

	//check if email exists
	q := bson.M{"email": input.Email}
	if err := r.DB.Users.FindOne(ctx, q).Decode(&result); err == nil{
		return nil, fmt.Errorf("Create user failed. User email already exists.")
	}

	//collect current data
	UserID := *input.UserID
	Email := *input.Email
	NickName := *input.NickName
	CreatedAt := time.Now()
	LastQuery := time.Now()

	newUser := dbModel.UserModel{
		ID:			primitive.NewObjectID(),
		UserId: 	UserID,
		Email:		Email,
		NickName:	NickName,
		CreatedAt:	CreatedAt,
		LastQuery:	LastQuery,
	}

	//insert to db
	_, err = r.DB.Users.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	//retrun graph result to server
	result, _ = getUserById(ctx, r.DB, UserID)

	return result, nil

}

func (r *mutationResolver) UpdateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddFriends(ctx context.Context, userID string) (*models.User, error) {
	panic("not implemented")
}

//helper functions

func getUserById(ctx context.Context, DB *mongodb.MongoDB, userID string) (*models.User, error) {

	q := bson.M{"userid": userID}
	result := dbModel.UserModel{}
	
	if err := DB.Users.FindOne(ctx, q).Decode(&result); err != nil{
		return nil, fmt.Errorf("UserID doesn't exist.")
	}

	r := &models.User{
		ID:       	result.ID.Hex(),
		Email:     	result.Email,
		UserID:   	result.UserId,
		NickName:	&result.NickName,
		CreatedAt: 	result.CreatedAt,
		APIKey:		&result.APIKey,
		//Friends:	
		//Followers:	
		LastQuery:	result.LastQuery,
	}
	
	return r, nil
}
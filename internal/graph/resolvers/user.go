package resolvers

import (
	"context"
	"time"
	"fmt"

	"github.com/linkc0829/go-ics/internal/graph/models"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if result, err := getUserById(ctx, UserID); err != nil{
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
	result, err := getUserById(input.UserID)
	if err == nil{
		return nil, fmt.Errorf("Create user failed. UserID already exists.")
	}

	//check if email exists
	q = bson.M{"email": input.Email}
	if err := r.db.users.FindOne(ctx, q).Decode(&result); err == nil{
		return nil, fmt.Errorf("Create user failed. User email already exists.")
	}

	//collect current data
	UserID := input.UserID
	Email := input.Email
	NickName := input.NickName
	CreateAt := time.Now()
	LastQuery := time.Now()

	result = dbModel.UserModel{
		UserID: 	UserID,
		Email:		Email,
		NickName:	NickName,
		CreateAt:	CreateAt,
		LastQuery:	LastQuery,
	}

	//insert to db
	_, err := r.db.users.InsertOne(ctx, result)
	if err != nil {
		return nil, err
	}

	//retrun graph result to server
	result, _ = getUserById(ctx, UserID)

	return result, nil

}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddFriends(ctx context.Context, userID string) (*models.User, error) {
	panic("not implemented")
}

//helper functions

func getUserById(ctx context.Context, userID string) (*models.User, error) {

	q := bson.M{"userid": input.UserID}
	result := dbModel.UserModel{}
	
	if err := r.db.users.FindOne(ctx, q).Decode(&result); err != nil{
		return nil, fmt.Errorf("UserID doesn't exist.")
	}
	return result, nil
}
package resolvers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/internal/mongodb"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
	"github.com/linkc0829/go-ics/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)

	result, err := getUserByID(ctx, r.DB, me.ID.Hex())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *queryResolver) GetUser(ctx context.Context, ID string) (*models.User, error) {
	if result, err := getUserByID(ctx, r.DB, ID); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (r *queryResolver) MyFriends(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) MyFollowers(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {

	//check if user exists
	var result dbModel.UserModel
	q := bson.M{"email": input.Email, "provider": "ics", "userId": input.UserID}
	if err := r.DB.Users.FindOne(ctx, q).Decode(&result); err == nil {
		return nil, fmt.Errorf("Create user failed. User already exists.")
	}

	//collect current data
	UserID := *input.UserID
	Email := *input.Email
	NickName := input.NickName
	CreatedAt := time.Now()
	LastQuery := time.Now()

	newUser := &dbModel.UserModel{
		ID:        primitive.NewObjectID(),
		UserID:    UserID,
		Email:     Email,
		NickName:  NickName,
		CreatedAt: CreatedAt,
		LastQuery: LastQuery,
		Provider:  "ics",
		Friends:   nil,
	}

	//insert to db
	_, err := r.DB.Users.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	ret := &models.User{
		ID:        newUser.ID.Hex(),
		UserID:    newUser.UserID,
		Email:     newUser.Email,
		NickName:  newUser.NickName,
		Friends:   nil,
		Followers: nil,
	}

	return ret, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input models.UserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

//AddFriends add id to my Friends
func (r *mutationResolver) AddFriend(ctx context.Context, id string) (*models.User, error) {

	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)

	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		//not a valid objectID
		return nil, err
	}

	//add to friend
	q := bson.M{"_id": me.ID}
	result := dbModel.UserModel{}
	if err := r.DB.Users.FindOne(ctx, q).Decode(&result); err != nil {
		return nil, fmt.Errorf("UserID doesn't exist.")
	}
	result.Friends = append(result.Friends, hexID)

	//update DB
	upd := bson.M{"$set": result}
	_, err = r.DB.Users.UpdateOne(ctx, q, upd)
	if err != nil {
		return nil, err
	}

	//return user
	update, err := getUserByID(ctx, r.DB, me.ID.Hex())
	if err != nil {
		return nil, err
	}

	return update, nil
}

func (r *mutationResolver) AddFollower(ctx context.Context, id string) (*models.User, error) {
	panic("not implemented")
}

//helper functions

func getUserByID(ctx context.Context, DB *mongodb.MongoDB, ID string) (*models.User, error) {

	hexID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		//not a valid objectID
		return nil, err
	}

	q := bson.M{"_id": hexID}
	result := dbModel.UserModel{}

	if err := DB.Users.FindOne(ctx, q).Decode(&result); err != nil {
		return nil, fmt.Errorf("UserID doesn't exist.")
	}

	friends, err := getUserFriends(ctx, DB, &result)
	if err != nil {
		return nil, err
	}
	followers, err := getUserFollowers(ctx, DB, &result)
	if err != nil {
		return nil, err
	}

	r := &models.User{
		ID:        result.ID.Hex(),
		Email:     result.Email,
		UserID:    result.UserID,
		NickName:  result.NickName,
		CreatedAt: result.CreatedAt,
		Friends:   friends,
		Followers: followers,
	}

	return r, nil
}

func getUserFriends(ctx context.Context, DB *mongodb.MongoDB, user *dbModel.UserModel) (friends []*models.User, err error) {

	for _, f_id := range user.Friends {
		f, err := getUserByID(ctx, DB, f_id.Hex())
		if err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return
}

func getUserFollowers(ctx context.Context, DB *mongodb.MongoDB, me *dbModel.UserModel) (followers []*models.User, err error) {
	//find users that have me as friend
	q := bson.M{"friends": me.ID}
	cursor, err := DB.Users.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		f := dbModel.UserModel{}
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &f)
		follower, err := getUserByID(ctx, DB, f.ID.Hex())
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	return
}

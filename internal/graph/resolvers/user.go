package resolvers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/linkc0829/go-icsharing/internal/db/mongodb"
	dbModel "github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/internal/graph/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	result, err := r.resolveUsers(ctx, me.ID.Hex())
	return result[0], err
}

func (r *queryResolver) GetUser(ctx context.Context, ID string) (*models.User, error) {
	if result, err := getUserByID(ctx, r.DB, ID); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (r *queryResolver) MyFriends(ctx context.Context) ([]*models.User, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	friends := getUserFriends(me)
	result, err := r.resolveUsers(ctx, friends...)
	return result, err
}

func (r *queryResolver) MyFollowers(ctx context.Context) ([]*models.User, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	followers, err := getUserFollowers(ctx, r.DB, me)
	if err != nil {
		return nil, err
	}

	result, err := r.resolveUsers(ctx, followers...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	if !isAdmin(ctx) {
		return nil, errors.New("permission denied, only Admin oculd create user")
	}
	//check if user exists
	var result dbModel.UserModel
	q := bson.M{"email": input.Email, "provider": "ics", "userid": input.UserID}

	if err := r.DB.Users.FindOne(ctx, q).Decode(&result); err == nil {
		if err != mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Create user failed. User already exists.")
		} else {
			return nil, err
		}
	}

	//collect current data
	UserID := input.UserID
	Email := input.Email
	NickName := input.NickName
	CreatedAt := time.Now()
	LastIncomeQuery := time.Now()
	LastCostQuery := time.Now()

	newUser := &dbModel.UserModel{
		ID:              primitive.NewObjectID(),
		UserID:          UserID,
		Email:           Email,
		NickName:        NickName,
		CreatedAt:       CreatedAt,
		LastIncomeQuery: LastIncomeQuery,
		LastCostQuery:   LastCostQuery,
		Provider:        "ics",
		Friends:         []primitive.ObjectID{},
		Role:            dbModel.USER,
	}

	//insert to db
	_, err := r.DB.Users.InsertOne(ctx, newUser)
	if err != nil {
		log.Fatal(err)
	}

	ret := &models.User{
		ID:        newUser.ID.Hex(),
		UserID:    newUser.UserID,
		Email:     newUser.Email,
		NickName:  newUser.NickName,
		CreatedAt: newUser.CreatedAt,
		Friends:   []string{},
		Followers: []string{},
		Role:      models.RoleUser,
	}

	return ret, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input models.UpdateUserInput) (*models.User, error) {

	//check editorial permission
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	if !isAdmin(ctx) && me.ID.Hex() != id {
		return nil, errors.New("permission denied")
	}

	user, err := getUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}
	//update user profile
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.UserID != nil {
		user.UserID = *input.UserID
	}
	if input.NickName != nil {
		user.NickName = input.NickName
	}

	//check if user exists
	var result dbModel.UserModel
	q := bson.M{"email": user.Email, "provider": "ics", "userid": user.UserID}
	err = r.DB.Users.FindOne(ctx, q).Decode(&result)
	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("Update user failed. There's one user has the same info.")
	}

	//update db
	primID, _ := primitive.ObjectIDFromHex(id)
	q = bson.M{"_id": primID}
	upd := bson.M{"$set": input}
	_, err = r.DB.Users.UpdateOne(ctx, q, upd)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {

	if !isAdmin(ctx) {
		return false, errors.New("permission denied")
	}
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	q := bson.M{"_id": primID}
	result, err := r.DB.Users.DeleteOne(ctx, q)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 1 {
		return true, nil
	}
	return false, nil
}

//AddFriends add id to my Friends
func (r *mutationResolver) AddFriend(ctx context.Context, id string) (bool, error) {

	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	fID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		//not a valid objectID
		return false, err
	}
	if fID == me.ID {
		return false, errors.New("Cannot add yourself to friend")
	}
	//if already friend, remove fID
	length := len(me.Friends)
	for i, f := range me.Friends {
		if f == fID {
			if length == 1 {
				me.Friends = me.Friends[:0]
			} else {
				me.Friends[i] = me.Friends[length-1]
				me.Friends = me.Friends[:length-1]
			}
			break
		}
	}
	//previous not friend
	if length == len(me.Friends) {
		//add friend
		me.Friends = append(me.Friends, fID)
	}
	//update DB
	q := bson.M{"_id": me.ID}
	upd := bson.M{"$set": bson.M{"friends": me.Friends}}
	_, err = r.DB.Users.UpdateOne(ctx, q, upd)
	if err != nil {
		return false, err
	}
	return true, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Friends(ctx context.Context, obj *models.User) ([]*models.User, error) {
	return r.resolveUsers(ctx, obj.Friends...)
}

func (r *userResolver) Followers(ctx context.Context, obj *models.User) ([]*models.User, error) {
	return r.resolveUsers(ctx, obj.Followers...)
}

func (r *userResolver) Role(ctx context.Context, obj *models.User) (models.Role, error) {
	panic("not implemented")
}

//helper functions

func getDBUserByID(ctx context.Context, DB *mongodb.MongoDB, id string) (*dbModel.UserModel, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	qUser := bson.M{"_id": userID}
	user := dbModel.UserModel{}
	if err := DB.Users.FindOne(ctx, qUser).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserByID(ctx context.Context, DB *mongodb.MongoDB, ID string) (*models.User, error) {

	result, err := getDBUserByID(ctx, DB, ID)
	if err != nil {
		return nil, err
	}
	friends := getUserFriends(result)
	if err != nil {
		return nil, err
	}
	followers, err := getUserFollowers(ctx, DB, result)
	if err != nil {
		return nil, err
	}
	role := models.RoleUser
	if result.Role == string(models.RoleAdmin) {
		role = models.RoleAdmin
	}

	r := &models.User{
		ID:        result.ID.Hex(),
		Email:     result.Email,
		UserID:    result.UserID,
		NickName:  result.NickName,
		CreatedAt: result.CreatedAt,
		Friends:   friends,
		Followers: followers,
		Role:      role,
	}

	return r, nil
}

func getUserFriends(user *dbModel.UserModel) (friends []string) {
	for _, f_id := range user.Friends {
		f := f_id.Hex()
		friends = append(friends, f)
	}
	return
}

func getUserFollowers(ctx context.Context, DB *mongodb.MongoDB, me *dbModel.UserModel) (followers []string, err error) {
	//find users that have me as friend
	q := bson.M{"friends": me.ID}
	cursor, err := DB.Users.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		f := dbModel.UserModel{}
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &f)
		follower := f.ID.Hex()
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	return
}

//is user admin?
func isAdmin(ctx context.Context) bool {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	if me.Role == dbModel.ADMIN {
		return true
	}
	return false
}

//could userA see userB's friend content?
//or did userB add userA to friend?
func couldViewFriendContent(userA *dbModel.UserModel, userB *dbModel.UserModel) bool {
	for _, f := range userB.Friends {
		if userA.ID == f {
			return true
		}
	}
	return false
}

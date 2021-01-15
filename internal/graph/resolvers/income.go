package resolvers

import (
	"context"
	"errors"

	dbModel "github.com/linkc0829/go-ics/internal/db/mongodb/models"
	"github.com/linkc0829/go-ics/internal/graph/models"
	tf "github.com/linkc0829/go-ics/internal/graph/resolvers/transformer"
	"github.com/linkc0829/go-ics/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateIncome(ctx context.Context, input models.CreateIncomeInput) (*models.Income, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)

	cat := (models.PortfolioCategory)(input.Category)
	newIncome := dbModel.IncomeModel{
		ID:          primitive.NewObjectID(),
		Owner:       me.ID,
		Amount:      input.Amount,
		OccurDate:   input.OccurDate,
		Description: input.Description,
		Vote:        nil,
		Category:    cat,
		Privacy:     input.Privacy,
	}
	//insert to db
	_, err := r.DB.Income.InsertOne(ctx, newIncome)
	if err != nil {
		return nil, err
	}

	result := tf.DBPortfolioToGQLPortfolio(newIncome).(models.Income)

	return &result, nil
}

func (r *mutationResolver) UpdateIncome(ctx context.Context, id string, input models.UpdateIncomeInput) (*models.Income, error) {
	incomeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	q := bson.M{"_id": incomeID}
	result := dbModel.IncomeModel{}
	if err := r.DB.Income.FindOne(ctx, q).Decode(&result); err != nil {
		return nil, err
	}

	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	if !isAdmin(ctx) && me.ID != result.Owner {
		return nil, errors.New("permission denied")
	}

	if input.Amount != nil {
		result.Amount = *input.Amount
	}
	if input.Category != nil {
		cat := (models.PortfolioCategory)(*input.Category)
		result.Category = cat
	}
	if input.Description != nil {
		result.Description = *input.Description
	}
	if input.OccurDate != nil {
		result.OccurDate = *input.OccurDate
	}
	if input.Privacy != nil {
		result.Privacy = *input.Privacy
	}

	upd := bson.M{"$set": input}
	_, err = r.DB.Income.UpdateOne(ctx, q, upd)
	if err != nil {
		return nil, err
	}

	ret := tf.DBPortfolioToGQLPortfolio(result).(models.Income)
	return &ret, nil
}

func (r *mutationResolver) DeleteIncome(ctx context.Context, id string) (bool, error) {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	q := bson.M{"_id": primID}
	result := dbModel.IncomeModel{}
	if err := r.DB.Income.FindOne(ctx, q).Decode(&result); err != nil {
		return false, err
	}

	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	if !isAdmin(ctx) && me.ID != result.Owner {
		return false, errors.New("permission denied")
	}
	delete, err := r.DB.Income.DeleteOne(ctx, q)
	if err != nil {
		return false, err
	}
	if delete.DeletedCount == 1 {
		return true, nil
	}
	return false, nil
}

func (r *mutationResolver) VoteIncome(ctx context.Context, id string) (int, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	incomeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	q := bson.M{"_id": incomeID}
	income := dbModel.IncomeModel{}
	if err := r.DB.Income.FindOne(ctx, q).Decode(&income); err != nil {
		return -1, err
	}
	owner := dbModel.UserModel{}
	if err := r.DB.Users.FindOne(ctx, bson.M{"_id": income.Owner}).Decode(&owner); err != nil {
		return -1, err
	}
	//deny private
	if !isAdmin(ctx) && me.ID != income.Owner && income.Privacy == models.PrivacyPrivate {
		return -1, errors.New("it's private, permission denied")
	}
	//deny non-friend
	if !isAdmin(ctx) && !couldViewFriendContent(me, &owner) && me.ID != income.Owner && income.Privacy == models.PrivacyFriend {
		return -1, errors.New("it's only for friend, permission denied")
	}
	//if already voted, revoke
	length := len(income.Vote)
	for i, v := range income.Vote {
		if v == me.ID {
			if length == 1 {
				income.Vote = income.Vote[:0]
			} else {
				income.Vote[i] = income.Vote[length-1]
				income.Vote = income.Vote[:length-1]
			}
			break
		}
	}
	if length == len(income.Vote) {
		//add to vote
		income.Vote = append(income.Vote, me.ID)
	}

	//update DB
	upd := bson.M{"$set": bson.M{"vote": income.Vote}}
	_, err = r.DB.Income.UpdateOne(ctx, q, upd)
	if err != nil {
		return -1, err
	}
	return len(income.Vote), nil
}

type incomeResolver struct{ *Resolver }

func (r *incomeResolver) Vote(ctx context.Context, obj *models.Income) ([]*models.User, error) {
	return r.resolveUsers(ctx, obj.Vote...)
}

func (r *incomeResolver) Owner(ctx context.Context, obj *models.Income) (*models.User, error) {
	owner, err := r.resolveUsers(ctx, obj.Owner)
	if err != nil {
		return nil, err
	}
	return owner[0], nil
}

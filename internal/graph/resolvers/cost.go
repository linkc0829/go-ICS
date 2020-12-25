package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/models"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
	"github.com/linkc0829/go-ics/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	tf "github.com/linkc0829/go-ics/internal/graph/resolvers/transformer"
)

func (r *mutationResolver) CreateCost(ctx context.Context, input models.CreateCostInput) (*models.Cost, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)

	newCost := dbModel.CostModel{
		ID:          primitive.NewObjectID(),
		Owner:       me.ID,
		Amount:      input.Amount,
		OccurDate:   input.OccurDate,
		Description: input.Description,
		Vote:        nil,
		Category:    input.Category,
	}
	//insert to db
	_, err := r.DB.Cost.InsertOne(ctx, newCost)
	if err != nil {
		return nil, err
	}
	result := tf.DBPortfolioToGQLPortfolio(newCost).(models.Cost)

	return &result, nil
}

func (r *mutationResolver) UpdateCost(ctx context.Context, id string, input models.UpdateCostInput) (*models.Cost, error) {
	costID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	q := bson.M{"_id": costID}
	result := dbModel.CostModel{}
	if err := r.DB.Cost.FindOne(ctx, q).Decode(&result); err != nil {
		return nil, err
	}

	if input.Amount != nil {
		result.Amount = *input.Amount
	}
	if input.Category != nil {
		result.Category = *input.Category
	}
	if input.Description != nil {
		result.Description = input.Description
	}
	if input.OccurDate != nil {
		result.OccurDate = *input.OccurDate
	}

	upd := bson.M{"$set": input}
	_, err = r.DB.Cost.UpdateOne(ctx, q, upd)
	if err != nil {
		return nil, err
	}

	ret := tf.DBPortfolioToGQLPortfolio(result).(models.Cost)

	return &ret, nil
}

func (r *mutationResolver) DeleteCost(ctx context.Context, id string) (bool, error) {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	q := bson.M{"_id": primID}
	result, err := r.DB.Cost.DeleteOne(ctx, q)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 1 {
		return true, nil
	}
	return false, nil
}

func (r *mutationResolver) VoteCost(ctx context.Context, id string) (int, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	costID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	q := bson.M{"_id": costID}
	cost := dbModel.CostModel{}
	if err := r.DB.Cost.FindOne(ctx, q).Decode(&cost); err != nil {
		return -1, err
	}

	//if already voted, revoke
	length := len(cost.Vote)
	for i, v := range cost.Vote {
		if v == me.ID {
			if length == 1 {
				cost.Vote = cost.Vote[:0]
			} else {
				cost.Vote[i] = cost.Vote[length-1]
				cost.Vote = cost.Vote[:length-1]
			}
			break
		}
	}
	if length == len(cost.Vote) {
		//add to vote
		cost.Vote = append(cost.Vote, me.ID)
	}

	//update DB
	upd := bson.M{"$set": bson.M{"vote": cost.Vote}}
	_, err = r.DB.Cost.UpdateOne(ctx, q, upd)
	if err != nil {
		return -1, err
	}
	return len(cost.Vote), nil
}

type costResolver struct{ *Resolver }

func (r *costResolver) Vote(ctx context.Context, obj *models.Cost) ([]*models.User, error) {
	return r.resolveUsers(ctx, obj.Vote...)
}

func (r *costResolver) Owner(ctx context.Context, obj *models.Cost) (*models.User, error) {
	owner, err := r.resolveUsers(ctx, obj.Owner)
	if err != nil {
		return nil, err
	}
	return owner[0], nil
}

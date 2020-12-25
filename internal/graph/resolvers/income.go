package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/models"
	tf "github.com/linkc0829/go-ics/internal/graph/resolvers/transformer"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
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

	if input.Amount != nil {
		result.Amount = *input.Amount
	}
	if input.Category != nil {
		cat := (*models.PortfolioCategory)(input.Category)
		result.Category = *cat
	}
	if input.Description != nil {
		result.Description = input.Description
	}
	if input.OccurDate != nil {
		result.OccurDate = *input.OccurDate
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
	result, err := r.DB.Income.DeleteOne(ctx, q)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 1 {
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

func (r *incomeResolver) Category(ctx context.Context, obj *models.Income) (models.PortfolioCategory, error) {
	panic("not implemented")
}

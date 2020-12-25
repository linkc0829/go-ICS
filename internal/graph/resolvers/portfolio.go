package resolvers

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/linkc0829/go-ics/internal/graph/models"
	_ "github.com/linkc0829/go-ics/internal/mongodb"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
	"github.com/linkc0829/go-ics/pkg/utils"

	tf "github.com/linkc0829/go-ics/internal/graph/resolvers/transformer"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) MyCostHistory(ctx context.Context, rangeArg int) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	return r.GetUserCostHistory(ctx, me.ID.Hex(), rangeArg)
}

func (r *queryResolver) MyIncomeHistory(ctx context.Context, rangeArg int) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	return r.GetUserIncomeHistory(ctx, me.ID.Hex(), rangeArg)
}

func (r *queryResolver) MyIncome(ctx context.Context) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	return r.GetUserIncome(ctx, me.ID.Hex())
}

func (r *queryResolver) MyCost(ctx context.Context) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	return r.GetUserCost(ctx, me.ID.Hex())
}

//GetUserCost returns user id's cost
func (r *queryResolver) GetUserCost(ctx context.Context, id string) ([]models.Portfolio, error) {
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	q := bson.M{"owner": user.ID}
	cursor, err := r.DB.Cost.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	return decodeFindResult(results), nil
}

func (r *queryResolver) GetUserIncome(ctx context.Context, id string) ([]models.Portfolio, error) {
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	q := bson.M{"owner": user.ID}
	cursor, err := r.DB.Income.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	return decodeFindResult(results), nil
}

func (r *queryResolver) GetUserIncomeHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {

	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	fromDate, toDate := getDateRange(rangeArg)
	q := bson.M{"owner": user.ID, "OccurDate": bson.M{
		"$gte": fromDate,
		"$lte": toDate,
	}}
	cursor, err := r.DB.IncomeHistory.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	return decodeFindResult(results), nil
}

func (r *queryResolver) GetUserCostHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	fromDate, toDate := getDateRange(rangeArg)
	q := bson.M{"owner": user.ID, "OccurDate": bson.M{
		"$gte": fromDate,
		"$lte": toDate,
	}}
	cursor, err := r.DB.CostHistory.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	return decodeFindResult(results), nil
}

//helper function

func getDateRange(days int) (time.Time, time.Time) {

	duration, _ := time.ParseDuration(strconv.Itoa(-days*24-24) + "h")
	preDay, _ := time.ParseDuration(strconv.Itoa(-24) + "h")

	to := time.Now().Add(preDay)
	toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.Local)

	from := time.Now().Add(duration)
	fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.Local)

	return toDate, fromDate
}

func decodeFindResult(results []bson.M) []models.Portfolio {
	hist := []models.Portfolio{}
	for _, result := range results {
		c := dbModel.CostModel{}
		//encode mongodb result to JSON format, then decode
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &c)

		//convert DB User struct to GQL User struct
		gqlCost := tf.DBPortfolioToGQLPortfolio(c)
		hist = append(hist, gqlCost)
	}
	return hist
}

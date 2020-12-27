package resolvers

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/linkc0829/go-ics/internal/graph/models"
	_ "github.com/linkc0829/go-ics/internal/mongodb"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
	"github.com/linkc0829/go-ics/pkg/utils"

	tf "github.com/linkc0829/go-ics/internal/graph/resolvers/transformer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	findopt := options.Find().SetSort(bson.M{"occurDate": -1})
	cursor, err := r.DB.Cost.Find(ctx, q, findopt)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	qy, qm, qd := user.LastCostQuery.Date()
	yy, mm, dd := time.Now().Date()
	//if today has quired, return
	if (qy == yy) && (qm == mm) && (qd == dd) {
		return decodeFindResult(results), nil
	}

	//move expired cost to history if user hasn't queried today
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, time.Local)
	rets := []models.Portfolio{}
	for _, result := range results {
		c := dbModel.CostModel{}
		//encode mongodb result to JSON format, then decode
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &c)
		if c.OccurDate.Before(today) {
			//insert to history
			_, err := r.DB.CostHistory.InsertOne(ctx, c)
			if err != nil {
				return nil, err
			}
			//delete old entry
			result, err := r.DB.Cost.DeleteOne(ctx, bson.M{"_id": c.ID})
			if err != nil {
				return nil, err
			}
			if result.DeletedCount != 1 {
				return nil, errors.New("delete failed when move old cost entry")
			}
		} else {
			rets = append(rets, tf.DBPortfolioToGQLPortfolio(c))
		}
	}
	//update user LastCostQuery
	q = bson.M{"_id": user.ID}
	upd := bson.M{"$set": bson.M{"lastCostQuery": time.Now()}}
	_, err = r.DB.Users.UpdateOne(ctx, q, upd)
	if err != nil {
		return nil, err
	}
	return rets, nil
}

func (r *queryResolver) GetUserIncome(ctx context.Context, id string) ([]models.Portfolio, error) {
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	q := bson.M{"owner": user.ID}
	findopt := options.Find().SetSort(bson.M{"occurDate": -1})
	cursor, err := r.DB.Income.Find(ctx, q, findopt)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	qy, qm, qd := user.LastIncomeQuery.Date()
	yy, mm, dd := time.Now().Date()

	//if today has quired, return
	if (qy == yy) && (qm == mm) && (qd == dd) {
		return decodeFindResult(results), nil
	}

	//move expired cost to history if user hasn't queried today
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, time.Local)
	rets := []models.Portfolio{}
	for _, result := range results {
		c := dbModel.CostModel{}
		//encode mongodb result to JSON format, then decode
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &c)
		if c.OccurDate.Before(today) {
			//insert to history
			_, err := r.DB.IncomeHistory.InsertOne(ctx, c)
			if err != nil {
				return nil, err
			}
			//delete old entry
			result, err := r.DB.Income.DeleteOne(ctx, bson.M{"_id": c.ID})
			if err != nil {
				return nil, err
			}
			if result.DeletedCount != 1 {
				return nil, errors.New("delete failed when move old cost entry")
			}
		} else {
			rets = append(rets, tf.DBPortfolioToGQLPortfolio(c))
		}
	}
	//update user LastIncomeQuery
	q = bson.M{"_id": user.ID}
	upd := bson.M{"$set": bson.M{"lastIncomeQuery": time.Now()}}
	_, err = r.DB.Users.UpdateOne(ctx, q, upd)
	if err != nil {
		return nil, err
	}

	return rets, nil
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
	findopt := options.Find().SetSort(bson.M{"occurDate": -1})
	cursor, err := r.DB.IncomeHistory.Find(ctx, q, findopt)
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
	findopt := options.Find().SetSort(bson.M{"occurDate": -1})
	cursor, err := r.DB.CostHistory.Find(ctx, q, findopt)
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

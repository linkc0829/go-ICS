package resolvers

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	_ "github.com/linkc0829/go-icsharing/internal/db/mongodb"
	dbModel "github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/internal/graph/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"

	tf "github.com/linkc0829/go-icsharing/internal/graph/resolvers/transformer"
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
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	q := bson.M{"owner": user.ID}
	findopt := options.Find().SetSort(bson.M{"occurDate": 1})
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
		return decodeAndFilterPrivacy(me, user, results), nil
	}

	//move expired cost to history if user hasn't queried today
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, time.Local)
	rets := []models.Portfolio{}
	for _, result := range results {
		c := dbModel.CostModel{}
		//encode mongodb result to JSON format, then decode
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &c)
		if c.OccurDate.Before(today) || c.OccurDate.Equal(today) {
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
	return privacyFilter(me, user, rets), nil
}

//GetUserIncome find user's income portfolio, sort them by date and move outdate entries to history
func (r *queryResolver) GetUserIncome(ctx context.Context, id string) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	q := bson.M{"owner": user.ID}
	findopt := options.Find().SetSort(bson.M{"occurDate": 1})
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
		return decodeAndFilterPrivacy(me, user, results), nil
	}

	//move expired cost to history if user hasn't queried today
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, time.Local)
	rets := []models.Portfolio{}
	for _, result := range results {
		c := dbModel.CostModel{}
		//encode mongodb result to JSON format, then decode
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &c)
		if c.OccurDate.Before(today) || c.OccurDate.Equal(today) {
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
	return privacyFilter(me, user, rets), nil
}

func (r *queryResolver) GetUserIncomeHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	fromDate, toDate := getDateRange(rangeArg)
	q := bson.M{"owner": user.ID, "occurDate": bson.M{
		"$gte": fromDate,
		"$lte": toDate,
	}}
	findopt := options.Find().SetSort(bson.M{"occurDate": 1})
	cursor, err := r.DB.IncomeHistory.Find(ctx, q, findopt)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	return decodeAndFilterPrivacy(me, user, results), nil

}

func (r *queryResolver) GetUserCostHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	fromDate, toDate := getDateRange(rangeArg)
	q := bson.M{"owner": user.ID, "occurDate": bson.M{
		"$gte": fromDate,
		"$lte": toDate,
	}}

	findopt := options.Find().SetSort(bson.M{"occurDate": 1})
	cursor, err := r.DB.CostHistory.Find(ctx, q, findopt)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	return decodeAndFilterPrivacy(me, user, results), nil
}

//helper functions

func getDateRange(days int) (time.Time, time.Time) {

	duration, _ := time.ParseDuration(strconv.Itoa(-days*24-24) + "h")
	preDay, _ := time.ParseDuration(strconv.Itoa(-24) + "h")

	to := time.Now().Add(preDay)
	toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.Local)

	from := time.Now().Add(duration)
	fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.Local)

	return fromDate, toDate
}

func decodeAndFilterPrivacy(me *dbModel.UserModel, owner *dbModel.UserModel, results []bson.M) []models.Portfolio {
	hist := []models.Portfolio{}
	for _, result := range results {
		c := dbModel.CostModel{}
		//encode mongodb result to JSON format, then decode
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &c)
		//filter by privacy
		ok := true
		if me.Role != string(models.RoleAdmin) && me.ID != c.Owner { //if not admin and not look at myown portfolio
			//if friend, return friend content by filtering out private content
			if couldViewFriendContent(me, owner) {
				if c.Privacy == models.PrivacyPrivate {
					ok = false
				}
				//not friend, only could see public content
			} else {
				if c.Privacy != models.PrivacyPublic {
					ok = false
				}
			}
		}
		//convert DB User struct to GQL User struct
		if ok {
			gqlCost := tf.DBPortfolioToGQLPortfolio(c)
			hist = append(hist, gqlCost)
		}
	}
	return hist
}

func privacyFilter(me *dbModel.UserModel, owner *dbModel.UserModel, input []models.Portfolio) (result []models.Portfolio) {
	//if admin or lookup myown portfolio
	if me.Role == string(models.RoleAdmin) || me.ID == owner.ID {
		return input
	}
	//if friend, return friend content by filtering out private content
	if couldViewFriendContent(me, owner) {
		for _, in := range input {
			inCost, _ := in.(models.Cost)
			if inCost.Privacy == models.PrivacyPrivate {
				continue
			}
			result = append(result, in)
		}
	} else { //if not friend, return public content
		for _, in := range input {
			inCost, _ := in.(models.Cost)
			if inCost.Privacy != models.PrivacyPublic {
				continue
			}
			result = append(result, in)
		}
	}
	return result
}

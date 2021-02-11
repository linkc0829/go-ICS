package resolvers

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/linkc0829/go-icsharing/internal/db/mongodb"
	dbModel "github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/internal/graph/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"

	tf "github.com/linkc0829/go-icsharing/internal/graph/resolvers/transformer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	return r.GetUserPortfolio(ctx, id, "cost")
}

//GetUserIncome find user's income portfolio, sort them by date and move outdate entries to history
func (r *queryResolver) GetUserIncome(ctx context.Context, id string) ([]models.Portfolio, error) {
	return r.GetUserPortfolio(ctx, id, "income")
}

//GetUserIncome find user's income portfolio, sort them by date and move outdate entries to history
func (r *queryResolver) GetUserPortfolio(ctx context.Context, id string, pType string) ([]models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	user, err := getDBUserByID(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}
	col := r.getPortfolioDB(pType)
	histCol := r.getHistoryDB(pType)
	q := bson.M{"owner": user.ID}
	findopt := options.Find().SetSort(bson.M{"occurDate": 1})
	cursor, err := col.Find(ctx, q, findopt)
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
		return decodeAndFilterPrivacy(me, user, results, pType), nil
	}

	//move expired cost to history if user hasn't queried today
	today := time.Date(yy, mm, dd, 0, 0, 0, 0, time.Local)
	rets := []models.Portfolio{}

	for _, result := range results {
		c := dbModel.PortfolioModel{}
		//encode mongodb result to JSON format, then decode
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &c)
		if c.OccurDate.Before(today) || c.OccurDate.Equal(today) {
			//insert to history
			_, err := histCol.InsertOne(ctx, c)
			if err != nil {
				return nil, err
			}
			//delete old entry
			result, err := col.DeleteOne(ctx, bson.M{"_id": c.ID})
			if err != nil {
				return nil, err
			}
			if result.DeletedCount != 1 {
				return nil, errors.New("delete failed when move old cost entry")
			}
		} else {
			rets = append(rets, tf.DBPortfolioToGQLPortfolio(c, pType))
		}
	}
	//update user LastIncomeQuery
	q = bson.M{"_id": user.ID}
	upd := bson.M{"$set": bson.M{"lastIncomeQuery": time.Now()}}
	if pType == "cost" {
		upd = bson.M{"$set": bson.M{"lastCostQuery": time.Now()}}
	}
	_, err = r.DB.Users.UpdateOne(ctx, q, upd)
	if err != nil {
		return nil, err
	}
	return privacyFilter(me, user, rets), nil
}

func (r *queryResolver) GetUserIncomeHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {
	return r.GetUserPortfolioHistory(ctx, id, rangeArg, "income")
}

func (r *queryResolver) GetUserCostHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {
	return r.GetUserPortfolioHistory(ctx, id, rangeArg, "cost")
}

func (r *queryResolver) GetUserPortfolioHistory(ctx context.Context, id string, rangeArg int, pType string) ([]models.Portfolio, error) {
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

	col := r.getHistoryDB(pType)
	findopt := options.Find().SetSort(bson.M{"occurDate": 1})
	cursor, err := col.Find(ctx, q, findopt)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	return decodeAndFilterPrivacy(me, user, results, pType), nil
}

//helper functions

func (r *mutationResolver) CreatePortfolio(ctx context.Context, input models.CreatePortfolioInput, pType string) (*models.Portfolio, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)

	newPortfolio := dbModel.PortfolioModel{
		ID:          primitive.NewObjectID(),
		Owner:       me.ID,
		Amount:      input.GetAmount(),
		OccurDate:   input.GetOccurDate(),
		Description: input.GetDescription(),
		Vote:        nil,
		Category:    input.GetCategory(),
		Privacy:     input.GetPrivacy(),
	}
	col := r.getDB(pType)
	//insert to db
	_, err := col.InsertOne(ctx, newPortfolio)
	if err != nil {
		return nil, err
	}
	result := tf.DBPortfolioToGQLPortfolio(newPortfolio, pType)

	return &result, nil
}

func (r *mutationResolver) UpdatePortfolio(ctx context.Context, id string, input models.UpdatePortfolioInput, pType string) (*models.Portfolio, error) {

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)

	q := bson.M{"_id": ID}
	col := r.getDB(pType)
	result := dbModel.PortfolioModel{}
	portInput := models.UpdatePortfolio{
		Amount:      input.GetAmount(),
		OccurDate:   input.GetOccurDate(),
		Category:    input.GetCategory(),
		Description: input.GetDescription(),
		Privacy:     input.GetPrivacy(),
	}

	if err := col.FindOne(ctx, q).Decode(&result); err != nil {
		return nil, err
	}
	if !isAdmin(ctx) && me.ID != result.Owner {
		return nil, errors.New("permission denied")
	}
	if portInput.Amount != nil {
		result.Amount = *portInput.Amount
	}
	if portInput.Category != nil {
		result.Category = *portInput.Category
	}
	if portInput.Description != nil {
		result.Description = *portInput.Description
	}
	if portInput.OccurDate != nil {
		result.OccurDate = *portInput.OccurDate
	}
	if portInput.Privacy != nil {
		result.Privacy = *portInput.Privacy
	}

	upd := bson.M{"$set": portInput}
	_, err = col.UpdateOne(ctx, q, upd)
	if err != nil {
		return nil, err
	}

	ret := tf.DBPortfolioToGQLPortfolio(result, pType)
	return &ret, nil
}

func (r *mutationResolver) DeletePortfolio(ctx context.Context, id string, pType string) (bool, error) {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	q := bson.M{"_id": primID}
	result := dbModel.PortfolioModel{}
	col := r.getDB(pType)
	if err := col.FindOne(ctx, q).Decode(&result); err != nil {
		return false, err
	}

	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	if !isAdmin(ctx) && me.ID != result.Owner {
		return false, errors.New("permission denied")
	}

	delete, err := col.DeleteOne(ctx, q)
	if err != nil {
		return false, err
	}
	if delete.DeletedCount == 1 {
		return true, nil
	}
	return false, nil
}

func (r *mutationResolver) VotePortfolio(ctx context.Context, id string, pType string) (int, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	q := bson.M{"_id": ID}
	target := dbModel.PortfolioModel{}
	col := r.getDB(pType)
	if err := col.FindOne(ctx, q).Decode(&target); err != nil {
		return -1, err
	}
	owner := dbModel.UserModel{}
	if err := r.DB.Users.FindOne(ctx, bson.M{"_id": target.Owner}).Decode(&owner); err != nil {
		return -1, err
	}
	//deny private
	if !isAdmin(ctx) && me.ID != target.Owner && target.Privacy == models.PrivacyPrivate {
		return -1, errors.New("it's private, permission denied")
	}
	//deny non-friend
	if !isAdmin(ctx) && !couldViewFriendContent(me, &owner) && me.ID != target.Owner && target.Privacy == models.PrivacyFriend {
		return -1, errors.New("it's only for friend, permission denied")
	}
	//if already voted, revoke
	length := len(target.Vote)
	for i, v := range target.Vote {
		if v == me.ID {
			if length == 1 {
				target.Vote = target.Vote[:0]
			} else {
				target.Vote[i] = target.Vote[length-1]
				target.Vote = target.Vote[:length-1]
			}
			break
		}
	}
	if length == len(target.Vote) {
		//add to vote
		target.Vote = append(target.Vote, me.ID)
	}

	//update DB
	q = bson.M{"_id": target.ID, "voteVer": target.VoteVer}
	upd := bson.M{"$set": bson.M{"vote": target.Vote, "voteVer": (target.VoteVer + 1)}}
	result, err := col.UpdateOne(ctx, q, upd)
	if err != nil {
		log.Fatal(err)
	}
	if result.ModifiedCount == 0 {
		log.Println("CommitIncomeVote modify unsucceed, retry")
		MongoPortfolioOCT(&target, me.ID)
	}
	return len(target.Vote), nil
}

//MongoDB optimistic concurancy transaction
func MongoPortfolioOCT(portfolio *dbModel.PortfolioModel, voter primitive.ObjectID) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		channelNumber := rand.Intn(100) % 10
		result := make(chan []primitive.ObjectID)
		mongodb.PortfolioChan[channelNumber] <- mongodb.PortfolioData{
			Portfolio: *portfolio,
			Voter:     voter,
			Result:    &result,
		}
		select {
		case portfolio.Vote = <-result:
			wg.Done()
		}
	}(&wg)
	wg.Wait()
}

func getDateRange(days int) (time.Time, time.Time) {

	duration, _ := time.ParseDuration(strconv.Itoa(-days*24-24) + "h")
	preDay, _ := time.ParseDuration(strconv.Itoa(-24) + "h")

	to := time.Now().Add(preDay)
	toDate := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.Local)

	from := time.Now().Add(duration)
	fromDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.Local)

	return fromDate, toDate
}

func decodeAndFilterPrivacy(me *dbModel.UserModel, owner *dbModel.UserModel, results []bson.M, pType string) []models.Portfolio {
	hist := []models.Portfolio{}
	for _, result := range results {
		c := dbModel.PortfolioModel{}
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
			gqlCost := tf.DBPortfolioToGQLPortfolio(c, pType)
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

func (r *mutationResolver) getDB(pType string) (col *mongo.Collection) {
	if pType == "cost" {
		col = r.DB.Cost
	} else {
		col = r.DB.Income
	}
	return
}

func (r *queryResolver) getPortfolioDB(pType string) (col *mongo.Collection) {
	if pType == "cost" {
		col = r.DB.Cost
	} else {
		col = r.DB.Income
	}
	return
}

func (r *queryResolver) getHistoryDB(pType string) (col *mongo.Collection) {
	if pType == "cost" {
		col = r.DB.CostHistory
	} else {
		col = r.DB.IncomeHistory
	}
	return
}

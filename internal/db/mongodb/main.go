package mongodb

import (
	"context"
	"log"

	"github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbUser, pwd, dbName, dsn string
var IncomeChan []chan IncomeData
var CostChan []chan CostData

type IncomeData struct {
	Income models.IncomeModel
	Voter  primitive.ObjectID
	Result *chan []primitive.ObjectID
}

type CostData struct {
	Cost   models.CostModel
	Voter  primitive.ObjectID
	Result *chan []primitive.ObjectID
}

type MongoDB struct {
	Session       *mongo.Client
	Users         *mongo.Collection
	DeletedUsers  *mongo.Collection
	Income        *mongo.Collection
	DeletedIncome *mongo.Collection
	Cost          *mongo.Collection
	DeletedCost   *mongo.Collection
	IncomeHistory *mongo.Collection
	CostHistory   *mongo.Collection
}

//ConnectDB will build connection to MongoDB Atlas
func ConnectMongoDB(cfg *utils.ServerConfig) (db *MongoDB) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.DSN))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	IncomeChan = make([]chan IncomeData, 10)

	db = &MongoDB{
		Session:       client,
		Users:         client.Database("ics").Collection("users"),
		DeletedUsers:  client.Database("ics").Collection("deletedUsers"),
		Income:        client.Database("ics").Collection("income"),
		DeletedIncome: client.Database("ics").Collection("deletedIncome"),
		Cost:          client.Database("ics").Collection("cost"),
		DeletedCost:   client.Database("ics").Collection("deletedCost"),
		IncomeHistory: client.Database("ics").Collection("incomeHistory"),
		CostHistory:   client.Database("ics").Collection("costHistory"),
	}
	for i := range IncomeChan {
		IncomeChan[i] = make(chan IncomeData)
		go CommitIncomeVote(context.Background(), &IncomeChan[i], db)
	}
	return
}

//CloseDB will dissconnect to MongoDB
func CloseMongoDB(db *MongoDB) {
	err := db.Session.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

//implement mongodb transaction for vote income
func CommitIncomeVote(ctx context.Context, in *chan IncomeData, db *MongoDB) {
	for {
		select {
		case data := <-*in:
		Loop:
			//step1: reload data from mongo
			q := bson.M{"_id": data.Income.ID}
			income := models.IncomeModel{}
			if err := db.Income.FindOne(context.TODO(), q).Decode(&income); err != nil {
				panic("mongoDB err during commitIncomeVote: " + err.Error())
			}
			//step2: exec vote
			//if already voted, revoke
			length := len(income.Vote)
			for i, v := range income.Vote {
				if v == data.Voter {
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
				income.Vote = append(income.Vote, data.Voter)
			}
			//step3: update DB
			q = bson.M{"_id": income.ID, "voteVer": income.VoteVer}
			upd := bson.M{"$set": bson.M{"vote": income.Vote, "voteVer": income.VoteVer + 1}}
			_, err := db.Income.UpdateOne(ctx, q, upd)
			if err != nil {
				log.Println("got error during commitVoteIncome: " + err.Error())
				goto Loop
			}
			*data.Result <- income.Vote
		}
	}
}

//implement mongodb transaction for vote cost
func CommitCostVote(ctx context.Context, in *chan CostData, db *MongoDB) {
	for {
		select {
		case data := <-*in:
		Loop:
			//step1: reload data from mongo
			q := bson.M{"_id": data.Cost.ID}
			cost := models.CostModel{}
			if err := db.Cost.FindOne(context.TODO(), q).Decode(&cost); err != nil {
				panic("mongoDB err during commitCostVote: " + err.Error())
			}
			//step2: exec vote
			//if already voted, revoke
			length := len(cost.Vote)
			for i, v := range cost.Vote {
				if v == data.Voter {
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
				cost.Vote = append(cost.Vote, data.Voter)
			}
			//step3: update DB
			q = bson.M{"_id": cost.ID, "voteVer": cost.VoteVer}
			upd := bson.M{"$set": bson.M{"vote": cost.Vote, "voteVer": cost.VoteVer + 1}}
			_, err := db.Cost.UpdateOne(ctx, q, upd)
			if err != nil {
				log.Println("got error during commitVoteCost: " + err.Error())
				goto Loop
			}
			*data.Result <- cost.Vote
		}
	}
}

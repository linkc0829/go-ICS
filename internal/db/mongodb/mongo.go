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
var PortfolioChan []chan PortfolioData

type PortfolioData struct {
	Portfolio models.PortfolioModel
	Voter     primitive.ObjectID
	DB        *mongo.Collection
	Result    *chan []primitive.ObjectID
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
func ConnectDB(cfg *utils.ServerConfig) (db *MongoDB) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.DSN))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

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
	db.initMultipleQueue(1)
	return db
}

//init multiple queue for mongoDB optimistic concurrency transaction
func (db *MongoDB) initMultipleQueue(maxThread int) {
	PortfolioChan = make([]chan PortfolioData, maxThread)

	for i := range PortfolioChan {
		PortfolioChan[i] = make(chan PortfolioData)
		go db.CommitPortfolioVote(context.Background(), &PortfolioChan[i], i)
	}

}

//CloseDB will dissconnect to MongoDB
func (db *MongoDB) CloseDB() {
	err := db.Session.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

//implement mongodb transaction for vote income
func (db *MongoDB) CommitPortfolioVote(ctx context.Context, in *chan PortfolioData, i int) {
	for {
		select {
		case data := <-*in:
		Loop:
			//step1: reload data from mongo
			q := bson.M{"_id": data.Portfolio.ID}
			target := models.PortfolioModel{}
			if err := db.Income.FindOne(context.TODO(), q).Decode(&target); err != nil {
				panic("mongoDB err during commitIncomeVote: " + err.Error())
			}
			//step2: exec vote
			//if already voted, revoke
			length := len(target.Vote)
			for i, v := range target.Vote {
				if v == data.Voter {
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
				target.Vote = append(target.Vote, data.Voter)
			}
			//step3: update DB
			q = bson.M{"_id": target.ID, "voteVer": target.VoteVer}
			upd := bson.M{"$set": bson.M{"vote": target.Vote, "voteVer": (target.VoteVer + 1)}}
			result, err := data.DB.UpdateOne(ctx, q, upd)
			if err != nil {
				log.Fatal(err)
			}
			if result.ModifiedCount == 0 {
				log.Printf("CommitIncomeVote modify %d voteVer unsucceed in channel #%d, retry", target.VoteVer, i)
				goto Loop
			}
			*data.Result <- target.Vote
		}
	}
}

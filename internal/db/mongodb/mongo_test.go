package mongodb

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_CommitPortfolioVote(t *testing.T) {
	runs := 50

	mongoRoot := utils.MustGet("MONGO_INITDB_ROOT_USERNAME")
	mongoRootPWD := utils.MustGet("MONGO_INITDB_ROOT_PASSWORD")
	mongoHost := utils.MustGet("MONGO_HOST")
	connectDB := utils.MustGet("MONGO_INITDB_DATABASE")
	mongoDSN := "mongodb://" + mongoRoot + ":" + mongoRootPWD + "@" + mongoHost + "/" + connectDB + "?authSource=admin"

	demo := utils.MustGet("DEMO_MODE")
	if demo == "on" {
		mongoDSN = utils.MustGet("MONGO_CONNECTION_DSN")
	}

	conf := &utils.ServerConfig{
		MongoDB: utils.MGDBConfig{
			DSN: mongoDSN,
		},
	}
	db := ConnectDB(conf)
	defer db.CloseDB()
	ctx := context.Background()

	t.Run(fmt.Sprintf("test Vote Income for %d times", runs), func(t *testing.T) {
		income := models.PortfolioModel{
			ID:          primitive.NewObjectID(),
			Description: "TEST VOTE",
			VoteVer:     0,
		}
		_, err := db.Income.InsertOne(ctx, income)
		if err != nil {
			t.Fatal(err)
		}

		voter := []primitive.ObjectID{}
		for i := 0; i < runs; i++ {
			voter = append(voter, primitive.NewObjectID())
		}
		//vote many times in background
		wg := sync.WaitGroup{}
		wg.Add(runs)
		for _, v := range voter {
			go MongoPortfolioOCT(&wg, income, v, db.Income)
		}
		wg.Wait()

		q := bson.M{"_id": income.ID}
		err = db.Income.FindOne(ctx, q).Decode(&income)
		if err != nil {
			t.Fatal(err)
		}
		if income.VoteVer != runs {
			t.Error("wrong voteVer")
		}
		if len(income.Vote) != runs {
			t.Error("wrong Vote")
		}

		//vote again, should get zero vote
		wg = sync.WaitGroup{}
		wg.Add(runs)
		for _, v := range voter {
			go MongoPortfolioOCT(&wg, income, v, db.Income)
		}
		wg.Wait()

		q = bson.M{"_id": income.ID}
		err = db.Income.FindOne(ctx, q).Decode(&income)
		if err != nil {
			t.Fatal(err)
		}

		if len(income.Vote) != 0 {
			t.Error("wrong Vote, should be zero")
		}

		_, err = db.Income.DeleteOne(ctx, q)
		if err != nil {
			t.Fatal(err)
		}
	})
}

//MongoDB optimistic concurancy transaction
func MongoPortfolioOCT(wg *sync.WaitGroup, portfolio models.PortfolioModel, voter primitive.ObjectID, db *mongo.Collection) {

	channelNumber := rand.Intn(100) % 10
	result := make(chan []primitive.ObjectID)
	PortfolioChan[channelNumber] <- PortfolioData{
		Portfolio: portfolio,
		Voter:     voter,
		DB:        db,
		Result:    &result,
	}
	select {
	case portfolio.Vote = <-result:
		wg.Done()
	}
}

func BenchmarkInitMultipleQueue(b *testing.B) {
	runs := 1000

	mongoRoot := utils.MustGet("MONGO_INITDB_ROOT_USERNAME")
	mongoRootPWD := utils.MustGet("MONGO_INITDB_ROOT_PASSWORD")
	mongoHost := utils.MustGet("MONGO_HOST")
	connectDB := utils.MustGet("MONGO_INITDB_DATABASE")
	mongoDSN := "mongodb://" + mongoRoot + ":" + mongoRootPWD + "@" + mongoHost + "/" + connectDB + "?authSource=admin"

	demo := utils.MustGet("DEMO_MODE")
	if demo == "on" {
		mongoDSN = utils.MustGet("MONGO_CONNECTION_DSN")
	}

	conf := &utils.ServerConfig{
		MongoDB: utils.MGDBConfig{
			DSN: mongoDSN,
		},
	}
	db := ConnectDB(conf)
	defer db.CloseDB()
	ctx := context.Background()

	income := models.PortfolioModel{
		ID:          primitive.NewObjectID(),
		Description: "TEST VOTE",
		VoteVer:     0,
	}
	_, err := db.Income.InsertOne(ctx, income)
	if err != nil {
		fmt.Println(err)
	}
	voter := []primitive.ObjectID{}
	for i := 0; i < runs; i++ {
		voter = append(voter, primitive.NewObjectID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//vote many times in background
		wg := sync.WaitGroup{}
		wg.Add(runs)
		for _, v := range voter {
			go MongoPortfolioOCT(&wg, income, v, db.Income)
		}
		wg.Wait()
	}
	q := bson.M{"_id": income.ID}
	_, err = db.Income.DeleteOne(ctx, q)
	if err != nil {
		fmt.Println(err)
	}
}

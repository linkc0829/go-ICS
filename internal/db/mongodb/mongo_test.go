package mongodb

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestConnectMongoDB(t *testing.T) {
	runs := 50
	conf := &utils.ServerConfig{
		MongoDB: utils.MGDBConfig{
			DSN: utils.MustGet("MONGO_CONNECTION_DSN"),
		},
	}
	db := ConnectMongoDB(conf)
	ctx := context.Background()

	income := models.IncomeModel{
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
	wg := sync.WaitGroup{}
	wg.Add(runs)
	for _, v := range voter {
		go MongoIncomeOCT(&wg, income, v)
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
		go MongoIncomeOCT(&wg, income, v)
	}
	wg.Wait()

	q = bson.M{"_id": income.ID}
	err = db.Income.FindOne(ctx, q).Decode(&income)
	if err != nil {
		t.Fatal(err)
	}

	if len(income.Vote) != 0 {
		t.Error("wrong Vote, should be zero")
		log.Println(len(income.Vote))
	}

	_, err = db.Income.DeleteOne(ctx, q)
	if err != nil {
		t.Fatal(err)
	}
}

//MongoDB optimistic concurancy transaction
func MongoIncomeOCT(wg *sync.WaitGroup, income models.IncomeModel, voter primitive.ObjectID) {

	channelNumber := rand.Intn(100) % 10
	result := make(chan []primitive.ObjectID)
	IncomeChan[channelNumber] <- IncomeData{
		Income: income,
		Voter:  voter,
		Result: &result,
	}
	select {
	case income.Vote = <-result:
		wg.Done()
	}
}

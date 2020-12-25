package mongodb

import (
	"context"
	"log"

	"github.com/linkc0829/go-ics/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbUser, pwd, dbName, dsn string

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
func ConnectDB(cfg *utils.ServerConfig) *MongoDB {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.DSN))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return &MongoDB{
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
}

//CloseDB will dissconnect to MongoDB
func CloseDB(db *MongoDB) {
	err := db.Session.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

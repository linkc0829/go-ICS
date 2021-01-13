package main

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/linkc0829/go-ics/internal/db/mongodb"
	"github.com/linkc0829/go-ics/internal/db/redisdb"
	"github.com/linkc0829/go-ics/internal/db/sqlitedb"
	"github.com/linkc0829/go-ics/pkg/server"
	"github.com/linkc0829/go-ics/pkg/utils/datasource"
)

func TestRestAPI(t *testing.T) {

	mongoDB := mongodb.ConnectMongoDB(serverconf)
	defer mongodb.CloseMongoDB(mongoDB)

	sqlite := sqlitedb.ConnectSqlite()
	defer sqlitedb.CloseSqlite(sqlite)

	redis := redisdb.ConnectRedis(serverconf)
	defer redisdb.CloseRedis(redis)

	db := &datasource.DB{
		Mongo:  mongoDB,
		Sqlite: sqlite,
		Redis:  redis,
	}

	r := server.SetupServer(serverconf, db)

	ts := httptest.NewServer(r)
	defer ts.Close()

	fmt.Println(ts.URL)
}

func TestGraphAPI(t *testing.T) {

}

func TestAuthAPI(t *testing.T) {

}

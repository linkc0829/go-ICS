package main

import(
	"github.com/linkc0829/go-ics/pkg/server"
	"github.com/linkc0829/go-ics/internal/mongodb"
)

func main(){

	db := mongodb.ConnectDB()
	defer mongodb.CloseDB(db)
	server.Run(db)
}
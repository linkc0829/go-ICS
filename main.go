package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/linkc0829/go-icsharing/internal/db/mongodb"
	"github.com/linkc0829/go-icsharing/internal/db/redisdb"
	"github.com/linkc0829/go-icsharing/internal/db/sqlitedb"
	"github.com/linkc0829/go-icsharing/pkg/server"
	"github.com/linkc0829/go-icsharing/pkg/utils"
	"github.com/linkc0829/go-icsharing/pkg/utils/datasource"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var serverconf *utils.ServerConfig

func init() {
	demo := utils.MustGet("DEMO_MODE")

	mongoRoot := utils.MustGet("MONGO_INITDB_ROOT_USERNAME")
	mongoRootPWD := utils.MustGet("MONGO_INITDB_ROOT_PASSWORD")
	mongoHost := utils.MustGet("MONGO_HOST")
	connectDB := utils.MustGet("MONGO_INITDB_DATABASE")
	mongoDSN := "mongodb://" + mongoRoot + ":" + mongoRootPWD + "@" + mongoHost + "/" + connectDB + "?authSource=admin"
	redisEndpoint := utils.MustGet("REDIS_HOST")
	port := utils.MustGet("SERVER_PORT")
	heroku := utils.MustGet("ISHEROKU")
	//heroku network setting
	if heroku == "true" {
		port = utils.MustGet("PORT")
		log.Println("Deploy in Heroku, port: " + port)
	}

	if demo == "on" {
		mongoDSN = utils.MustGet("MONGO_CONNECTION_DSN")
		redisEndpoint = utils.MustGet("REDIS_ENDPOINT")
	}

	serverconf = &utils.ServerConfig{
		Host:          utils.MustGet("SERVER_HOST"),
		Port:          port,
		URISchema:     utils.MustGet("SERVER_URI_SCHEMA"),
		ApiVer:        utils.MustGet("SERVER_PATH_VERSION"),
		SessionSecret: utils.MustGet("SESSION_SECRET"),
		StaticPath:    utils.MustGet("SERVER_STATIC_PATH"),
		CookiesDomain: utils.MustGet("COOKIES_DOMAIN"),
		JWT: utils.JWTConfig{
			Secret:             utils.MustGet("AUTH_JWT_SECRET"),
			Algorithm:          utils.MustGet("AUTH_JWT_SIGNING_ALGORITHM"),
			AccessTokenExpire:  utils.MustGet("AUTH_JWT_ACCESSTOKEN_EXPIRE"),
			RefreshTokenExpire: utils.MustGet("AUTH_JWT_REFRESHTOKEN_EXPIRE"),
		},
		GraphQL: utils.GQLConfig{
			Path:                utils.MustGet("GQL_SERVER_GRAPHQL_PATH"),
			PlaygroundPath:      utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH"),
			IsPlaygroundEnabled: utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED"),
		},
		MongoDB: utils.MGDBConfig{
			DSN: mongoDSN,
		},
		Redis: utils.RedisConfig{
			EndPoint: redisEndpoint,
			PWD:      utils.MustGet("REDIS_PWD"),
		},
		AuthProviders: []utils.AuthProvider{
			utils.AuthProvider{
				Provider:  "google",
				ClientKey: utils.MustGet("PROVIDER_GOOGLE_KEY"),
				Secret:    utils.MustGet("PROVIDER_GOOGLE_SECRET"),
			},
			utils.AuthProvider{
				Provider:  "auth0",
				ClientKey: utils.MustGet("PROVIDER_AUTH0_KEY"),
				Secret:    utils.MustGet("PROVIDER_AUTH0_SECRET"),
				Domain:    utils.MustGet("PROVIDER_AUTH0_DOMAIN"),
				Scopes:    strings.Split(utils.MustGet("PROVIDER_AUTH0_SCOPES"), ","),
			},
		},
	}
}

func main() {

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

	go func() {
		id := primitive.NewObjectID()
		mongoDB.Users.InsertOne(context.TODO(), bson.M{"_id": id, "userid": "testgoroutine"})
	}()

	if serverconf.URISchema == "https://" {
		r := server.SetupServer(serverconf, db)
		//using staging environment
		//https://community.letsencrypt.org/t/golang-autocert-staging-environment/128912/3
		m := &autocert.Manager{
			Client: &acme.Client{
				DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
			},
			Cache:      autocert.DirCache("secret-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("icsharing.com", "www.icsharing.com"),
		}
		//http
		go r.Run(":80")
		log.Fatal(http.Serve(m.Listener(), r))

		//use above if exceed renewal limits
		//log.Fatal(http.Serve(autocert.NewListener("icsharing.com", "www.icsharing.com"), r))
	} else {
		log.Fatal(server.SetupServer(serverconf, db).Run(":" + serverconf.Port))
	}
}

//helper function for testing
func getServer() http.Handler {
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
	return server.SetupServer(serverconf, db)
}

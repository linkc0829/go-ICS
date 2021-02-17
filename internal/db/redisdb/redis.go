package redisdb

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/linkc0829/go-icsharing/pkg/utils"
)

type RedisDB struct {
	Conn redis.Conn
}

func ConnectDB(cfg *utils.ServerConfig) *RedisDB {
	conn, err := redis.Dial("tcp", cfg.Redis.EndPoint, redis.DialPassword(cfg.Redis.PWD))
	if err != nil {
		log.Fatal(err)
	}
	db := &RedisDB{
		Conn: conn,
	}
	return db
}

func (db *RedisDB) CloseDB() {
	err := db.Conn.Close()
	if err != nil {
		log.Fatal(err)
	}

}

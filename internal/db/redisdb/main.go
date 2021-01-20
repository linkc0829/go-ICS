package redisdb

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func ConnectRedis(cfg *utils.ServerConfig) redis.Conn {
	conn, err := redis.Dial("tcp", "redis:6379", redis.DialPassword(cfg.Redis.PWD))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func CloseRedis(c redis.Conn) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}

}

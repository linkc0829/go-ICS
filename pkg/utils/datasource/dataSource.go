package datasource

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/linkc0829/go-icsharing/internal/db/mongodb"
)

type DB struct {
	Mongo  *mongodb.MongoDB
	Sqlite *gorm.DB
	Redis  redis.Conn
}

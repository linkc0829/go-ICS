package sqlitedb

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Portfolio struct {
	ID          *uint      `json:"id,omitempty" gorm:"id,primaryKey,autoIncrement"`
	Amount      *int       `json:"amount,string,omitempty" gorm:"amount"`
	OccurDate   *time.Time `json:"occurDate,omitempty" gorm:"occurDate"`
	Category    *string    `json:"category,omitempty" gorm:"category"`
	Description *string    `json:"description,omitempty" gorm:"description"`
	Type        *string    `json:"type,omitempty" gorm:"type"`
	Privacy     *string    `json:"privacy,omitempty" gorm:"privacy"`
}

type Portfolios struct {
	Portfolios []Portfolio `json:"portfolios"`
}

type SqliteDB struct {
	Client *gorm.DB
}

// TableName set the name of the table.
func (Portfolio) TableName() string {
	return "portfolio"
}

func ConnectDB() *SqliteDB {
	client, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	if !client.HasTable(&Portfolio{}) {
		client.CreateTable(&Portfolio{})
	}
	db := &SqliteDB{
		Client: client,
	}
	return db
}

func (db *SqliteDB) CloseDB() {
	if err := db.Client.Close(); err != nil {
		log.Fatal(err)
	}
}

package models

import (
	"time"

	"github.com/linkc0829/go-ics/internal/graph/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserModel struct
type UserModel struct {
	ID              primitive.ObjectID   `bson:"_id" json:"_id,omitempty"`
	UserID          string               `bson:"userid" json:"userid,omitempty"`
	Password        *string              `bson:"password" json:"password,omitempty"`
	Email           string               `bson:"email" json:"email,omitempty"`
	NickName        *string              `bson:"nickname" json:"nickname,omitempty"`
	CreatedAt       time.Time            `bson:"createAt" json:"createAt,omitempty"`
	Friends         []primitive.ObjectID `bson:"friends" json:"friends,omitempty"`
	RefreshToken    string               `bson:"refreshToken" json:"refreshToken,omitempty"`
	LastIncomeQuery time.Time            `bson:"lastIncomeQuery" json:"lastIncomeQuery,omitempty"`
	LastCostQuery   time.Time            `bson:"lastCostQuery" json:"lastCostQuery,omitempty"`
	Provider        string               `bson:"provider" json:"provider,omitempty"`
	AvatarURL       string               `bson:"avatarURL" json:"avatarURL,omitempty"`
	Role            string               `bson:"role" json:"role,omitempty"`
}

//role constants
const (
	ADMIN string = "ADMIN"
	USER  string = "USER"
)

type PortfolioModel interface {
	IsPortfolio()
}

type IncomeModel struct {
	ID          primitive.ObjectID       `bson:"_id" json:"_id,omitempty"`
	Owner       primitive.ObjectID       `bson:"owner" json:"owner,omitempty"`
	Amount      int                      `bson:"amount" json:"amount,omitempty"`
	OccurDate   time.Time                `bson:"occurDate" json:"occurDate,omitempty"`
	Category    models.PortfolioCategory `bson:"category" json:"category,omitempty"`
	Description *string                  `bson:"description" json:"description,omitempty"`
	Vote        []primitive.ObjectID     `bson:"vote" json:"vote,omitempty"`
}

func (IncomeModel) IsPortfolio() {}

type CostModel struct {
	ID          primitive.ObjectID       `bson:"_id" json:"_id,omitempty"`
	Owner       primitive.ObjectID       `bson:"owner" json:"owner,omitempty"`
	Amount      int                      `bson:"amount" json:"amount,omitempty"`
	OccurDate   time.Time                `bson:"occurDate" json:"occurDate,omitempty"`
	Category    models.PortfolioCategory `bson:"category" json:"category,omitempty"`
	Description *string                  `bson:"description" json:"description,omitempty"`
	Vote        []primitive.ObjectID     `bson:"vote" json:"vote,omitempty"`
}

func (CostModel) IsPortfolio() {}

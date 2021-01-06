package models

import (
	"time"
)

type Portfolio interface {
	IsPortfolio()
}

type Cost struct {
	ID          string            `json:"id"`
	Owner       string            `json:"owner"`
	Amount      int               `json:"amount,string"`
	OccurDate   time.Time         `json:"occurDate"`
	Description string            `json:"description"`
	Vote        []string          `json:"vote"`
	Category    PortfolioCategory `json:"category"`
}

func (Cost) IsPortfolio() {}

type CreateCostInput struct {
	Amount      int          `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   time.Time    `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    CostCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description string       `bson:"description,omitempty" json:"description,omitempty"`
}

type CreateIncomeInput struct {
	Amount      int            `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   time.Time      `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    IncomeCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description string         `bson:"description,omitempty" json:"description,omitempty"`
}

type CreateUserInput struct {
	Email    string  `bson:"email,omitempty" json:"email,omitempty"`
	UserID   string  `bson:"userID,omitempty" json:"userId,omitempty"`
	NickName *string `bson:"nickName,omitempty" json:"nickName,omitempty"`
}

type Income struct {
	ID          string            `json:"id"`
	Owner       string            `json:"owner"`
	Amount      int               `json:"amount,string"`
	OccurDate   time.Time         `json:"occurDate"`
	Description string            `json:"description"`
	Vote        []string          `json:"vote"`
	Category    PortfolioCategory `json:"category"`
}

func (Income) IsPortfolio() {}

type UpdateCostInput struct {
	Amount      *int          `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   *time.Time    `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    *CostCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description *string       `bson:"description,omitempty" json:"description,omitempty"`
}

type UpdateIncomeInput struct {
	Amount      *int            `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   *time.Time      `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    *IncomeCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description *string         `bson:"description,omitempty" json:"description,omitempty"`
}

type UpdateUserInput struct {
	Email    *string `bson:"email,omitempty" json:"email,omitempty"`
	UserID   *string `bson:"userID,omitempty" json:"userId,omitempty"`
	NickName *string `bson:"nickName,omitempty" json:"nickName,omitempty"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	UserID    string    `json:"userId"`
	NickName  *string   `json:"nickName"`
	CreatedAt time.Time `json:"createdAt"`
	// granted permission to friends to view portfolio
	Friends []string `json:"friends"`
	// permission to view followers portfolio
	Followers []string `json:"followers"`
	Role      Role     `json:"role"`
}

type UserInput struct {
	Email    *string `json:"email"`
	UserID   *string `json:"userId"`
	NickName *string `json:"nickName"`
}

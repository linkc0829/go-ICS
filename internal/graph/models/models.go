package models

import (
	"time"
)

type Portfolio interface {
	IsPortfolio()
}

type Cost struct {
	ID          string       `json:"id"`
	Owner       string       `json:"owner"`
	Amount      int          `json:"amount"`
	OccurDate   time.Time    `json:"occurDate"`
	Description *string      `json:"description"`
	Vote        []string     `json:"vote"`
	Category    PortfolioCategory `json:"category"`
}

func (Cost) IsPortfolio() {}

type CreateCostInput struct {
	Amount      int          `json:"amount"`
	OccurDate   time.Time    `json:"occurDate"`
	Category    CostCategory `json:"category"`
	Description *string      `json:"description"`
}

type CreateIncomeInput struct {
	Amount      int            `json:"amount"`
	OccurDate   time.Time      `json:"occurDate"`
	Category    IncomeCategory `json:"category"`
	Description *string        `json:"description"`
}

type CreateUserInput struct {
	Email    string  `json:"email"`
	UserID   string  `json:"userId"`
	NickName *string `json:"nickName"`
}

type Income struct {
	ID          string         `json:"id"`
	Owner       string         `json:"owner"`
	Amount      int            `json:"amount"`
	OccurDate   time.Time      `json:"occurDate"`
	Description *string        `json:"description"`
	Vote        []string       `json:"vote"`
	Category    PortfolioCategory `json:"category"`
}

func (Income) IsPortfolio() {}

type UpdateCostInput struct {
	Amount      *int          `bson:"amount" json:"amount,omitempty"`
	OccurDate   *time.Time    `bson:"occurDate" json:"occurDate,omitempty"`
	Category    *CostCategory `bson:"category" json:"category,omitempty"`
	Description *string       `bson:"description" json:"description,omitempty"`
}

type UpdateIncomeInput struct {
	Amount      *int            `json:"amount"`
	OccurDate   *time.Time      `json:"occurDate"`
	Category    *IncomeCategory `json:"category"`
	Description *string         `json:"description"`
}

type UpdateUserInput struct {
	Email    *string `json:"email"`
	UserID   *string `json:"userId"`
	NickName *string `json:"nickName"`
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
}

type UserInput struct {
	Email    *string `json:"email"`
	UserID   *string `json:"userId"`
	NickName *string `json:"nickName"`
}

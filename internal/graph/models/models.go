package models

import (
	"time"
)

type Portfolio interface {
	IsPortfolio()
}

type Cost struct {
	ID          string    `json:"id"`
	Owner       string    `json:"owner"`
	Amount      int       `json:"amount"`
	OccurDate   time.Time `json:"occurDate"`
	Description *string   `json:"description"`
	Vote        []string  `json:"vote"`
	Category    CostCategory  `json:"category"`

}

func (Cost) IsPortfolio() {}

type CostInput struct {
	Amount      *int       		`json:"amount"`
	OccurDate   time.Time 		`json:"occurDate"`
	Category    *CostCategory   `json:"category"`
	Description *string    		`json:"description"`
}

type Income struct {
	ID          string    `json:"id"`
	Owner       string    `json:"owner"`
	Amount      int       `json:"amount"`
	OccurDate   time.Time `json:"occurDate"`
	Description *string   `json:"description"`
	Vote        []string  `json:"vote"`
	Category    IncomeCategory  `json:"category"`
}

func (Income) IsPortfolio() {}

type IncomeInput struct {
	Amount      *int       		`json:"amount"`
	OccurDate   time.Time 		`json:"occurDate"`
	Category    *IncomeCategory `json:"category"`
	Description *string    		`json:"description"`
}


type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	UserID    string    `json:"userId"`
	NickName  *string   `json:"nickName"`
	CreatedAt time.Time `json:"createdAt"`
	// granted permission to friends to view portfolio
	Friends   []string  `json:"friends"`
	// permission to view followers portfolio
	Followers []string  `json:"followers"`
}

type UserInput struct {
	Email    *string `json:"email"`
	UserID   *string `json:"userId"`
	NickName *string `json:"nickName"`
}

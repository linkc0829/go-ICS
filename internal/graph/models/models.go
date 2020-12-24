package models

import (
	"time"
)

type Cost struct {
	ID          string    `json:"id"`
	Owner       string    `json:"owner"`
	Amount      int       `json:"amount"`
	OccurDate   time.Time `json:"occurDate"`
	Category    Category  `json:"category"`
	Description *string   `json:"description"`
	Vote        []string `json:"vote"`
}

type CostInput struct {
	Amount      *int       `json:"amount"`
	Date        *time.Time `json:"date"`
	Category    *Category  `json:"category"`
	Description *string    `json:"description"`
}

type Income struct {
	ID          string    `json:"id"`
	Owner       string    `json:"owner"`
	Amount      int       `json:"amount"`
	OccurDate   time.Time `json:"occurDate"`
	Category    Category  `json:"category"`
	Description *string   `json:"description"`
	Vote        []string `json:"vote"`
}

type IncomeInput struct {
	Amount      *int       `json:"amount"`
	Date        *time.Time `json:"date"`
	Category    *Category  `json:"category"`
	Description *string    `json:"description"`
}

// List current or historical portfolio
type Portfolio struct {
	Total  int       `json:"total"`
	Income []string  `json:"income"`
	Cost   []string  `json:"cost"`
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

type Category string

const (
	CategoryInvestment Category = "INVESTMENT"
	CategorySalory     Category = "SALORY"
	CategoryOthers     Category = "OTHERS"
	CategoryDaily      Category = "DAILY"
	CategoryLearning   Category = "LEARNING"
	CategoryCharity    Category = "CHARITY"
)

var AllCategory = []Category{
	CategoryInvestment,
	CategorySalory,
	CategoryOthers,
	CategoryDaily,
	CategoryLearning,
	CategoryCharity,
}

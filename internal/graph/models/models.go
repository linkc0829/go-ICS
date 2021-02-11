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
	Privacy     Privacy           `json:"privacy"`
}

func (Cost) IsPortfolio() {}

type CreateCostInput struct {
	Amount      int          `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   time.Time    `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    CostCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description string       `bson:"description,omitempty" json:"description,omitempty"`
	Privacy     Privacy      `bson:"privacy,omitempty" json:"privacy,omitempty"`
}

type CreateIncomeInput struct {
	Amount      int            `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   time.Time      `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    IncomeCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description string         `bson:"description,omitempty" json:"description,omitempty"`
	Privacy     Privacy        `bson:"privacy,omitempty" json:"privacy,omitempty"`
}

type CreatePortfolioInput interface {
	GetAmount() int
	GetOccurDate() time.Time
	GetCategory() PortfolioCategory
	GetDescription() string
	GetPrivacy() Privacy
}

func (c CreateCostInput) GetAmount() int {
	return c.Amount
}

func (c CreateCostInput) GetOccurDate() time.Time {
	return c.OccurDate
}
func (c CreateCostInput) GetCategory() PortfolioCategory {
	return PortfolioCategory(c.Category)
}
func (c CreateCostInput) GetDescription() string {
	return c.Description
}
func (c CreateCostInput) GetPrivacy() Privacy {
	return c.Privacy
}
func (c CreateIncomeInput) GetAmount() int {
	return c.Amount
}
func (c CreateIncomeInput) GetOccurDate() time.Time {
	return c.OccurDate
}
func (c CreateIncomeInput) GetCategory() PortfolioCategory {
	return PortfolioCategory(c.Category)
}
func (c CreateIncomeInput) GetDescription() string {
	return c.Description
}
func (c CreateIncomeInput) GetPrivacy() Privacy {
	return c.Privacy
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
	Privacy     Privacy           `json:"privacy"`
}

func (Income) IsPortfolio() {}

type UpdateCostInput struct {
	Amount      *int          `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   *time.Time    `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    *CostCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description *string       `bson:"description,omitempty" json:"description,omitempty"`
	Privacy     *Privacy      `bson:"privacy,omitempty" json:"privacy,omitempty"`
}

type UpdateIncomeInput struct {
	Amount      *int            `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   *time.Time      `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    *IncomeCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description *string         `bson:"description,omitempty" json:"description,omitempty"`
	Privacy     *Privacy        `bson:"privacy,omitempty" json:"privacy,omitempty"`
}

type UpdatePortfolio struct {
	Amount      *int               `bson:"amount,omitempty" json:"amount,string,omitempty"`
	OccurDate   *time.Time         `bson:"occurDate,omitempty" json:"occurDate,omitempty"`
	Category    *PortfolioCategory `bson:"category,omitempty" json:"category,omitempty"`
	Description *string            `bson:"description,omitempty" json:"description,omitempty"`
	Privacy     *Privacy           `bson:"privacy,omitempty" json:"privacy,omitempty"`
}

type UpdatePortfolioInput interface {
	GetAmount() *int
	GetOccurDate() *time.Time
	GetCategory() *PortfolioCategory
	GetDescription() *string
	GetPrivacy() *Privacy
}

func (c UpdateCostInput) GetAmount() *int {
	return c.Amount
}
func (c UpdateCostInput) GetOccurDate() *time.Time {
	return c.OccurDate
}
func (c UpdateCostInput) GetCategory() *PortfolioCategory {
	return (*PortfolioCategory)(c.Category)
}
func (c UpdateCostInput) GetDescription() *string {
	return c.Description
}
func (c UpdateCostInput) GetPrivacy() *Privacy {
	return c.Privacy
}
func (c UpdateIncomeInput) GetAmount() *int {
	return c.Amount
}
func (c UpdateIncomeInput) GetOccurDate() *time.Time {
	return c.OccurDate
}
func (c UpdateIncomeInput) GetCategory() *PortfolioCategory {
	return (*PortfolioCategory)(c.Category)
}
func (c UpdateIncomeInput) GetDescription() *string {
	return c.Description
}
func (c UpdateIncomeInput) GetPrivacy() *Privacy {
	return c.Privacy
}

func (UpdateCostInput) IsPortfolioInput()   {}
func (UpdateIncomeInput) IsPortfolioInput() {}

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

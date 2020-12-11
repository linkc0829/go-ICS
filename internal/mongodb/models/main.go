package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserModel struct
type UserModel struct {
	ID       	primitive.ObjectID		`bson:"_id" json:"_id,omitempty"`
	UserId  	string             		`bson:"userid" json:"username,omitempty"`
	Email     	string             		`bson:"email" json:"email,omitempty"`
	Nickname  	string  				`bson:"nickname" json:"nickname,omitempty"`
	CreatedAt 	time.Time          		`bson:"createAt" json:"createAt,omitempty"`
	Friends   	[]primitive.ObjectID 	`bson:"friends" json:"friends,omitempty"`
	Token     	string             		`bson:"token" json:"token,omitempty"` // graphql only
	LastQuery 	time.Time 		 		`bson:"lastQuery" json:"lastQuery, omitempty"`
}

type category string

const (
	CategoryInvestment category = "INVESTMENT"
	CategorySalory     category = "SALORY"
	CategoryOthers     category = "OTHERS"
	CategoryDaily      category = "DAILY"
	CategoryLearning   category = "LEARNING"
	CategoryCharity    category = "CHARITY"
)

type IncomeModel struct{
	ID 			primitive.ObjectID		`bson:"_id" json:"_id,omitempty"`
	Owner		primitive.ObjectID		`bson:"owner" json:"owner,omitempty"`
	Amount		int 					`bson:"amount" json:"amount,omitempty"`
	OccurDate	time.Time 				`bson:"date" json:"date,omitempty"`
	Category	category				`bson:"category" json:"category,omitempty"`
	Description	string  				`bson:"description" json:"description,omitempty"`
	Like		[]primitive.ObjectID  	`bson:"like" json:"like,omitempty"`

}

type CostModel struct{
	ID 			primitive.ObjectID		`bson:"_id" json:"_id,omitempty"`
	Owner		primitive.ObjectID		`bson:"owner" json:"owner,omitempty"`
	Amount		int 					`bson:"amount" json:"amount,omitempty"`
	OccurDate	time.Time 				`bson:"date" json:"date,omitempty"`
	Category	category				`bson:"category" json:"category,omitempty"`
	Description	string  				`bson:"description" json:"description,omitempty"`
	Like		[]primitive.ObjectID  	`bson:"like" json:"like,omitempty"`
}
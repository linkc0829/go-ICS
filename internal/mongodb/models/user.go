package models

// UserModel struct
type UserModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Username  string             `bson:"username" json:"username,omitempty"`
	Nickname  string  			 `bson:"nickname" json:"nickname"`
	Password  string             `bson:"password" json:"password,omitempty"`
	Email     string             `bson:"email" json:"email,omitempty"`
	CreatedAt time.Time          `bson:"createAt" json:"createAt,omitempty"`
	Friends   primitive.ObjectID `bson:"friends" json:"friends"`
	Token     string             `bson:"token" json:"token,omitempty"` // graphql only
}
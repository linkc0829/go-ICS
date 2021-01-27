package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkc0829/go-icsharing/internal/db/mongodb"
	"github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	"github.com/linkc0829/go-icsharing/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserProfileHandler(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		data := struct {
			Title string
		}{
			Title: "User Profile | ICS",
		}

		c.HTML(http.StatusOK, "profile", data)

	}
}

func UserHistoryHandler(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		data := struct {
			Title string
		}{
			Title: "User History | ICS",
		}

		c.HTML(http.StatusOK, "history", data)

	}
}

//UserFriendsHandlers handle GET /friends/:id
func UserFriendsHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {

		type friend struct {
			Title string
			User  []string
			Info  [][]string
		}

		id := c.Param("id")
		hexID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		q := bson.M{"_id": hexID}
		user := models.UserModel{}
		if err := db.Users.FindOne(c, q).Decode(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		friends := []string{}
		infos := [][]string{}
		for _, f_id := range user.Friends {
			f := f_id.Hex()
			friends = append(friends, f)
			fInfo := models.UserModel{}
			if err := db.Users.FindOne(c, bson.M{"_id": f_id}).Decode(&fInfo); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			info := []string{fInfo.UserID, *fInfo.NickName, fInfo.Email}
			infos = append(infos, info)
		}

		data := friend{
			Title: user.UserID + "'s Friends List",
			User:  friends,
			Info:  infos,
		}
		c.HTML(http.StatusOK, "relationship", data)
	}
}

func UserFollowersHandler(cfg *utils.ServerConfig, db *mongodb.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {

		type friend struct {
			Title string
			User  []string
			Info  [][]string
		}

		id := c.Param("id")
		hexID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		q := bson.M{"_id": hexID}
		user := models.UserModel{}
		if err := db.Users.FindOne(c, q).Decode(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		followers, err := getUserFollowers(c, db, hexID)
		infos := [][]string{}
		for _, f_id := range followers {
			fInfo := models.UserModel{}
			if err := db.Users.FindOne(c, bson.M{"_id": f_id}).Decode(&fInfo); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			info := []string{fInfo.UserID, *fInfo.NickName, fInfo.Email}
			infos = append(infos, info)
		}

		data := friend{
			Title: user.UserID + "'s Friends List",
			User:  followers,
			Info:  infos,
		}
		c.HTML(http.StatusOK, "relationship", data)

	}
}

func getUserFollowers(ctx *gin.Context, DB *mongodb.MongoDB, id primitive.ObjectID) (followers []string, err error) {
	//find users that have me as friend
	q := bson.M{"friends": id}
	cursor, err := DB.Users.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		f := models.UserModel{}
		bsonBytes, _ := bson.Marshal(result)
		bson.Unmarshal(bsonBytes, &f)
		follower := f.ID.Hex()
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	return
}

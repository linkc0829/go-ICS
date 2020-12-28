package gqlclient

import (
	//"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"

	"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/pkg/utils"
)

//GetUser handle request GET /user/:id
func GetUser(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		ID := c.Param("id")
		var query struct {
			GetUser struct {
				Id        graphql.ID
				UserId    graphql.String
				Email     graphql.String
				NickName  graphql.String
				CreatedAt graphql.String
				Friends   []struct {
					Id graphql.ID
				}
				Followers []struct {
					Id graphql.ID
				}
			} `graphql:"getUser(id: $ID)"`
		}
		variables := map[string]interface{}{
			"ID": graphql.ID(ID),
		}
		err := client.Query(c, &query, variables)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		log.Println(query)
		c.JSON(http.StatusOK, query)
	}
}

//CreateUser handle request POST /user
func CreateUser(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($createUserInput: CreateUserInput!){
				createUser(input: $createUserInput){
					id
					userId
					email
					nickName
					createdAt
					friends{
						id
					}
					followers{
						id
					}
				}
			}
		*/
		var mutation struct {
			CreateUser struct {
				Id        graphql.ID
				UserId    graphql.String
				Email     graphql.String
				NickName  graphql.String
				CreatedAt graphql.String
				Friends   []struct {
					Id graphql.ID
				}
				Followers []struct {
					Id graphql.ID
				}
			} `graphql:"createUser(input: $createUserInput)"`
		}
		nickName := c.PostForm("nickName")
		variables := map[string]interface{}{
			"createUserInput": models.CreateUserInput{
				UserID:   c.PostForm("userID"),
				Email:    c.PostForm("email"),
				NickName: &nickName,
			},
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
	}
}

//UpdateUser handle request PATCH /user/:id
func UpdateUser(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)

		/*
			mutation ($id: ID!, $updateUserInput: UpdateUserInput!){
				updateUser(id: $id, input: $updateUserInput){
					id
					userId
					email
					nickName
					createdAt
					friends{
						id
					}
					followers{
						id
					}
				}
			}
		*/
		var mutation struct {
			UpdateUser struct {
				Id        graphql.ID
				UserId    graphql.String
				Email     graphql.String
				NickName  graphql.String
				CreatedAt graphql.String
				Friends   []struct {
					Id graphql.ID
				}
				Followers []struct {
					Id graphql.ID
				}
			} `graphql:"updateUser(id: $id, input: $updateUserInput)"`
		}
		id := c.Param("id")
		nickName := c.PostForm("nickName")
		userId := c.PostForm("userID")
		email := c.PostForm("email")
		variables := map[string]interface{}{
			"id": id,
			"updateUserInput": models.UpdateUserInput{
				UserID:   &userId,
				Email:    &email,
				NickName: &nickName,
			},
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
	}
}

//DeleteUser handle request DELETE /user/:id
func DeleteUser(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//GetUser handle request PUT /user/addfriend/:id
func AddFriend(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

//GetUserIncome handle request GET /user/:id/income/
func GetUserIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		client := newClient(c, cfg)
		ID := c.Param("id")
		var query struct {
			GetUserIncome []struct {
				Id    graphql.ID
				Owner struct {
					Id graphql.ID
				}
				Amount      graphql.Int
				Category    graphql.String
				OccurDate   graphql.String
				Description graphql.String
				Vote        []struct {
					Id graphql.ID
				}
			} `graphql:"getUserIncome(id: $ID)"`
		}
		variables := map[string]interface{}{
			"ID": graphql.ID(ID),
		}
		err := client.Query(c, &query, variables)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		log.Println(query)
		c.JSON(http.StatusOK, query)

	}
}

//GetUserCost handle request GET /user/:id/cost
func GetUserCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Not implement yet.")
	}
}

package gqlclient

import (
	//"encoding/json"

	"net/http"
	"strconv"

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
				Role graphql.String
			} `graphql:"getUser(id: $ID)"`
		}
		variables := map[string]interface{}{
			"ID": graphql.ID(ID),
		}
		err := client.Query(c, &query, variables)
		if err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
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
				Roles graphql.String
			} `graphql:"createUser(input: $createUserInput)"`
		}
		createUserInput := models.CreateUserInput{}
		err := c.ShouldBind(&createUserInput)
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}
		variables := map[string]interface{}{
			"createUserInput": createUserInput,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
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
				Role graphql.String
			} `graphql:"updateUser(id: $id, input: $updateUserInput)"`
		}
		id := c.Param("id")
		updateUserInput := models.UpdateUserInput{}
		err := c.ShouldBind(&updateUserInput)
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}

		variables := map[string]interface{}{
			"id":              id,
			"updateUserInput": updateUserInput,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
		return
	}
}

//DeleteUser handle request DELETE /user/:id
func DeleteUser(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!){
				deleteUser(id: $id)
			}
		*/
		var mutation struct {
			DeleteUser graphql.Boolean `graphql:"deleteUser(id: $id)"`
		}

		id := c.Param("id")
		variables := map[string]interface{}{
			"id": id,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
		return
	}
}

//GetUser handle request PUT /user/addfriend/:id
func AddFriend(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!){
				addFriend(id: $id)
			}
		*/
		var mutation struct {
			AddFriend graphql.Boolean `graphql:"addFriend(id: $id)"`
		}
		id := c.Param("id")
		variables := map[string]interface{}{
			"id": id,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
		return
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
				Privacy graphql.String
			} `graphql:"getUserIncome(id: $ID)"`
		}
		variables := map[string]interface{}{
			"ID": graphql.ID(ID),
		}
		err := client.Query(c, &query, variables)
		if err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, query)

	}
}

//GetUserCost handle request GET /user/:id/cost
func GetUserCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		ID := c.Param("id")
		var query struct {
			GetUserCost []struct {
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
				Privacy graphql.String
			} `graphql:"getUserCost(id: $ID)"`
		}
		variables := map[string]interface{}{
			"ID": graphql.ID(ID),
		}
		err := client.Query(c, &query, variables)
		if err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, query)
	}
}

//GetUserIncomeHistory handle request GET /user/:id/income/history?range=
func GetUserIncomeHistory(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		client := newClient(c, cfg)
		ID := c.Param("id")
		Range, err := strconv.Atoi(c.Query("range"))
		if err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		var query struct {
			GetUserIncomeHistory []struct {
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
				Privacy graphql.String
			} `graphql:"getUserIncomeHistory(id: $ID, range: $Range)"`
		}
		variables := map[string]interface{}{
			"ID":    graphql.ID(ID),
			"Range": graphql.Int(Range),
		}
		err = client.Query(c, &query, variables)
		if err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, query)

	}
}

//GetUserCostHistory handle request GET /user/:id/income/history?range=
func GetUserCostHistory(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		client := newClient(c, cfg)
		ID := c.Param("id")
		Range, err := strconv.Atoi(c.Query("range"))
		if err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		var query struct {
			GetUserCostHistory []struct {
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
				Privacy graphql.String
			} `graphql:"getUserCostHistory(id: $ID, range: $Range)"`
		}
		variables := map[string]interface{}{
			"ID":    graphql.ID(ID),
			"Range": graphql.Int(Range),
		}
		err = client.Query(c, &query, variables)
		if err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, query)

	}
}

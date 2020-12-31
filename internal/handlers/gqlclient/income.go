package gqlclient

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"

	"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/pkg/utils"
)

//"encoding/json"
//"github.com/linkc0829/go-ics/internal/graph/models"

//CreateIncome handle request POST /income
func CreateIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($createIncomeInput: CreateIncomeInput!){
				createIncome(input: $createIncomeInput){
					id
					owner{
						id
					}
					amount
					category
					occurDate
					description
					vote{
						id
					}
				}
			}
		*/
		var mutation struct {
			CreateIncome struct {
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
			} `graphql:"createIncome(input: $createIncomeInput)"`
		}
		createIncomeInput := models.CreateIncomeInput{}
		err := c.ShouldBind(&createIncomeInput)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		variables := map[string]interface{}{
			"createIncomeInput": createIncomeInput,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
	}
}

//UpdateIncome handle request POST /income/:id
func UpdateIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!, $updateIncomeInput: updateIncomeInput!){
				updateIncome(id: $id, input: $updateIncomeInput){
					id
					owner{
						id
					}
					amount
					category
					occurDate
					description
					vote{
						id
					}
				}
			}
		*/
		var mutation struct {
			UpdateIncome struct {
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
			} `graphql:"updateIncome(id: $id, input: $updateIncomeInput)"`
		}
		id := c.Param("id")
		UpdateIncomeInput := models.UpdateIncomeInput{}
		err := c.ShouldBind(&UpdateIncomeInput)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		variables := map[string]interface{}{
			"id":                id,
			"updateIncomeInput": UpdateIncomeInput,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
	}
}

//DeleteIncome handle request POST /income
func DeleteIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!){
				deleteIncome(id: $id)
			}
		*/
		var mutation struct {
			DeleteIncome graphql.Boolean `graphql:"deleteIncome(id: $id)"`
		}

		id := c.Param("id")
		variables := map[string]interface{}{
			"id": id,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
		return
	}
}

//VoteIncome handle request PUT /income/vote/:id
func VoteIncome(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!){
				VoteIncome(id: $id)
			}
		*/
		var mutation struct {
			VoteIncome graphql.Int `graphql:"voteIncome(id: $id)"`
		}
		id := c.Param("id")
		variables := map[string]interface{}{
			"id": id,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
		return
	}
}

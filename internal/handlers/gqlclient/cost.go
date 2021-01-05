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

//CreateCost handle request POST /cost
func CreateCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($createCostInput: CreateCostInput!){
				createCost(input: $createCostInput){
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
			CreateCost struct {
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
			} `graphql:"createCost(input: $createCostInput)"`
		}
		createCostInput := models.CreateCostInput{}
		err := c.ShouldBind(&createCostInput)
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}
		variables := map[string]interface{}{
			"createCostInput": createCostInput,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
	}
}

//UpdateCost handle request POST /cost/:id
func UpdateCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!, $updateCostInput: updateCostInput!){
				updateCost(id: $id, input: $updateCostInput){
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
			UpdateCost struct {
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
			} `graphql:"updateCost(id: $id, input: $updateCostInput)"`
		}
		id := c.Param("id")
		UpdateCostInput := models.UpdateCostInput{}
		err := c.ShouldBind(&UpdateCostInput)
		if err != nil {
			ErrorWriter(c, http.StatusInternalServerError, err)
			return
		}
		variables := map[string]interface{}{
			"id":              id,
			"updateCostInput": UpdateCostInput,
		}

		if err := client.Mutate(c, &mutation, variables); err != nil {
			ErrorWriter(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, mutation)
	}
}

//DeleteCost handle request POST /cost
func DeleteCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!){
				deleteCost(id: $id)
			}
		*/
		var mutation struct {
			DeleteCost graphql.Boolean `graphql:"deleteCost(id: $id)"`
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

//VoteCost handle request PUT /cost/vote/:id
func VoteCost(cfg *utils.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := newClient(c, cfg)
		/*
			mutation ($id: ID!){
				VoteCost(id: $id)
			}
		*/
		var mutation struct {
			VoteCost graphql.Int `graphql:"voteCost(id: $id)"`
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

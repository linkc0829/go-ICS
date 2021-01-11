package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/linkc0829/go-ics/internal/db/sqlitedb"
	"github.com/linkc0829/go-ics/pkg/utils"
)

func FreeTrialHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		data := struct {
			Title        string
			ErrorMessage string
		}{
			Title:        "Free try Income & Cost SFunction here. Signup to share with friends.",
			ErrorMessage: "",
		}

		c.HTML(http.StatusOK, "layout", data)
	}

}

//GetPortfolioHandlers handle GET request from /api/v1/trial
func GetPortfolioHandlers(cfg *utils.ServerConfig, sqlite *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := sqlite.Table("portfolio").Select("*").Order("occurDate").Rows()
		if err != nil {
			c.AbortWithError(http.StatusBadGateway, errors.New("Unable to retrieve database"))
			return
		}
		portfolios := []sqlitedb.Portfolio{}

		for rows.Next() {
			p := sqlitedb.Portfolio{}
			sqlite.ScanRows(rows, &p)
			portfolios = append(portfolios, p)
		}
		c.JSON(http.StatusOK, gin.H{"Portfolios": portfolios})

	}
}

//GetPortfolioHandlers handle POST request from /api/v1/trial
func CreatePortfolioHandlers(cfg *utils.ServerConfig, sqlite *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := sqlitedb.Portfolio{}
		err := c.ShouldBind(&input)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		log.Println(input)
		sqlite.Table("portfolio").Create(&input)

		output := sqlitedb.Portfolio{}
		sqlite.Table("portfolio").Last(&output)
		c.JSON(http.StatusOK, gin.H{"portfolio": output})
	}
}

//GetPortfolioHandlers handle PATCH request from /api/v1/trial/:id
func UpdatePortfolioHandlers(cfg *utils.ServerConfig, sqlite *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := sqlitedb.Portfolio{}
		err := c.ShouldBind(&input)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		uintID := uint(id)
		sqlite.Table("portfolio").Model(&sqlitedb.Portfolio{ID: &uintID}).Updates(input)
		portfolio := sqlitedb.Portfolio{}
		sqlite.Table("portfolio").First(&portfolio, id)
		json, _ := json.Marshal(portfolio)
		c.JSON(http.StatusOK, json)
	}
}

func DeletePortfolioHandlers(cfg *utils.ServerConfig, sqlite *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		portfolio := sqlitedb.Portfolio{}
		sqlite.Table("portfolio").Delete(&portfolio, id)
		c.JSON(http.StatusOK, gin.H{"data": "True"})
	}
}

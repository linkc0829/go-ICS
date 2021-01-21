package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/linkc0829/go-ics/internal/db/mongodb"
	dbModel "github.com/linkc0829/go-ics/internal/db/mongodb/models"
	"github.com/linkc0829/go-ics/internal/db/redisdb"
	"github.com/linkc0829/go-ics/internal/db/sqlitedb"
	"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/pkg/server"
	"github.com/linkc0829/go-ics/pkg/utils/datasource"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	testUser           user
	admin              user
	createIncomeInputs []models.CreateIncomeInput
	createCostInputs   []models.CreateCostInput
	myIncomes          []portfolio
	myCosts            []portfolio
	incomeIDs          map[string]bool
	costIDs            map[string]bool
	N                  int = 1
)

type portfolio struct {
	Id    string
	Owner struct {
		Id string
	}
	Amount      int
	Category    string
	OccurDate   time.Time
	Description string
	Vote        []struct {
		Id string
	}
	Privacy string
}

type user struct {
	ID     string
	UserID string
	Email  string
	PWD    string
}

type token struct {
	Token       string `json:"token"`
	TokenExpiry string `json:"token_expiry"`
	Type        string `json:"type"`
}

//init mock data
func init() {
	for i := 0; i < N; i++ {
		add, _ := time.ParseDuration(strconv.Itoa(i+1) + "h")
		createIncomeInput := models.CreateIncomeInput{
			Amount:      1111 * (i + 1),
			OccurDate:   time.Now().Add(add).Truncate(time.Second),
			Category:    models.IncomeCategoryInvestment,
			Description: "Income no." + strconv.Itoa(i),
			Privacy:     models.PrivacyFriend,
		}
		createIncomeInputs = append(createIncomeInputs, createIncomeInput)

		createCostInput := models.CreateCostInput{
			Amount:      1111 * (i + 1),
			OccurDate:   time.Now().Add(add).Truncate(time.Second),
			Category:    models.CostCategoryInvestment,
			Description: "Cost no." + strconv.Itoa(i),
			Privacy:     models.PrivacyFriend,
		}
		createCostInputs = append(createCostInputs, createCostInput)
	}
	incomeIDs = make(map[string]bool)
	costIDs = make(map[string]bool)

	testUser.UserID = "test99"
	testUser.Email = "test999@gmail.com"
	testUser.PWD = "123456"

	admin.UserID = "admin_test"
	admin.Email = "admin_test@icsharing.com"
	admin.PWD = "123456"

}

func TestGraphQLAPI(t *testing.T) {

	mongoDB := mongodb.ConnectMongoDB(serverconf)
	defer mongodb.CloseMongoDB(mongoDB)
	log.Println("Mongo connected")

	sqlite := sqlitedb.ConnectSqlite()
	defer sqlitedb.CloseSqlite(sqlite)

	redis := redisdb.ConnectRedis(serverconf)
	defer redisdb.CloseRedis(redis)
	log.Println("Redis connected")

	db := &datasource.DB{
		Mongo:  mongoDB,
		Sqlite: sqlite,
		Redis:  redis,
	}

	r := server.SetupServer(serverconf, db)
	ts := httptest.NewTLSServer(r)
	defer ts.Close()

	log.Println("Test on: ", ts.URL)

	client := ts.Client()

	t.Run("create admin user", func(t *testing.T) {
		hashPWD, err := bcrypt.GenerateFromPassword([]byte(admin.PWD), 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		pwd := string(hashPWD)
		id := primitive.NewObjectID()
		admin.ID = id.Hex()

		newUser := &dbModel.UserModel{
			ID:              id,
			UserID:          admin.UserID,
			Password:        &pwd,
			Email:           admin.Email,
			NickName:        &admin.UserID,
			CreatedAt:       time.Now(),
			LastIncomeQuery: time.Now(),
			LastCostQuery:   time.Now(),
			Provider:        "ics",
			Role:            dbModel.ADMIN,
		}
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		_, err = mongoDB.Users.InsertOne(ctx, newUser)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		test := &dbModel.UserModel{}
		mongoDB.Users.FindOne(ctx, bson.M{"_id": newUser.ID}).Decode(test)

	})

	var refresh_token *http.Cookie
	tokenJson := &token{}

	t.Run("user signup", func(t *testing.T) {
		data := url.Values{"userID": {testUser.UserID}, "email": {testUser.Email}, "password": {testUser.PWD}}
		reader := strings.NewReader(data.Encode())
		req, err := http.NewRequest("POST", ts.URL+"/api/v1/auth/ics/signup", reader)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("accept", "json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		refresh_token = resp.Cookies()[0]

	})

	t.Run("user repeated signup", func(t *testing.T) {
		data := url.Values{"userID": {testUser.UserID}, "email": {testUser.Email}, "password": {testUser.PWD}}
		reader := strings.NewReader(data.Encode())
		req, err := http.NewRequest("POST", ts.URL+"/api/v1/auth/ics/signup", reader)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("accept", "json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 400 {
			t.Fatalf("Expected status code 400, got %v", resp.StatusCode)
		}
	})

	t.Run("token soft refresh", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL+"/api/v1/auth/ics/refresh_token", nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		req.Header.Set("accept", "json")
		req.AddCookie(refresh_token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		//decode resp body
		err = json.NewDecoder(resp.Body).Decode(&tokenJson)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if refresh_token.Value == resp.Cookies()[0].Value {
			t.Fatal("refresh token should update after soft resresh")
		}
	})

	for _, tt := range createIncomeInputs {
		createIncomeInput := tt
		t.Run("create_"+createIncomeInput.Description, func(t *testing.T) {
			//t.Parallel()
			query := `mutation ($createIncomeInput: CreateIncomeInput!){
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
					privacy
				}
			}`
			variables := map[string]interface{}{
				"createIncomeInput": createIncomeInput,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			//send := map[string]string{"Query": query, "Variables": income}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			out := struct {
				Data struct {
					CreateIncome portfolio
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			//not work on drone
			// outCheck := models.CreateIncomeInput{
			// 	Amount:      out.Data.CreateIncome.Amount,
			// 	Category:    models.IncomeCategory(out.Data.CreateIncome.Category),
			// 	Description: out.Data.CreateIncome.Description,
			// 	Privacy:     models.Privacy(out.Data.CreateIncome.Privacy),
			// 	OccurDate:   out.Data.CreateIncome.OccurDate.Truncate(time.Second),
			// }
			// log.Println(outCheck)
			// log.Println(createIncomeInput)
			// if outCheck != createIncomeInput {
			// 	t.Errorf("input output mismatch")
			// }
			incomeIDs[out.Data.CreateIncome.Id] = true
		})
	}
	t.Run("getMyIncome", func(t *testing.T) {
		query := `{
			myIncome{
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
				privacy
			}
		}`
		send := struct {
			Query string `json:"query"`
		}{
			Query: query,
		}
		//send := map[string]string{"Query": query, "Variables": income}
		sendJson, err := json.Marshal(send)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		tokenString := tokenJson.Type + " " + tokenJson.Token
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("accept", "json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		out := struct {
			Data struct {
				MyIncome []portfolio
			}
		}{}
		err = json.Unmarshal(body, &out)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		for _, income := range out.Data.MyIncome {
			if !incomeIDs[income.Id] {
				t.Errorf("get wrong user income")
			}
		}
		myIncomes = append(myIncomes, out.Data.MyIncome...)
		testUser.ID = out.Data.MyIncome[0].Owner.Id
	})

	for _, income := range myIncomes {
		income := income
		t.Run("update"+income.Description, func(t *testing.T) {
			query := `mutation ($id: ID!, $updateIncomeInput: UpdateIncomeInput!){
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
					privacy
				}
			}`
			category := models.IncomeCategorySalary
			updateIncomeInput := models.UpdateIncomeInput{
				Category: &category,
			}
			variables := map[string]interface{}{
				"id":                income.Id,
				"updateIncomeInput": updateIncomeInput,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			//send := map[string]string{"Query": query, "Variables": income}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			out := struct {
				Data struct {
					UpdateIncome portfolio
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if out.Data.UpdateIncome.Category != string(category) {
				t.Errorf("input output mismatch")
			}
		})
	}

	for _, income := range myIncomes {
		income := income
		t.Run("vote "+income.Description, func(t *testing.T) {
			query := `mutation ($id: ID!){
					voteIncome(id: $id)
				}`
			variables := map[string]interface{}{
				"id": income.Id,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			//log.Println(string(body))
			out := struct {
				Data struct {
					VoteIncome int
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if out.Data.VoteIncome == 0 {
				t.Errorf("Vote " + income.Description + " failed")
			}

			//second vote should reset vote
			resp, err = client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if out.Data.VoteIncome == 1 {
				t.Errorf("Vote " + income.Description + " failed")
			}
		})
	}

	for _, income := range myIncomes {
		income := income
		t.Run("delete "+income.Description, func(t *testing.T) {
			query := `mutation ($id: ID!){
					deleteIncome(id: $id)
				}`

			variables := map[string]interface{}{
				"id": income.Id,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			//send := map[string]string{"Query": query, "Variables": income}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			out := struct {
				Data struct {
					DeleteIncome bool
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if out.Data.DeleteIncome == false {
				t.Errorf("Delete " + income.Description + " failed")
			}
		})
	}
	for _, tt := range createCostInputs {
		createCostInput := tt
		t.Run("create_"+createCostInput.Description, func(t *testing.T) {
			//t.Parallel()
			query := `mutation ($createCostInput: CreateCostInput!){
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
					privacy
				}
			}`
			variables := map[string]interface{}{
				"createCostInput": createCostInput,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			//send := map[string]string{"Query": query, "Variables": cost}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			out := struct {
				Data struct {
					CreateCost portfolio
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			//not work on drone
			// outCheck := models.CreateCostInput{
			// 	Amount:      out.Data.CreateCost.Amount,
			// 	Category:    models.CostCategory(out.Data.CreateCost.Category),
			// 	Description: out.Data.CreateCost.Description,
			// 	Privacy:     models.Privacy(out.Data.CreateCost.Privacy),
			// 	OccurDate:   out.Data.CreateCost.OccurDate.Truncate(time.Second),
			// }
			// log.Println(outCheck)
			// log.Println(createCostInput)
			// if outCheck != createCostInput {
			// 	t.Errorf("input output mismatch")
			// }
			costIDs[out.Data.CreateCost.Id] = true
		})
	}
	t.Run("getMyCost", func(t *testing.T) {
		query := `{
			myCost{
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
				privacy
			}
		}`
		send := struct {
			Query string `json:"query"`
		}{
			Query: query,
		}
		//send := map[string]string{"Query": query, "Variables": cost}
		sendJson, err := json.Marshal(send)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		tokenString := tokenJson.Type + " " + tokenJson.Token
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("accept", "json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		out := struct {
			Data struct {
				MyCost []portfolio
			}
		}{}
		err = json.Unmarshal(body, &out)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		for _, cost := range out.Data.MyCost {
			if !costIDs[cost.Id] {
				t.Errorf("get wrong user cost")
			}
		}
		myCosts = append(myCosts, out.Data.MyCost...)
		testUser.ID = out.Data.MyCost[0].Owner.Id
	})

	for _, cost := range myCosts {
		cost := cost
		t.Run("update"+cost.Description, func(t *testing.T) {
			query := `mutation ($id: ID!, $updateCostInput: UpdateCostInput!){
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
					privacy
				}
			}`
			category := models.CostCategoryCharity
			updateCostInput := models.UpdateCostInput{
				Category: &category,
			}
			variables := map[string]interface{}{
				"id":              cost.Id,
				"updateCostInput": updateCostInput,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			//send := map[string]string{"Query": query, "Variables": cost}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			out := struct {
				Data struct {
					UpdateCost portfolio
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if out.Data.UpdateCost.Category != string(category) {
				t.Errorf("input output mismatch")
			}
		})
	}

	for _, cost := range myCosts {
		cost := cost
		t.Run("vote "+cost.Description, func(t *testing.T) {
			query := `mutation ($id: ID!){
					voteCost(id: $id)
				}`
			variables := map[string]interface{}{
				"id": cost.Id,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			//log.Println(string(body))
			out := struct {
				Data struct {
					VoteCost int
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if out.Data.VoteCost == 0 {
				t.Errorf("Vote " + cost.Description + " failed")
			}

			//second vote should reset vote
			resp, err = client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if out.Data.VoteCost == 1 {
				t.Errorf("Vote " + cost.Description + " failed")
			}
		})
	}

	for _, cost := range myCosts {
		cost := cost
		t.Run("delete "+cost.Description, func(t *testing.T) {
			query := `mutation ($id: ID!){
					deleteCost(id: $id)
				}`

			variables := map[string]interface{}{
				"id": cost.Id,
			}
			send := struct {
				Query     string                 `json:"query"`
				Variables map[string]interface{} `json:"variables"`
			}{
				Query:     query,
				Variables: variables,
			}
			//send := map[string]string{"Query": query, "Variables": cost}
			sendJson, err := json.Marshal(send)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, err := http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			tokenString := tokenJson.Type + " " + tokenJson.Token
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("accept", "json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.StatusCode != 200 {
				t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			out := struct {
				Data struct {
					DeleteCost bool
				}
			}{}
			err = json.Unmarshal(body, &out)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if out.Data.DeleteCost == false {
				t.Errorf("Delete " + cost.Description + " failed")
			}
		})
	}

	t.Run("test login and ADMIN Delete Test User", func(t *testing.T) {

		//ADMIN login
		data := url.Values{"userID": {admin.UserID}, "email": {admin.Email}, "password": {admin.PWD}}
		reader := strings.NewReader(data.Encode())
		req, err := http.NewRequest("POST", ts.URL+"/api/v1/auth/ics/login", reader)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("accept", "json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		refresh_token := resp.Cookies()[0]

		req, err = http.NewRequest("GET", ts.URL+"/api/v1/auth/ics/refresh_token", nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		req.Header.Set("accept", "json")
		req.AddCookie(refresh_token)
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		//decode resp body
		err = json.NewDecoder(resp.Body).Decode(&tokenJson)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if refresh_token.Value == resp.Cookies()[0].Value {
			t.Fatal("refresh token should update after soft resresh")
		}

		//delete user
		query := `mutation ($id: ID!){
			deleteUser(id: $id)
		}`

		variables := map[string]interface{}{
			"id": testUser.ID,
		}
		send := struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}{
			Query:     query,
			Variables: variables,
		}
		sendJson, err := json.Marshal(send)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		req, err = http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		tokenString := tokenJson.Type + " " + tokenJson.Token
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("accept", "json")

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		out := struct {
			Data struct {
				DeleteUser bool
			}
		}{}
		err = json.Unmarshal(body, &out)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if out.Data.DeleteUser == false {
			t.Errorf("Delete " + testUser.UserID + " failed")
		}
		//delete test admin
		variables = map[string]interface{}{
			"id": admin.ID,
		}
		send = struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}{
			Query:     query,
			Variables: variables,
		}
		sendJson, err = json.Marshal(send)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		req, err = http.NewRequest("POST", ts.URL+"/api/v1/graph", bytes.NewBuffer(sendJson))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		tokenString = tokenJson.Type + " " + tokenJson.Token
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenString)
		req.Header.Set("accept", "json")

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		out = struct {
			Data struct {
				DeleteUser bool
			}
		}{}
		err = json.Unmarshal(body, &out)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if out.Data.DeleteUser == false {
			t.Errorf("Delete " + admin.UserID + " failed")
		}

	})

}

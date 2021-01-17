package main

import (
	"bytes"
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
	"github.com/linkc0829/go-ics/internal/db/redisdb"
	"github.com/linkc0829/go-ics/internal/db/sqlitedb"
	"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/pkg/server"
	"github.com/linkc0829/go-ics/pkg/utils/datasource"
)

var (
	createIncomeInputs []models.CreateIncomeInput
	createCostInputs   []models.CreateCostInput
	myIncomes          []portfolio
	myCosts            []portfolio
	incomeIDs          map[string]bool
	costIDs            map[string]bool
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

//init mock data
func init() {
	for i := 0; i < 10; i++ {
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
}

func TestGraphQLAPI(t *testing.T) {

	mongoDB := mongodb.ConnectMongoDB(serverconf)
	defer mongodb.CloseMongoDB(mongoDB)

	sqlite := sqlitedb.ConnectSqlite()
	defer sqlitedb.CloseSqlite(sqlite)

	redis := redisdb.ConnectRedis(serverconf)
	defer redisdb.CloseRedis(redis)

	db := &datasource.DB{
		Mongo:  mongoDB,
		Sqlite: sqlite,
		Redis:  redis,
	}

	r := server.SetupServer(serverconf, db)
	ts := httptest.NewServer(r)
	defer ts.Close()

	log.Println("Test on: ", ts.URL)

	client := ts.Client()
	//test login/signup
	data := url.Values{"userID": {"test9999"}, "email": {"test9999@ics.com"}, "password": {"123456"}}
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

	//test soft refresh
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
	type token struct {
		Token       string `json:"token"`
		TokenExpiry string `json:"token_expiry"`
		Type        string `json:"type"`
	}
	tokenJson := &token{}
	err = json.NewDecoder(resp.Body).Decode(&tokenJson)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if refresh_token.Value == resp.Cookies()[0].Value {
		t.Fatal("refresh token should update after soft resresh")
	}
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
			outCheck := models.CreateIncomeInput{
				Amount:      out.Data.CreateIncome.Amount,
				Category:    models.IncomeCategory(out.Data.CreateIncome.Category),
				Description: out.Data.CreateIncome.Description,
				Privacy:     models.Privacy(out.Data.CreateIncome.Privacy),
				OccurDate:   out.Data.CreateIncome.OccurDate.Truncate(time.Second),
			}
			log.Println(outCheck)
			log.Println(createIncomeInput)
			if outCheck != createIncomeInput {
				t.Errorf("input output mismatch")
			}
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
			log.Println(string(body))
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

}

func TestAuthAPI(t *testing.T) {

}

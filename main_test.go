package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestRestAPI(t *testing.T) {

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

	client := &http.Client{}
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

	t.Run("create income", func(t *testing.T) {
		add, _ := time.ParseDuration("36h")
		income := models.CreateIncomeInput{
			Amount:      8888,
			OccurDate:   time.Now().Add(add),
			Category:    models.IncomeCategoryInvestment,
			Description: "create income testing",
			Privacy:     models.PrivacyFriend,
		}
		data, err := json.Marshal(income)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		req, err := http.NewRequest("POST", ts.URL+"/api/v1/income", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenJson.Token)
		req.Header.Set("accept", "json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
		}
		incomeJson := &models.Income{}
		err = json.NewDecoder(resp.Body).Decode(&incomeJson)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("update income", func(t *testing.T) {

	})

	t.Run("vote income", func(t *testing.T) {

	})

	t.Run("delete income", func(t *testing.T) {

	})

	// content, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Fatalf("Expected no error, got %v", err)
	// }
	// fmt.Println(string(content))

	// client := &http.Client{}
	// req, _ = http.NewRequest("POST", ts.URL+"/api/v1/auth/ics/signup", reader)

	// res, err := client.Get(ts.URL + "/api/v1/auth/ics/refresh_token")
	// if err != nil {
	// 	t.Fatalf("Expected no error, got %v", err)
	// }
	// if res.StatusCode != 200 {
	// 	t.Fatalf("Expected status code 200, got %v", res.StatusCode)
	// }

	// token = struct{
	// 	token string
	// 	token_expiry string
	// 	token_type string
	// }{
	// 	token: res.Body.token,
	// 	token_expiry: res.Body.token_expiry,

	// }
}

func TestGraphAPI(t *testing.T) {

}

func TestAuthAPI(t *testing.T) {

}

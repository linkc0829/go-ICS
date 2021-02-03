<img src="https://github.com/linkc0829/go-icsharing/blob/master/golang+grqphql.png?raw=true" />

# Go-ICSharing
This repo demo an income/cost planing and sharing fullstack web app, using a GraphQL server written in Go with GQLGen, Gin, Goth, GORM, MongoDB, Redis, JWT. The demo UI is written with JavaScript and AJAX.


## System Requirements
* [Golang](https://golang.org/)
* [GCC](https://gcc.gnu.org/) for SQLite
* [git](https://git-scm.com/)
* [Docker](https://www.docker.com/)
* A modern browser to run the app

For Windows users, install GCC provided by [MSYS2](https://www.msys2.org/)

You can also build environment with Docker(see followings)

## Basic Usage

### Clone

`$ git clone https://github.com/linkc0829/go-icsharing.git`

### Deployment With docker-compose

`$ docker-compose up -d`

1. will build a multistaged build using docker/prod.dockerfile to produce a minimal graphql server image
2. will create 3 containers: ICSharing server, mongoDB and redis


### Build Docker image using Dockerfile

`$ docker build -f ./docker/prod.dockerfile -t icsharing .`


### Deploy with Docker image pulling from Dockerhub

`$ docker run -d --name icsharing -p 8080:8080 --env-file .env linkc0829/go-icsharing`


### Deploy to GCP Compute Engine
https://icsharing.com/
* could have invalid certificate if renewal limits exceed, switch to staging environment and using [Fake LE Intermediate X1](https://letsencrypt.org/certs/fakeleintermediatex1.pem)


### Deploy to GCP Cloud Run
https://icsharing-amyzrpzc2q-de.a.run.app/


### Deploy to Heroku
https://icsharing.herokuapp.com/


### Demo Account
* Username: test567
* email: test567@gmail.com
* password: 123456


### CI/CD with drone
<img src="https://github.com/linkc0829/go-icsharing/blob/master/droneScreenshot.png?raw=true" /><br>


## 軟體介紹

這是一個使用Golang撰寫，與家人朋友分享個人收入支出規劃的社交web app<br>

你有在記帳嗎? 每次記帳只能得到一個結論: 入不敷出？<br>
那是因為記帳是被動的，她只能協助你看到過去，無法規畫未來<br>
但財富的累積，靠的是未來規劃！<br><br>

本程式除了幫助您方便規劃未來的收入支出，提前預知可能的財務狀況<br>
也可以將規劃內容分享給親朋好友！大家一起幫你鼓勵打氣！<br>
也能練習提早建立理財觀念！<br><br>


### 主要特色
1. 收入支出分類: 將收入分成四大類別，支出分為五大類別，方便檢視金錢流向
2. 數據分析: 自動將已發生的收入支出移動到歷史紀錄，可以透過回顧特定區間內的規劃紀錄，分析與現實的實際收入支出差異，越接近代表規劃能力越強，財富執行力也越強
3. 登入/註冊/登出: 除了可在app上註冊使用者帳戶，也可連結google/facebook等帳戶登入
4. 社交功能: 可以加別人為好友(friend)，表示願意分享收入支出給朋友觀看，同時您的名字也會出現在對方的追隨者(follower)清單
5. 權限管理: 可以將收入支出設定PUBLIC(供所有人觀看)，FRIEND(僅供朋友觀看)，PRIVATE(僅自己能看)
6. 附單元測試，效能測試


### 後端使用技術
1. [Gin-gonic](https://github.com/gin-gonic/gin): 使用Gin作為web框架
2. [gqlgen](https://github.com/99designs/gqlgen): 使用gqlgen套件架設GraphQL Server端
3. [jwt-go](https://github.com/dgrijalva/jwt-go): 使用JSON Web Token做使用者認證，達成client端的狀態分離提升擴展性、在middleware透過token檢查權限、發放access token(存在momory)與refresh token(存在cookies)，並實作soft refresh，server端會將refresh token存在redis加速存取速度
4. [Goth](https://github.com/markbates/goth): 使用Goth的套件提供OAuth2使用者認證
5. [MongoDB](https://go.mongodb.org/mongo-driver): 儲存使用者資料與收入支出資料，並利用Goroutine實作Job Queue進行資料同步最佳化，方便擴展部屬
6. [GORM](https://github.com/jinzhu/gorm): 使用GORM在memory建立sqlite，作為軟體試用介面的資料庫
7. [Redigo](https://github.com/gomodule/redigo): 使用套件作為Redis Client，存取Redis中的refresh token
8. [graphql](https://github.com/shurcooL/graphql): 架設GraphQL Client端，提供restful API
9. [autocert](https://golang.org/x/crypto/acme/autocert): 自動跟Let's Encrypt要求SSL憑證，建立HTTPS server


### Restful + GraphQL 的API
1. 使用者安全資訊透過Restful API傳送
2. GraphQL Client端提供Restful API: 供不會GraphQL的開發者使用
3. GraphQL API: 提供高自由度的資料存取API
4. API權限管理: admin可進行任意CRUD，一般User可以對自己擁有的資料CRUD，經過授權(加好友)則可進行部分查詢與修改


### 前端使用技術
1. AJAX: 使用superagent對接API，進行非同步的資料存取，並透過原生Javascript語法控制互動介面<br>
2. Golang Template: 使用Golang內建的template做版面元件控制<br>
3. JWT: 取得server傳送的access token，並存在memory，在要求API存取資料時附加在Request Header<br><br>


### DevOps使用技術
1. Docker: 利用容器技術方便擴展及部屬, docker-compose up 完成架構部屬
2. Drone: CI/CD工具，自動執行單元測試，並發布映像檔至Dockerhub、Google與Heroku Container Registry



package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yangtaishan/core"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))

var db *sql.DB

func makeHandler(fn func(http.ResponseWriter, *http.Request, *core.Api)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		api := &core.Api{}
		fn(rw, r, api)
		result, _ := json.Marshal(api)
		rw.Write(result)
	}
}

func getArticles(w http.ResponseWriter, r *http.Request, api *core.Api) {
	articles, _ := GetArticles(db)

	api.Code = 0
	api.Message = "success"
	api.Data = map[string][]core.BaseArticle{"list": articles}
}

func getArticle(w http.ResponseWriter, r *http.Request, api *core.Api) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get("id"))
	article, _ := GetArticleByID(id)
	api.Code = 0
	api.Message = "success"
	api.Data = article
}

func indexHandler(w http.ResponseWriter, r *http.Request, api *core.Api) {
	api.Code = 0
	api.Message = "success"
	api.Data = "this is index page"
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Info("Error loading .env file")
	}

	db, err = InitDB()

	if err != nil {
		fmt.Printf("db err: %v", err)
	}

	http.HandleFunc("/", makeHandler(indexHandler))
	http.HandleFunc("/api/news", makeHandler(getArticles))
	http.HandleFunc("/api/news/detail", makeHandler(getArticle))

	logger.Info("server starting: http://localhost:8091")

	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		logger.Fatal("server error", zap.Error(err))
	}

}

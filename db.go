package main

import (
	"database/sql"
	"fmt"
	"os"

	"yangtaishan/core"

	"github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	cfg := mysql.Config{
		Addr:                 os.Getenv("HOST") + ":" + os.Getenv("PORT"),
		User:                 os.Getenv("USER_NAME"),
		Passwd:               os.Getenv("PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		Net:                  "tcp",
		AllowNativePasswords: true, // mysql8 required
	}

	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		fmt.Print("err", err)
		return nil, fmt.Errorf("error %v", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, fmt.Errorf("error %v", err)
	}

	return db, nil
}

func GetArticles(db *sql.DB, limit int, offset int) ([]core.BaseArticle, error) {
	var articles []core.BaseArticle
	rows, err := db.Query("SELECT id,title,date,image_url FROM news ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)

	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var article core.BaseArticle
		if err := rows.Scan(&article.ID, &article.Title, &article.Date, &article.ImageUrl); err != nil {
			return nil, fmt.Errorf("err: %v", err)
		}
		// fmt.Println("res", article)
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("err %v", err)
	}

	return articles, nil
}

func GetArticleByID(id int) (core.Article, error) {
	var article core.Article
	row := db.QueryRow("SELECT id, title,image_url,date,source,audio_url,transcript FROM news WHERE id = ?", id)
	if err := row.Scan(&article.ID, &article.Title, &article.ImageUrl, &article.Date, &article.Source, &article.AudioUrl, &article.Transcript); err != nil {
		if err == sql.ErrNoRows {
			return article, fmt.Errorf("articleByID %d: no such article", id)
		}
		return article, fmt.Errorf("articleByID %d: %v", id, err)
	}
	return article, nil
}

func AddAticle(db *sql.DB, article core.Article) (int64, error) {
	res, err := db.Exec("INSERT INTO news (title) VALUES (?)", article.Title)

	if err != nil {
		return 0, fmt.Errorf("err: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("err: %v", err)
	}

	return id, nil
}

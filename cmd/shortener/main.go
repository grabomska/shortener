package main

import (
	"database/sql"
	"github.com/grabomska/shortener/internal/handler"
	"github.com/grabomska/shortener/internal/repository"
	"github.com/grabomska/shortener/internal/service"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./shortener.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqLiteRepository := repository.NewSQLiteRepository(db)
	urlShortener := service.NewShortenerService(sqLiteRepository)

	createHandler := handler.NewCreateShortURLHandler(urlShortener)
	getHandler := handler.NewGetURLHandler(urlShortener)

	mux := http.NewServeMux()
	mux.Handle("POST /", createHandler)
	mux.Handle("GET /{id}", getHandler)

	err = http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

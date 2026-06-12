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
	httpHandler := handler.NewHandler(urlShortener)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", httpHandler.CreateShort)
	mux.HandleFunc("GET /{id}", httpHandler.GetUrl)

	err = http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

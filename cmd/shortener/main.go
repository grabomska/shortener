package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/config"
	"github.com/grabomska/shortener/internal/handler"
	"github.com/grabomska/shortener/internal/repository"
	"github.com/grabomska/shortener/internal/service"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("sqlite", "./shortener.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqLiteRepository := repository.NewSQLiteRepository(db)
	urlShortener := service.NewShortenerService(sqLiteRepository)
	httpHandler := handler.NewHandler(cfg, urlShortener)

	r := gin.Default()

	r.POST("/", httpHandler.CreateShort)
	r.GET("/:id", httpHandler.GetUrl)

	r.Run(cfg.ServerAddress)
}

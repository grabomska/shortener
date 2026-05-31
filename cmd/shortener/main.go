package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/handler"
	"github.com/grabomska/shortener/internal/repository"
	"github.com/grabomska/shortener/internal/service"
	"log"
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

	r := gin.Default()

	r.POST("/", httpHandler.CreateShort)
	r.GET("/:id", httpHandler.GetUrl)

	r.Run(":8080")
}

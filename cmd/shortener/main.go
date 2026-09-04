package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/config"
	"github.com/grabomska/shortener/internal/handler"
	"github.com/grabomska/shortener/internal/middleware"
	"github.com/grabomska/shortener/internal/repository"
	"github.com/grabomska/shortener/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("sqlite", "./shortener.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//services
	sqLiteRepository := repository.NewSQLiteRepository(db)
	urlShortener := service.NewShortenerService(sqLiteRepository)

	//logger
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	//handler
	httpHandler := handler.NewHandler(cfg, urlShortener)

	r := gin.New()

	//middlewares
	r.Use(middleware.Logger(sugar))
	r.Use(gin.Recovery())

	// routes
	r.POST("/", httpHandler.CreateShort)
	r.GET("/:id", httpHandler.GetURL)
	r.POST("api/shorten", httpHandler.CreateShorten)

	r.Run(cfg.ServerAddress)
}

package handler

import (
	"github.com/grabomska/shortener/internal/config"
	"github.com/grabomska/shortener/internal/service"
)

type Handler struct {
	cfg     *config.Config
	service service.ShortenerServiceInterface
}

func NewHandler(cfg *config.Config, service service.ShortenerServiceInterface) *Handler {
	return &Handler{
		cfg:     cfg,
		service: service,
	}
}

package handler

import "github.com/grabomska/shortener/internal/service"

type Handler struct {
	service service.ShortenerServiceInterface
}

func NewHandler(service service.ShortenerServiceInterface) *Handler {
	return &Handler{service: service}
}

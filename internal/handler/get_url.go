package handler

import (
	"github.com/grabomska/shortener/internal/service"
	"net/http"
)

type GetURLHandler struct {
	service *service.ShortenerService
}

func NewGetURLHandler(service *service.ShortenerService) *GetURLHandler {
	return &GetURLHandler{service: service}
}

func (h *GetURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	short := r.PathValue("id")
	shortUrl, err := h.service.GetFullUrlByShort(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, shortUrl.Url, http.StatusTemporaryRedirect)
}

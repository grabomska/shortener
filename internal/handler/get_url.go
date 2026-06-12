package handler

import (
	"net/http"
)

func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("id")
	shortUrl, err := h.service.GetFullUrlByShort(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, shortUrl.Url, http.StatusTemporaryRedirect)
}

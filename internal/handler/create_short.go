package handler

import (
	"io"
	"net/http"
)

func (h *Handler) CreateShort(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 1024)
	body, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := string(buf[:body])
	shorted, err := h.service.CreateShortUrl(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	resp := []byte(scheme + "://" + r.Host + "/" + shorted.Short)

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

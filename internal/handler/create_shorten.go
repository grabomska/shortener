package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createShortenRequest struct {
	URL string `json:"url"`
}

type createShortenResponse struct {
	Result string `json:"result"`
}

func (h *Handler) CreateShorten(c *gin.Context) {
	if c.ContentType() != "application/json" {
		c.String(http.StatusUnsupportedMediaType, "Unsupported Media Type")
		return
	}

	var request createShortenRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	shorted, err := h.service.CreateShortURL(request.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "create short url failed")
		return
	}

	response := createShortenResponse{
		Result: h.cfg.BaseURL + "/" + shorted.Short,
	}

	c.JSON(http.StatusOK, response)
}

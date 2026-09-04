package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateShort(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	url := string(body)
	shorted, err := h.service.CreateShortURL(url)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, h.cfg.BaseURL+"/"+shorted.Short)
}

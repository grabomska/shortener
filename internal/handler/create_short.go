package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateShort(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	url := string(body)
	shorted, err := h.service.CreateShortUrl(url)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("content-type", "text/plain")

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	c.String(http.StatusCreated, scheme+"://"+c.Request.Host+"/"+shorted.Short)
}

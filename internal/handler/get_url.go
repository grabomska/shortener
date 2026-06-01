package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/service"
	"net/http"
)

func (h *Handler) GetUrl(c *gin.Context) {
	short := c.Param("id")
	shortUrl, err := h.service.GetFullUrlByShort(short)
	if err != nil {
		if errors.Is(err, service.ErrShortURLNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, shortUrl.Url)
}

package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/service"
	"net/http"
)

func (h *Handler) GetURL(c *gin.Context) {
	short := c.Param("id")
	shortURL, err := h.service.GetFullURLByShort(short)
	if err != nil {
		if errors.Is(err, service.ErrShortURLNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, shortURL.URL)
}

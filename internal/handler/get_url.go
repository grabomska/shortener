package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetUrl(c *gin.Context) {
	short := c.Param("id")
	shortUrl, err := h.service.GetFullUrlByShort(short)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, shortUrl.Url)
}

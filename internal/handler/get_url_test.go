package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/config"
	"github.com/grabomska/shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grabomska/shortener/internal/model"
	"github.com/grabomska/shortener/internal/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandlerGetUrlSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		GetFullUrlByShort("ABC123").
		Return(&model.ShortUrl{Short: "ABC123", Url: "https://example.com"}, nil)

	cfg := &config.Config{}
	h := NewHandler(cfg, mockService)

	req := httptest.NewRequest(http.MethodGet, "/ABC123", nil)
	req.SetPathValue("id", "ABC123")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "ABC123"}}

	h.GetUrl(c)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))
}

func TestHandlerGetUrlNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		GetFullUrlByShort("ABC123").
		Return(nil, service.ErrShortURLNotFound)

	cfg := &config.Config{}
	h := NewHandler(cfg, mockService)

	req := httptest.NewRequest(http.MethodGet, "/ABC123", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "ABC123"}}
	c.Request = req

	h.GetUrl(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, service.ErrShortURLNotFound.Error(), strings.TrimSpace(w.Body.String()))
}

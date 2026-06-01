package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grabomska/shortener/internal/model"
	"github.com/grabomska/shortener/internal/service/mocks"
	"go.uber.org/mock/gomock"
)

var (
	errReadBody = errors.New("read body error")
	errService  = errors.New("service error")
)

type errReader struct{}

func (r *errReader) Read(_ []byte) (int, error) {
	return 0, errReadBody
}

func TestHandlerGetUrlSuccessHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		CreateShortUrl("https://example.com").
		Return(&model.ShortUrl{Short: "ABC123", Url: "https://example.com"}, nil)

	cfg := &config.Config{}
	h := NewHandler(cfg, mockService)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateShort(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, cfg.ResultAddress+"/ABC123", w.Body.String())
}

func TestHandlerCreateShortReadBodyError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	cfg := &config.Config{}
	h := NewHandler(cfg, mockService)

	req := httptest.NewRequest(http.MethodPost, "/", &errReader{})
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateShort(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "read body error", strings.TrimSpace(w.Body.String()))
}

func TestHandlerCreateShortServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		CreateShortUrl("https://example.com").
		Return(nil, errService)

	cfg := &config.Config{}
	h := NewHandler(cfg, mockService)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateShort(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "service error", strings.TrimSpace(w.Body.String()))
}

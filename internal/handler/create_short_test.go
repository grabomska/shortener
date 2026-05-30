package handler

import (
	"crypto/tls"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grabomska/shortener/internal/model"
	"github.com/grabomska/shortener/internal/service/mocks"
	"go.uber.org/mock/gomock"
)

type errReader struct{}

func (r *errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read body error")
}

func TestHandlerGetUrlSuccessHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		CreateShortUrl("https://example.com").
		Return(&model.ShortUrl{Short: "ABC123", Url: "https://example.com"}, nil)

	h := NewHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	req.Host = "short.local"
	rr := httptest.NewRecorder()

	h.CreateShort(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
	assert.Equal(t, "http://short.local/ABC123", rr.Body.String())
}

func TestHandlerCreateShortSuccessHTTPS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		CreateShortUrl("https://example.com").
		Return(&model.ShortUrl{Short: "ABC123"}, nil)

	h := NewHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	req.Host = "short.local"
	req.TLS = &tls.ConnectionState{}
	rr := httptest.NewRecorder()

	h.CreateShort(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
	assert.Equal(t, "https://short.local/ABC123", rr.Body.String())
}

func TestHandlerCreateShortReadBodyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	h := NewHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/", &errReader{})
	rr := httptest.NewRecorder()

	h.CreateShort(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "read body error", strings.TrimSpace(rr.Body.String()))
}

func TestHandlerCreateShortServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		CreateShortUrl("https://example.com").
		Return(nil, errors.New("service error"))

	h := NewHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	rr := httptest.NewRecorder()

	h.CreateShort(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "service error", strings.TrimSpace(rr.Body.String()))
}

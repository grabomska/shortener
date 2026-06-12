package handler

import (
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

func TestHandlerGetUrlSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		GetFullUrlByShort("ABC123").
		Return(&model.ShortUrl{Short: "ABC123", Url: "https://example.com"}, nil)

	h := NewHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/ABC123", nil)
	req.SetPathValue("id", "ABC123")
	req.Host = "short.local"
	rr := httptest.NewRecorder()

	h.GetUrl(rr, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Equal(t, "https://example.com", rr.Header().Get("Location"))
}

func TestHandlerGetUrlNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockShortenerServiceInterface(ctrl)
	mockService.EXPECT().
		GetFullUrlByShort("ABC123").
		Return(nil, errors.New("not found"))

	h := NewHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/ABC123", nil)
	req.SetPathValue("id", "ABC123")
	req.Host = "short.local"
	rr := httptest.NewRecorder()

	h.GetUrl(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "not found", strings.TrimSpace(rr.Body.String()))
}

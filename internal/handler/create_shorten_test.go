package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/config"
	"github.com/grabomska/shortener/internal/model"
	"github.com/grabomska/shortener/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlerCreateShorten(t *testing.T) {
	gin.SetMode(gin.TestMode)

	serviceErr := errors.New("service error")
	tests := []struct {
		name             string
		contentType      string
		body             string
		setup            func(service *mocks.MockShortenerServiceInterface)
		wantStatus       int
		wantResponseType string
		wantBody         string
		wantJSON         bool
	}{
		{
			name:        "success",
			contentType: "application/json",
			body:        `{"url":"https://example.com"}`,
			setup: func(service *mocks.MockShortenerServiceInterface) {
				service.EXPECT().
					CreateShortURL("https://example.com").
					Return(&model.ShortURL{URL: "https://example.com", Short: "ABC123"}, nil)
			},
			wantStatus:       http.StatusOK,
			wantResponseType: "application/json; charset=utf-8",
			wantBody:         `{"result":"http://localhost:8080/ABC123"}`,
			wantJSON:         true,
		},
		{
			name:             "unsupported media type",
			contentType:      "text/plain",
			body:             "https://example.com",
			setup:            func(_ *mocks.MockShortenerServiceInterface) {},
			wantStatus:       http.StatusUnsupportedMediaType,
			wantResponseType: "text/plain; charset=utf-8",
			wantBody:         "Unsupported Media Type",
		},
		{
			name:             "invalid JSON",
			contentType:      "application/json",
			body:             `{"url":`,
			setup:            func(_ *mocks.MockShortenerServiceInterface) {},
			wantStatus:       http.StatusBadRequest,
			wantResponseType: "text/plain; charset=utf-8",
		},
		{
			name:        "service error",
			contentType: "application/json",
			body:        `{"url":"https://example.com"}`,
			setup: func(service *mocks.MockShortenerServiceInterface) {
				service.EXPECT().
					CreateShortURL("https://example.com").
					Return(nil, serviceErr)
			},
			wantStatus:       http.StatusInternalServerError,
			wantResponseType: "text/plain; charset=utf-8",
			wantBody:         "create short url failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockShortenerServiceInterface(ctrl)
			test.setup(mockService)

			h := NewHandler(&config.Config{BaseURL: "http://localhost:8080"}, mockService)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.body))
			req.Header.Set("Content-Type", test.contentType)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h.CreateShorten(c)

			assert.Equal(t, test.wantStatus, w.Code)
			assert.Equal(t, test.wantResponseType, w.Header().Get("Content-Type"))
			if test.wantJSON {
				assert.JSONEq(t, test.wantBody, w.Body.String())
			} else {
				assert.Equal(t, test.wantBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}

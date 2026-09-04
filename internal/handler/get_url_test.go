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

func TestHandlerGetURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setup        func(service *mocks.MockShortenerServiceInterface)
		wantStatus   int
		wantLocation string
		wantBody     string
	}{
		{
			name: "success",
			setup: func(service *mocks.MockShortenerServiceInterface) {
				service.EXPECT().
					GetFullURLByShort("ABC123").
					Return(&model.ShortURL{Short: "ABC123", URL: "https://example.com"}, nil)
			},
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://example.com",
		},
		{
			name: "not found",
			setup: func(mockService *mocks.MockShortenerServiceInterface) {
				mockService.EXPECT().
					GetFullURLByShort("ABC123").
					Return(nil, service.ErrShortURLNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   service.ErrShortURLNotFound.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockShortenerServiceInterface(ctrl)
			test.setup(mockService)

			h := NewHandler(&config.Config{}, mockService)
			req := httptest.NewRequest(http.MethodGet, "/ABC123", nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "ABC123"}}

			h.GetURL(c)

			assert.Equal(t, test.wantStatus, w.Code)
			assert.Equal(t, test.wantLocation, w.Header().Get("Location"))
			if test.wantBody != "" {
				assert.Equal(t, test.wantBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}

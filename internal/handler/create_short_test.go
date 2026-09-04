package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/grabomska/shortener/internal/config"
	"github.com/stretchr/testify/assert"

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

func TestHandlerCreateShort(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		body            func() io.Reader
		setup           func(service *mocks.MockShortenerServiceInterface)
		wantStatus      int
		wantContentType string
		wantBody        string
	}{
		{
			name: "success",
			body: func() io.Reader {
				return strings.NewReader("https://example.com")
			},
			setup: func(service *mocks.MockShortenerServiceInterface) {
				service.EXPECT().
					CreateShortURL("https://example.com").
					Return(&model.ShortURL{Short: "ABC123", URL: "https://example.com"}, nil)
			},
			wantStatus:      http.StatusCreated,
			wantContentType: "text/plain",
			wantBody:        "/ABC123",
		},
		{
			name: "read body error",
			body: func() io.Reader {
				return &errReader{}
			},
			setup:      func(_ *mocks.MockShortenerServiceInterface) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   errReadBody.Error(),
		},
		{
			name: "service error",
			body: func() io.Reader {
				return strings.NewReader("https://example.com")
			},
			setup: func(service *mocks.MockShortenerServiceInterface) {
				service.EXPECT().
					CreateShortURL("https://example.com").
					Return(nil, errService)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   errService.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockShortenerServiceInterface(ctrl)
			test.setup(mockService)

			cfg := &config.Config{}
			h := NewHandler(cfg, mockService)
			req := httptest.NewRequest(http.MethodPost, "/", test.body())
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h.CreateShort(c)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantContentType != "" {
				assert.Equal(t, test.wantContentType, w.Header().Get("Content-Type"))
			}
			assert.Equal(t, test.wantBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

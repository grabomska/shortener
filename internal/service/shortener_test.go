package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/grabomska/shortener/internal/model"
	"github.com/grabomska/shortener/internal/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestShortenerServiceCreateShortUrl(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(repo *mocks.MockRepository)
		url     string
		wantErr bool
	}{
		{
			name: "success",
			setup: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByShort(gomock.Any()).
					Return(nil, nil)

				repo.EXPECT().
					Create(gomock.Any()).
					DoAndReturn(func(shortURL *model.ShortUrl) error {
						assert.Equal(t, "https://example.com", shortURL.Url)
						assert.NotEqual(t, "", shortURL.Short)
						assert.Equal(t, defaultLength, len(shortURL.Short))

						return nil
					})
			},
			url: "https://example.com",
		},
		{
			name: "repository get by short error",
			setup: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByShort(gomock.Any()).
					Return(nil, errors.New("get by short error"))
			},
			url:     "https://example.com",
			wantErr: true,
		},
		{
			name: "repository create error",
			setup: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByShort(gomock.Any()).
					Return(nil, nil)

				repo.EXPECT().
					Create(gomock.Any()).
					Return(errors.New("create error"))
			},
			url:     "https://example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			tt.setup(repo)

			s := NewShortenerService(repo)

			got, err := s.CreateShortUrl(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.url, got.Url)
			assert.NotEqual(t, "", got.Short)
			assert.Equal(t, defaultLength, len(got.Short))
		})
	}
}

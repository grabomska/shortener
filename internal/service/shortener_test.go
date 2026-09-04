package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/grabomska/shortener/internal/model"
	"github.com/grabomska/shortener/internal/repository/mocks"
	"go.uber.org/mock/gomock"
)

var (
	errGetByShort = errors.New("get by short error")
	errCreate     = errors.New("create error")
)

func TestShortenerServiceCreateShortURL(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *mocks.MockRepository)
		url     string
		wantErr bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByShort(gomock.Any()).
					Return(nil, nil)

				repo.EXPECT().
					Create(gomock.Any()).
					DoAndReturn(func(shortURL *model.ShortURL) error {
						assert.Equal(t, "https://example.com", shortURL.URL)
						assert.NotEqual(t, "", shortURL.Short)
						assert.Equal(t, defaultLength, len(shortURL.Short))

						return nil
					})
			},
			url: "https://example.com",
		},
		{
			name: "repository get by short error",
			setup: func(_ *testing.T, repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByShort(gomock.Any()).
					Return(nil, errGetByShort)
			},
			url:     "https://example.com",
			wantErr: true,
		},
		{
			name: "repository create error",
			setup: func(_ *testing.T, repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByShort(gomock.Any()).
					Return(nil, nil)

				repo.EXPECT().
					Create(gomock.Any()).
					Return(errCreate)
			},
			url:     "https://example.com",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			test.setup(t, repo)

			s := NewShortenerService(repo)

			got, err := s.CreateShortURL(test.url)
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, test.url, got.URL)
			assert.NotEqual(t, "", got.Short)
			assert.Equal(t, defaultLength, len(got.Short))
		})
	}
}

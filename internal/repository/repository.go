package repository

import "github.com/grabomska/shortener/internal/model"

type Repository interface {
	Create(shortURL *model.ShortURL) error
	GetByShort(short string) (*model.ShortURL, error)
	GetByURL(url string) (*model.ShortURL, error)
}

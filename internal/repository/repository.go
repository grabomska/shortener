package repository

import "github.com/grabomska/shortener/internal/model"

type Repository interface {
	Create(shortUrl *model.ShortUrl) error
	GetByShort(short string) (*model.ShortUrl, error)
	GetByUrl(url string) (*model.ShortUrl, error)
}

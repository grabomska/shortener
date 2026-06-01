package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/grabomska/shortener/internal/model"
	"github.com/grabomska/shortener/internal/repository"
	"math/big"
)

const (
	alphabet      = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultLength = 6
	maxAttempts   = 3
)

var ErrShortURLNotFound = errors.New("short url not found")

type ShortenerServiceInterface interface {
	CreateShortUrl(url string) (*model.ShortUrl, error)
	GetFullUrlByShort(short string) (*model.ShortUrl, error)
}

type ShortenerService struct {
	repo repository.Repository
}

func NewShortenerService(repo repository.Repository) ShortenerServiceInterface {
	return &ShortenerService{repo}
}

func (s *ShortenerService) CreateShortUrl(url string) (*model.ShortUrl, error) {
	shorted, err := s.generateUniqueShortURL(defaultLength)
	if err != nil {
		return nil, err
	}

	shortUrl := &model.ShortUrl{
		Url:   url,
		Short: shorted,
	}

	err = s.repo.Create(shortUrl)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return shortUrl, nil
}

func (s *ShortenerService) GetFullUrlByShort(short string) (*model.ShortUrl, error) {
	shortURL, err := s.repo.GetByShort(short)
	if err != nil {
		return nil, fmt.Errorf("get full url: %w", err)
	}

	if shortURL == nil {
		return nil, ErrShortURLNotFound
	}

	return shortURL, nil
}

// generateUniqueShortURL генерирует уникальный короткий URL
func (s *ShortenerService) generateUniqueShortURL(length int) (string, error) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Генерируем случайную строку
		short, err := s.generateRandomString(length)
		if err != nil {
			return "", err
		}

		// Проверяем уникальность
		isUnique, err := s.isShortUnique(short)
		if err != nil {
			return "", err
		}

		if isUnique {
			return short, nil
		}
	}

	// Если не удалось сгенерировать уникальный, увеличиваем длину
	return s.generateUniqueShortURL(length + 1)
}

// generateRandomString генерирует случайную строку заданной длины
func (s *ShortenerService) generateRandomString(length int) (string, error) {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		result[i] = alphabet[num.Int64()]
	}

	return string(result), nil
}

func (s *ShortenerService) isShortUnique(short string) (bool, error) {
	shortURL, err := s.repo.GetByShort(short)
	if err != nil {
		return false, err
	}

	if shortURL != nil {
		return false, nil
	}

	return true, nil
}

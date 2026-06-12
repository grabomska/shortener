package repository

import (
	"database/sql"
	"errors"

	"github.com/grabomska/shortener/internal/model"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) Repository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Create(url *model.ShortUrl) error {
	insert := `INSERT INTO short_url (url, short) VALUES (?, ?)`
	_, err := r.db.Exec(insert, url.Url, url.Short)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteRepository) GetByShort(short string) (*model.ShortUrl, error) {
	var shortUrl model.ShortUrl

	query := `SELECT url, short FROM short_url WHERE short = ?`
	err := r.db.QueryRow(query, short).Scan(&shortUrl.Url, &shortUrl.Short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &shortUrl, nil
}

func (r *SQLiteRepository) GetByUrl(url string) (*model.ShortUrl, error) {
	var shortUrl model.ShortUrl

	query := `SELECT url, short FROM short_url WHERE url = ?`
	err := r.db.QueryRow(query, url).Scan(&shortUrl.Url, &shortUrl.Short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &shortUrl, nil
}

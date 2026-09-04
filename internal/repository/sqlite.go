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

func (r *SQLiteRepository) Create(url *model.ShortURL) error {
	insert := `INSERT INTO short_url (url, short) VALUES (?, ?)`
	_, err := r.db.Exec(insert, url.URL, url.Short)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteRepository) GetByShort(short string) (*model.ShortURL, error) {
	var shortURL model.ShortURL

	query := `SELECT url, short FROM short_url WHERE short = ?`
	err := r.db.QueryRow(query, short).Scan(&shortURL.URL, &shortURL.Short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &shortURL, nil
}

func (r *SQLiteRepository) GetByURL(url string) (*model.ShortURL, error) {
	var shortURL model.ShortURL

	query := `SELECT url, short FROM short_url WHERE url = ?`
	err := r.db.QueryRow(query, url).Scan(&shortURL.URL, &shortURL.Short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &shortURL, nil
}

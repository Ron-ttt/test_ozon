package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(dbname string) (Storage, error) {

	db, err := sql.Open("postgres", dbname)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		return nil, err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &DBStorage{db: db}, nil
}

func (s *DBStorage) Add(key string, value string) error {
	_, err := s.db.Exec("INSERT INTO links (shorturl, originalurl) VALUES($1, $2)", key, value)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *DBStorage) Get(key string) (string, error) {
	rows := s.db.QueryRow("SELECT originalurl FROM links WHERE shorturl= $1", key)

	var originalURL string
	err := rows.Scan(&originalURL)

	if err != nil {
		return "", err
	}
	return originalURL, nil
}
func (s *DBStorage) Find(originalUrl string) (string, error) {
	rows := s.db.QueryRow("SELECT shorturl FROM links WHERE originalurl= $1", originalUrl)
	var short string
	err := rows.Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

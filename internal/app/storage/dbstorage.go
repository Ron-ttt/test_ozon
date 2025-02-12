package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(dbname string) (Storage, error) {

	db, err := sql.Open("postgres", dbname)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println(err)
		return nil, err
	}
	log.Println("DB connected")
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

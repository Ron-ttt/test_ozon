package storage

import "errors"

type Storage interface {
	Add(key string, value string) error
	Get(key string) (string, error)
	Find(originalurl string) (string, error)
}

type MapStorage struct {
	m map[string]string
}

func NewMapStorage() Storage {
	return &MapStorage{
		m: make(map[string]string),
	}
}

func (s *MapStorage) Add(key string, value string) error {
	s.m[key] = value
	return nil
}

func (s *MapStorage) Get(key string) (string, error) {
	value, found := s.m[key]
	if !found {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (s *MapStorage) Find(originalurl string) (string, error) {
	return "", errors.New("")
}

package storage

import "errors"

type MStorage struct{}

func NewMockStorage() Storage {
	return &MStorage{}
}

func (s *MStorage) Add(key string, value string) error {
	return nil
}

func (s *MStorage) Get(key string) (string, error) {
	if key == "invalid" {
		return "", errors.New("key not found")
	}
	return "http://love_nika", nil

}
func (s *MStorage) Find(originalurl string) (string, error) {
	return "", errors.New("")
}

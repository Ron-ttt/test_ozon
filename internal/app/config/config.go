package config

import (
	"flag"
	"os"
)

func Flags() (string, string, string) {
	// Определение флагов
	address := flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL") // порты должны совпадать
	db := flag.String("d", "", "адрес для бд")

	// Парсинг флагов
	flag.Parse()
	if envAddress := os.Getenv("SERVER_ADDRESS"); envAddress != "" {
		*address = envAddress
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		*baseURL = envBaseURL
	}

	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		*db = envDB
	}

	return *address, *baseURL, *db
}

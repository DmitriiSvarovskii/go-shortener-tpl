package config

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServiceURL       string
	BaseShortenerURL string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Флаги для конфигурации
	flag.StringVar(&cfg.ServiceURL, "a", "localhost:8080", "Адрес запуска HTTP-сервера (в формате host:port)")
	flag.StringVar(&cfg.BaseShortenerURL, "b", "http://localhost:8080", "Базовый адрес сокращённых URL")

	flag.Parse()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if envServiceURL := os.Getenv("SERVER_ADDRESS"); envServiceURL != "" {
		cfg.ServiceURL = envServiceURL
	}
	if envBaseShortenerURL := os.Getenv("BASE_URL"); envBaseShortenerURL != "" {
		cfg.BaseShortenerURL = envBaseShortenerURL
	}

	return cfg, nil
}

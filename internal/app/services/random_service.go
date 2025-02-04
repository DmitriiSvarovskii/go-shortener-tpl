package services

import (
	"errors"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/storage"
)

var ErrKeyNotFound = errors.New("key not found")

type RandomService struct {
	storage storage.Repository
}

func NewRandomService(storage storage.Repository) *RandomService {
	return &RandomService{
		storage: storage,
	}
}

func (r *RandomService) GenerateShortURL(value string) string {
	key := randStr()
	r.storage.Set(key, value)
	return key
}

func (r *RandomService) GetOriginURL(key string) (string, error) {
	val, exists := r.storage.Get(key)
	if !exists {
		return "", ErrKeyNotFound
	}
	return val, nil
}

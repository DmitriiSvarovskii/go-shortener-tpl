package storage

type Repository interface {
	Get(key string) (string, bool)
	Set(key, url string)
}
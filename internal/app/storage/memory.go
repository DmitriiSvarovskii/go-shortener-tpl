package storage

type MemoryStorage struct {
	urls map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]string),
	}
}

func (m *MemoryStorage) Get(key string) (string, bool) {
	url, ok := m.urls[key]
	return url, ok
}

func (m *MemoryStorage) Set(key, url string) {
	m.urls[key] = url
}

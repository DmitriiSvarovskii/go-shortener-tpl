package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/config"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	data map[string]string
}

func NewMockStorage() *MockStorage {
	return &MockStorage{data: make(map[string]string)}
}

func (m *MockStorage) Get(key string) (string, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *MockStorage) Set(key, url string) {
	fmt.Printf("Saving: %s -> %s\n", key, url)
	m.data[key] = url
}

// Запуск реального HTTP-сервера на 8888
func startRealServer() *http.Server {
	repo := NewMockStorage()
	service := services.NewRandomService(repo)
	cfg := &config.AppConfig{
		ServiceURL:       "http://localhost:8888",
		BaseShortenerURL: "http://localhost:8888",
	}
	handler := NewHandler(service, cfg)

	r := chi.NewRouter()
	r.Post("/", handler.CreateShortURLHandler)
	r.Get("/{shortURL}", handler.GetOriginalURLHandler)
	r.MethodNotAllowed(handler.MethodNotAllowedHandle)

	srv := &http.Server{Addr: "localhost:8888", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Даем серверу немного времени на запуск
	time.Sleep(500 * time.Millisecond)

	return srv
}

func TestHandlers(t *testing.T) {
	// Запускаем реальный сервер
	server := startRealServer()
	defer server.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Отключаем авто-редирект
		},
	}
	
	testURL := "https://example.com"
	resp, err := http.Post("http://localhost:8888/", "text/plain", strings.NewReader(testURL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	var shortPath string
	_, err = fmt.Fscanf(resp.Body, "%s", &shortPath)
	assert.NoError(t, err)

	t.Run("GET existing short URL", func(t *testing.T) {
		resp, err := client.Get(shortPath)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		defer resp.Body.Close()
	})

	t.Run("GET non-existing short URL", func(t *testing.T) {
		resp, err := client.Get("http://localhost:8888/notExist")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		defer resp.Body.Close()
	})

	t.Run("Invalid method PUT", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8888/", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		defer resp.Body.Close()
	})
}

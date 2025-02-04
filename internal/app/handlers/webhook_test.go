package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"

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
    m.data[key] = url
}

func TestWebhook(t *testing.T) {
	mockRepo := NewMockStorage() // Используем реальное in-memory хранилище
	service := services.NewRandomService(mockRepo)
	handler := NewHandler(service)

	// Создаем тестовый URL заранее, чтобы использовать в GET-запросах
	testURL := "https://example.com"
	shortKey := service.GenerateShortURL(testURL)

	testCases := []struct {
		name         string
		method       string
		body         string
		url          string
		expectedCode int
		expectBody   bool
	}{
		{"POST valid URL", http.MethodPost, testURL, "/", http.StatusCreated, true},
		{"GET existing short URL", http.MethodGet, "", "/" + shortKey, http.StatusTemporaryRedirect, false},
		{"GET non-existing short URL", http.MethodGet, "", "/notExist", http.StatusBadRequest, false},
		{"Invalid method PUT", http.MethodPut, "", "/", http.StatusBadRequest, false},
		{"Invalid method DELETE", http.MethodDelete, "", "/", http.StatusBadRequest, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			handler.Webhook(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")

			if tc.expectBody {
				body, _ := io.ReadAll(resp.Body)
				assert.Contains(t, string(body), "http://localhost:8080/", "Ответ не содержит короткий URL")
			}
		})
	}
}
package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	m.data[key] = url
}

func setupTestServer() *httptest.Server {
	repo := NewMockStorage()
	service := services.NewRandomService(repo)
	handler := NewHandler(service)

	r := chi.NewRouter()
	r.Post("/", handler.CreateShortURLHandler)
	r.Get("/{shortURL}", handler.GetOriginalURLHandler)
	r.MethodNotAllowed(handler.MethodNotAllowedHandle)

	return httptest.NewServer(r)
}

func TestHandlers(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Отключаем авто-редирект
		},
	}

	testURL := "https://example.com"
	resp, err := http.Post(ts.URL+"/", "text/plain", strings.NewReader(testURL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	var shortURL string
	_, err = fmt.Fscanf(resp.Body, "http://localhost:8080/%s", &shortURL)
	assert.NoError(t, err)

	t.Run("GET existing short URL", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/" + shortURL)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
		defer resp.Body.Close()
	})

	t.Run("GET non-existing short URL", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/notExist")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		defer resp.Body.Close()
	})

	t.Run("Invalid method PUT", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		defer resp.Body.Close()
	})
}

// package handlers

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
// 	"github.com/go-resty/resty/v2"

// 	"github.com/stretchr/testify/assert"
// )

// type MockStorage struct {
// 	data map[string]string
// }

// func NewMockStorage() *MockStorage {
// 	return &MockStorage{data: make(map[string]string)}
// }

// func (m *MockStorage) Get(key string) (string, bool) {
// 	val, ok := m.data[key]
// 	return val, ok
// }

// func (m *MockStorage) Set(key, url string) {
// 	m.data[key] = url
// }

// func TestWebhook(t *testing.T) {
// 	mockRepo := NewMockStorage() // Используем реальное in-memory хранилище
// 	service := services.NewRandomService(mockRepo)
// 	handler := NewHandler(service)

// 	srv := httptest.NewServer(handler)
// 	// останавливаем сервер после завершения теста
// 	defer srv.Close()

// 	// Создаем тестовый URL заранее, чтобы использовать в GET-запросах
// 	testURL := "https://example.com"
// 	shortKey := service.GenerateShortURL(testURL)

// 	testCases := []struct {
// 		name         string
// 		method       string
// 		body         string
// 		url          string
// 		expectedCode int
// 		expectBody   bool
// 	}{
// 		{"POST valid URL", http.MethodPost, testURL, "/", http.StatusCreated, true},
// 		{"GET existing short URL", http.MethodGet, "", "/" + shortKey, http.StatusTemporaryRedirect, false},
// 		{"GET non-existing short URL", http.MethodGet, "", "/notExist", http.StatusBadRequest, false},
// 		{"Invalid method PUT", http.MethodPut, "", "/", http.StatusBadRequest, false},
// 		{"Invalid method DELETE", http.MethodDelete, "", "/", http.StatusBadRequest, false},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
// 			// w := httptest.NewRecorder()

// 			// handler.Webhook(w, req)

// 			// resp := w.Result()
// 			// defer resp.Body.Close()
// 			req := resty.New().R()
// 			req.SetRedirectPolicy(resty.NoRedirect) // Отключаем следование за редиректами
// 			req.Method = tc.method
// 			req.URL = srv.URL + tc.url
// 			resp, err := req.Send()
// 			assert.NoError(t, err, "error making HTTP request")

// 			t.Logf("tc.expectedCode: %v", tc.expectedCode)
// 			t.Logf("Status Code: %d", resp.StatusCode())
// 			// t.Logf("resp.StatusCode: %v", resp.StatusCode)
// 			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Код ответа не совпадает с ожидаемым")

// 			if tc.expectBody {
// 				body := resp.Body()
// 				assert.Contains(t, string(body), "http://localhost:8080/", "Ответ не содержит короткий URL")
// 			}
// 		})
// 	}
// }

package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *services.RandomService
}

func NewHandler(service *services.RandomService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	key := h.service.GenerateShortURL(string(body))
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("http://localhost:8080/" + key))
}

func (h *Handler) GetOriginalURLHandler(rw http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "shortURL")

	if key == "" {
		http.Error(rw, "key param is missed", http.StatusBadRequest)
		return
	}

	value, err := h.service.GetOriginURL(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if value == "" {
		http.Error(rw, "original URL is empty", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Location", value)
	rw.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Println("Response Status Code:", http.StatusTemporaryRedirect)
}

func (h *Handler) MethodNotAllowedHandle(rw http.ResponseWriter, r *http.Request) {
	responseMessage := fmt.Sprintf("The method '%s' is not allowed for path '%s'.", r.Method, r.URL.Path)
	rw.WriteHeader(http.StatusMethodNotAllowed)
	io.WriteString(rw, responseMessage)
}

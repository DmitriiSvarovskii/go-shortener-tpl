package handlers

import (
	"io"
	"net/http"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
)

type Handler struct {
	service *services.RandomService
}

func NewHandler(service *services.RandomService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		key := h.service.GenerateShortURL(string(body))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + key))
	case http.MethodGet:
		URLId := r.URL.Path[len("/"):]
		value, err := h.service.GetOriginURL(URLId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", value)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
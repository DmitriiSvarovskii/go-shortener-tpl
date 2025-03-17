package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/config"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/logger"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/models"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *services.RandomService
	cfg     *config.AppConfig
}

func NewHandler(service *services.RandomService, cfg *config.AppConfig) *Handler {
	return &Handler{service: service, cfg: cfg}
}

func (h *Handler) CreateShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	key := h.service.GenerateShortURL(string(body))

	fullURL := fmt.Sprintf("%s/%s", h.cfg.BaseShortenerURL, key)

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fullURL))
}

func (h *Handler) CreateJSONShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	key := h.service.GenerateShortURL(req.Url)
	fullURL := fmt.Sprintf("%s/%s", h.cfg.BaseShortenerURL, key)

	resp := models.Response{
		Result: fullURL,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 201 response")
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
}

func (h *Handler) MethodNotAllowedHandle(rw http.ResponseWriter, r *http.Request) {
	responseMessage := fmt.Sprintf("The method '%s' is not allowed for path '%s'.", r.Method, r.URL.Path)
	rw.WriteHeader(http.StatusMethodNotAllowed)
	io.WriteString(rw, responseMessage)
}

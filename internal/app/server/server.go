package server

import (
	"fmt"
	"net/http"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/config"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/handlers"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/logger"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func ShortenerRouter(cfg *config.AppConfig) *Server {
	repo := storage.NewMemoryStorage()
	service := services.NewRandomService(repo)
	handler := handlers.NewHandler(service, cfg)

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return logger.RequestLogger(next.ServeHTTP)
	})

	fmt.Println("Setting up route for shortURL")

	r.Post("/", handler.CreateShortURLHandler)
	r.Get("/{shortURL}", handler.GetOriginalURLHandler)
	r.MethodNotAllowed(handler.MethodNotAllowedHandle)

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.ServiceURL,
			Handler: r,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

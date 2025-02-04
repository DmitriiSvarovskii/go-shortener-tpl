package server

import (
	"fmt"
	"net/http"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/handlers"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func ShortenerRouter() *Server {
	repo := storage.NewMemoryStorage()
	service := services.NewRandomService(repo)
	handler := handlers.NewHandler(service)

	r := chi.NewRouter()

	fmt.Println("Setting up route for shortURL")

	r.Post("/", handler.CreateShortURLHandler)
	r.Get("/{shortURL}", handler.GetOriginalURLHandler)
	r.MethodNotAllowed(handler.MethodNotAllowedHandle)

	return &Server{
		httpServer: &http.Server{
			Handler: r,
		},
	}
}

func (s *Server) Run(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

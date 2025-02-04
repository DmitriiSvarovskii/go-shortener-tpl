package server

import (
	"net/http"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/handlers"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/services"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/storage"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	repo := storage.NewMemoryStorage()
	service := services.NewRandomService(repo)
	handler := handlers.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Webhook)

	return &Server{
		httpServer: &http.Server{
			Handler: mux,
		},
	}
}

func (s *Server) Run(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

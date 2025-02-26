package main

import (
	"log"

	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/config"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	srv := server.ShortenerRouter(cfg)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

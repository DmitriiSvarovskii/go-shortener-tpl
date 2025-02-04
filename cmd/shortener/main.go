package main

import (
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/config"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	srv := server.ShortenerRouter(cfg)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}

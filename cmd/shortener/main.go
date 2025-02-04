package main

import (
	"log"
	"github.com/DmitriiSvarovskii/go-shortener-tpl.git/internal/app/server"
)

func main() {
	srv := server.ShortenerRouter()
	if err := srv.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
// пакеты исполняемых приложений должны называться main
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(webhook))
}

// функция webhook — обработчик HTTP-запроса
func webhook(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		key, valuer := randStr(), string(body)

		memoryURL[key] = valuer
		w.Write([]byte(`http://localhost:8080/` + key))
	case r.Method == http.MethodGet:
		URLId := r.URL.Path[len("/"):]

		value, exists := memoryURL[URLId]

		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("memoryURL[URLId]:", value)
		w.Header().Set("Location", value)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

var memoryURL = make(map[string]string)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStr() string {
	length := 8

	buf := make([]byte, length)

	for i := 0; i < length; i++ {
		buf[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	str := string(buf)
	return str
}

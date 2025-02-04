package services

import "math/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
const length = 8 // Максимальная длина сгенерированного ключа

func randStr() string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(buf)
}

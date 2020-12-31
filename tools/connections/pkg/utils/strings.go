package utils

import (
	"math/rand"
	"time"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandonBytes return bytes by length
func RandonBytes(length int) []byte {
	if length <= 0 {
		return []byte("")
	}

	rand.Seed(time.Now().UnixNano())
	byts := make([]byte, length)
	for i := 0; i < length; i++ {
		byts[i] = letters[rand.Intn(len(letters))]
	}

	return byts
}

// RandomString return string by length, return "" when length < 0
func RandomString(length int) string {
	return string(RandonBytes(length))
}

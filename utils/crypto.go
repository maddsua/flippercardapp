package utils

import (
	"crypto/rand"
)

func NewRandomToken(n int) []byte {

	buff := make([]byte, max(n, 8))
	if _, err := rand.Read(buff); err != nil {
		panic(err)
	}

	return buff
}

const passwordDict = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"

const BcrypMaxPasswordLength = 72

func NewRandomBcryptPassword(n int) string {

	buff := NewRandomToken(min(n, BcrypMaxPasswordLength))
	pass := make([]byte, len(buff))

	for idx, val := range buff {
		pass[idx] = passwordDict[int(val)%len(passwordDict)]
	}

	return string(pass)
}

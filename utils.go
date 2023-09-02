package main

import "crypto/rand"

var chars = "abcdefghijklmnopqrstuvwxyz1234567890"

func randomString(length int) string {
	charsLength := len(chars)
	b := make([]byte, length)
	rand.Read(b) // generates len(b) random bytes
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%charsLength]
	}
	return string(b)
}

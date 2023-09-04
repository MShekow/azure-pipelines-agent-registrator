package main

import (
	"crypto/rand"
	"strings"
)

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

func GetCapabilitiesMapFromString(capabilities string) *map[string]string {
	m := map[string]string{}

	capabilitiesList := strings.Split(capabilities, ";")
	for _, capabilityStr := range capabilitiesList {
		capabilityTuple := strings.Split(capabilityStr, "=")
		if len(capabilityTuple) == 2 {
			m[capabilityTuple[0]] = capabilityTuple[1]
		}
	}

	return &m
}

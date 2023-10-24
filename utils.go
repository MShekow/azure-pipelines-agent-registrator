package main

import (
	"crypto/rand"
	"fmt"
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
		// Note: strings.Cut() slices around the first(!) instance of sep, returning the text before and after sep
		// We need this because an arg might look like this:
		// ExtraAgentContainers=name=c,image=some-image:latest,cpu=500m,memory=2Gi
		key, value, separatorWasFound := strings.Cut(capabilityStr, "=")
		if !separatorWasFound {
			fmt.Printf("WARNING: ignoring capability %s because the '=' separator and value are not set correctly\n", key)
		} else {
			m[key] = value
		}
	}

	return &m
}

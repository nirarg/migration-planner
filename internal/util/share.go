package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateShareToken generates a secure random token for sharing
func GenerateShareToken() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	// Encode to base64 URL-safe format
	token := base64.RawURLEncoding.EncodeToString(bytes)
	return token, nil
}

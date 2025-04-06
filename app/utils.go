package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func generateShortCode(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	h.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))

	hash := h.Sum(nil)[:8]
	shortCode := base64.URLEncoding.EncodeToString(hash)
	return shortCode[:8]
}

func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

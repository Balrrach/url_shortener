package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "URL Shortener"})
}

func createShortURL(c *gin.Context) {
	var request URLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	shortCode := generateShortCode(request.LongURL)

	var exists bool
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM short_urls WHERE short_code = $1)", shortCode).Scan(&exists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		shortCode = shortCode[:6] + fmt.Sprintf("%d", time.Now().UnixNano()%1000)
	}

	var id int64
	err := db.QueryRow(
		"INSERT INTO short_urls (short_code, long_url) VALUES ($1, $2) RETURNING id",
		shortCode, request.LongURL,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	shortURL := fmt.Sprintf("%s/%s", getBaseURL(c.Request), shortCode)

	c.JSON(http.StatusOK, gin.H{
		"id":         id,
		"short_code": shortCode,
		"short_url":  shortURL,
		"long_url":   request.LongURL,
	})
}

func redirectToOriginal(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var longURL string
	err := db.QueryRow("SELECT long_url FROM short_urls WHERE short_code = $1", shortCode).Scan(&longURL)

	if err == sql.ErrNoRows {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Short URL not found"})
		return
	} else if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Database error"})
		return
	}

	_, err = db.Exec("UPDATE short_urls SET clicks = clicks + 1 WHERE short_code = $1", shortCode)
	if err != nil {
		c.Error(err) // Optional: log error for click increment
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}

func getURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var url ShortURL
	err := db.QueryRow(`
		SELECT id, short_code, long_url, created_at, clicks
		FROM short_urls
		WHERE short_code = $1
	`, shortCode).Scan(&url.ID, &url.ShortCode, &url.LongURL, &url.CreatedAt, &url.Clicks)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, url)
}

package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	initDB()
	defer db.Close()

	// Setup router
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Routes
	router.GET("/", homePage)
	router.POST("/api/shorten", createShortURL)
	router.GET("/:shortCode", redirectToOriginal)
	router.GET("/api/stats/:shortCode", getURLStats)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	router.Run(":" + port)
}

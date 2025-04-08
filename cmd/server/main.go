package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Mark-Life/smart-webp-resize/internal/api"
	"github.com/Mark-Life/smart-webp-resize/internal/handler"
	"github.com/Mark-Life/smart-webp-resize/internal/processor"
)

func main() {
	// Create dependencies
	imageHandler := handler.NewImageHandler()
	imageProcessor := processor.New()
	
	// Create API
	imageAPI := api.NewImageAPI(imageHandler, imageProcessor)
	
	// Set up HTTP routes
	http.HandleFunc("/health", imageAPI.Health)
	http.HandleFunc("/process", imageAPI.ProcessFromURL)
	http.HandleFunc("/upload", imageAPI.ProcessFromUpload)
	
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Start server
	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 
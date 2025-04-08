package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Mark-Life/smart-webp-resize/internal/api"
	"github.com/Mark-Life/smart-webp-resize/internal/handler"
	"github.com/Mark-Life/smart-webp-resize/internal/processor"
	"github.com/Mark-Life/smart-webp-resize/pkg/config"
)

func main() {
	log.Println("Starting Smart WebP Resizer service...")
	
	// Load configuration
	cfg := config.New()
	
	// Create dependencies
	imageHandler := handler.NewImageHandler()
	imageProcessor := processor.New()
	
	// Create API
	imageAPI := api.NewImageAPI(imageHandler, imageProcessor)
	
	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", imageAPI.Health)
	mux.HandleFunc("/process/url", imageAPI.ProcessFromURL)
	mux.HandleFunc("/process/upload", imageAPI.ProcessFromUpload)
	
	// Set up static file server for test pages
	testDir := getTestDataDir()
	log.Printf("Serving test files from %s", testDir)
	testFileServer := http.FileServer(http.Dir(testDir))
	mux.Handle("/test/", http.StripPrefix("/test/", testFileServer))
	
	// Add a redirect for root to the test page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/test/test_image.html", http.StatusFound)
		} else {
			http.NotFound(w, r)
		}
	})
	
	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 90 * time.Second, // Longer timeout for image processing
		IdleTimeout:  120 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s...\n", cfg.Port)
		log.Printf("Test page available at http://localhost:%s/test/test_image.html\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server gracefully stopped")
}

// getTestDataDir returns the path to the test data directory
func getTestDataDir() string {
	// Get the executable directory
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Could not determine executable path: %v, using current directory", err)
		return filepath.Join(".", "test", "testdata")
	}
	
	exeDir := filepath.Dir(exePath)
	
	// Look for test data in different possible locations
	possiblePaths := []string{
		filepath.Join(exeDir, "test", "testdata"),                    // Same directory as executable
		filepath.Join(exeDir, "..", "test", "testdata"),              // One level up
		filepath.Join(exeDir, "..", "..", "test", "testdata"),        // Two levels up (for dev environment)
		filepath.Join(exeDir, "..", "..", "..", "test", "testdata"),  // Three levels up
	}
	
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// Fallback to relative path from current working directory
	log.Println("Warning: Could not find test data directory, falling back to current directory")
	return filepath.Join(".", "test", "testdata")
} 
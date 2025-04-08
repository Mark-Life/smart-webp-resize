package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port          string
	MaxWidth      int
	MaxHeight     int
	DefaultQuality int
}

// New creates a new Config with values from environment or defaults
func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	return &Config{
		Port:          port,
		MaxWidth:      1920, // Default max width for resizing
		MaxHeight:     1080, // Default max height for resizing
		DefaultQuality: 80,  // Default WebP quality (0-100)
	}
} 
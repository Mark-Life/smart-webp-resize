package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mark-Life/smart-webp-resize/internal/api"
	"github.com/Mark-Life/smart-webp-resize/pkg/config"
)

func main() {
	cfg := config.New()
	
	server := api.NewServer(cfg)
	
	fmt.Printf("Starting Smart WebP Resizer server on :%s...\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, server); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 
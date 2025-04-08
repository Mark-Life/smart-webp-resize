package main

import (
	"fmt"
	"os"
	// Import the server package when we're ready to use it
	// _ "github.com/Mark-Life/smart-webp-resize/cmd/server"
)

func main() {
	fmt.Println("Smart WebP Resizer")
	fmt.Println("For development, run: go run cmd/server/main.go")
	os.Exit(0)
}

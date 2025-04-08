package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mark-Life/smart-webp-resize/internal/handler"
	"github.com/Mark-Life/smart-webp-resize/internal/processor"
)

// ImageAPI handles HTTP requests for image processing
type ImageAPI struct {
	imageHandler handler.ImageHandler
	processor    processor.ImageProcessor
}

// NewImageAPI creates a new ImageAPI with the provided dependencies
func NewImageAPI(imageHandler handler.ImageHandler, processor processor.ImageProcessor) *ImageAPI {
	return &ImageAPI{
		imageHandler: imageHandler,
		processor:    processor,
	}
}

// ProcessFromURL handles image processing requests where the image is specified by URL
func (api *ImageAPI) ProcessFromURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get URL parameter
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Get processing options from query parameters
	options := getProcessOptionsFromRequest(r)

	// Fetch the image
	imageData, err := api.imageHandler.GetImageFromURL(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch image: %v", err), http.StatusBadRequest)
		return
	}

	// Process the image
	processedData, metadata, err := api.processor.ProcessFromBytes(imageData, &options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process image: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if metadata response is requested
	if r.URL.Query().Get("metadata") == "true" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(metadata); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode metadata: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Set filename for download if requested
	if r.URL.Query().Get("download") == "true" || r.URL.Query().Get("format") == "webp" {
		// Extract original filename from URL or use a generic name
		filename := "image.webp"
		if lastSlashIndex := strings.LastIndex(url, "/"); lastSlashIndex != -1 && lastSlashIndex < len(url)-1 {
			origFilename := url[lastSlashIndex+1:]
			// Remove query parameters if any
			if queryIndex := strings.Index(origFilename, "?"); queryIndex != -1 {
				origFilename = origFilename[:queryIndex]
			}
			// Remove the original extension and replace with .webp
			if extIndex := strings.LastIndex(origFilename, "."); extIndex != -1 {
				filename = origFilename[:extIndex] + ".webp"
			} else {
				filename = origFilename + ".webp"
			}
		}
		
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	}

	// Return the processed image
	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Content-Length", strconv.Itoa(len(processedData)))
	w.WriteHeader(http.StatusOK)
	w.Write(processedData)
}

// ProcessFromUpload handles image upload requests
func (api *ImageAPI) ProcessFromUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get processing options from query parameters
	options := getProcessOptionsFromRequest(r)

	// Get the image from the form data
	imageData, err := api.imageHandler.GetImageFromUpload(r, "image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload image: %v", err), http.StatusBadRequest)
		return
	}

	// Process the image
	processedData, metadata, err := api.processor.ProcessFromBytes(imageData, &options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process image: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if metadata response is requested
	if r.URL.Query().Get("metadata") == "true" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(metadata); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode metadata: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Set filename for download if requested
	if r.URL.Query().Get("download") == "true" || r.URL.Query().Get("format") == "webp" {
		// Use the original filename from the form but change extension to .webp
		filename := "image.webp"
		
		// Try to get the original filename from the form
		file, fileHeader, err := r.FormFile("image")
		if err == nil && fileHeader != nil {
			defer file.Close()
			if fileHeader.Filename != "" {
				// Remove the original extension and replace with .webp
				if extIndex := strings.LastIndex(fileHeader.Filename, "."); extIndex != -1 {
					filename = fileHeader.Filename[:extIndex] + ".webp"
				} else {
					filename = fileHeader.Filename + ".webp"
				}
			}
		}
		
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	}

	// Return the processed image
	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Content-Length", strconv.Itoa(len(processedData)))
	w.WriteHeader(http.StatusOK)
	w.Write(processedData)
}

// Health endpoint returns status OK if the service is running
func (api *ImageAPI) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

// getProcessOptionsFromRequest extracts processing options from request query parameters
func getProcessOptionsFromRequest(r *http.Request) processor.ProcessOptions {
	options := processor.ProcessOptions{
		MaxWidth:      1920, // Default max width
		MaxHeight:     1080, // Default max height
		Quality:       85,   // Default quality
		PreserveRatio: true, // Default preserve ratio
	}

	// Parse max width
	if maxWidth := r.URL.Query().Get("max_width"); maxWidth != "" {
		if width, err := strconv.Atoi(maxWidth); err == nil && width > 0 {
			options.MaxWidth = width
		}
	}

	// Parse max height
	if maxHeight := r.URL.Query().Get("max_height"); maxHeight != "" {
		if height, err := strconv.Atoi(maxHeight); err == nil && height > 0 {
			options.MaxHeight = height
		}
	}

	// Parse quality
	if quality := r.URL.Query().Get("quality"); quality != "" {
		if q, err := strconv.Atoi(quality); err == nil && q > 0 && q <= 100 {
			options.Quality = q
		}
	}

	// Parse preserve ratio
	if preserveRatio := r.URL.Query().Get("preserve_ratio"); preserveRatio != "" {
		options.PreserveRatio = preserveRatio != "false"
	}

	return options
} 
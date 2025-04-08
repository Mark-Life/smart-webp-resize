package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// Common errors
var (
	ErrEmptyURL          = errors.New("URL cannot be empty")
	ErrInvalidURL        = errors.New("invalid URL format")
	ErrHTTPRequestFailed = errors.New("HTTP request failed")
	ErrEmptyFile         = errors.New("file content is empty")
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrNoFile            = errors.New("no file found in request")
)

// ImageHandler handles image input from different sources
type ImageHandler interface {
	// GetImageFromURL fetches an image from a URL
	GetImageFromURL(url string) ([]byte, error)
	
	// ValidateURL checks if a URL is valid
	ValidateURL(url string) error
	
	// GetImageFromUpload extracts an image from an HTTP file upload
	GetImageFromUpload(r *http.Request, fieldName string) ([]byte, error)
	
	// ValidateFileType checks if the file has a valid image extension
	ValidateFileType(fileName string) error
}

// NewImageHandler creates a new image handler
func NewImageHandler() ImageHandler {
	return &defaultImageHandler{}
}

// defaultImageHandler is the default implementation of ImageHandler
type defaultImageHandler struct{}

// GetImageFromURL fetches an image from a URL
func (h *defaultImageHandler) GetImageFromURL(imageURL string) ([]byte, error) {
	if err := h.ValidateURL(imageURL); err != nil {
		return nil, err
	}
	
	// Create HTTP client with sensible defaults
	client := &http.Client{
		Timeout: http.DefaultClient.Timeout,
	}
	
	// Make the request
	resp, err := client.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHTTPRequestFailed, err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: server returned status %d", ErrHTTPRequestFailed, resp.StatusCode)
	}
	
	// Read the response body
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	if len(imageData) == 0 {
		return nil, ErrEmptyFile
	}
	
	return imageData, nil
}

// ValidateURL checks if a URL is valid
func (h *defaultImageHandler) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return ErrEmptyURL
	}
	
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}
	
	// Check if the URL has a scheme (http:// or https://)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: scheme must be http or https", ErrInvalidURL)
	}
	
	return nil
}

// GetImageFromUpload extracts an image from an HTTP file upload
func (h *defaultImageHandler) GetImageFromUpload(r *http.Request, fieldName string) ([]byte, error) {
	// Parse the multipart form, with a reasonable max memory
	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}
	
	// Get the file from the form
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoFile, err)
	}
	defer file.Close()
	
	// Validate the file type
	if err := h.ValidateFileType(header.Filename); err != nil {
		return nil, err
	}
	
	// Read the file
	imageData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	if len(imageData) == 0 {
		return nil, ErrEmptyFile
	}
	
	return imageData, nil
}

// ValidateFileType checks if the file has a valid image extension
func (h *defaultImageHandler) ValidateFileType(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("%w: filename is empty", ErrInvalidFileType)
	}
	
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return fmt.Errorf("%w: file has no extension", ErrInvalidFileType)
	}
	
	// List of supported image formats
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
	}
	
	if !validExtensions[ext] {
		return fmt.Errorf("%w: extension %s is not supported", ErrInvalidFileType, ext)
	}
	
	return nil
} 
package processor

import (
	"errors"

	"github.com/Mark-Life/smart-webp-resize/pkg/models"
)

// Common errors
var (
	ErrProcessingFailed = errors.New("image processing failed")
	ErrInvalidImage     = errors.New("invalid or unsupported image format")
	ErrResizingFailed   = errors.New("image resizing failed")
	ErrEncodingFailed   = errors.New("WebP encoding failed")
)

// ImageProcessor defines the interface for processing images
type ImageProcessor interface {
	// ProcessFromURL processes an image from a URL
	ProcessFromURL(url string, options *ProcessOptions) ([]byte, *models.ImageMetadata, error)
	
	// ProcessFromBytes processes an image from bytes
	ProcessFromBytes(imageData []byte, options *ProcessOptions) ([]byte, *models.ImageMetadata, error)
}

// ProcessOptions contains options for image processing
type ProcessOptions struct {
	MaxWidth      int
	MaxHeight     int
	Quality       int
	PreserveRatio bool
}

// New creates a new image processor with default settings
func New() ImageProcessor {
	return &defaultProcessor{}
}

// defaultProcessor is the default implementation of ImageProcessor
type defaultProcessor struct{}

// ProcessFromURL implements the ImageProcessor interface
func (p *defaultProcessor) ProcessFromURL(url string, options *ProcessOptions) ([]byte, *models.ImageMetadata, error) {
	// To be implemented
	return nil, nil, nil
}

// ProcessFromBytes implements the ImageProcessor interface
func (p *defaultProcessor) ProcessFromBytes(imageData []byte, options *ProcessOptions) ([]byte, *models.ImageMetadata, error) {
	// To be implemented
	return nil, nil, nil
} 
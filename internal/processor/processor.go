package processor

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"io"
	"net/http"
	"time"

	"github.com/Mark-Life/smart-webp-resize/pkg/models"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp" // Register BMP format
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
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Fetch the image
	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("failed to fetch image, status: " + resp.Status)
	}
	
	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	
	// Process the image data
	return p.ProcessFromBytes(imageData, options)
}

// ProcessFromBytes implements the ImageProcessor interface
func (p *defaultProcessor) ProcessFromBytes(imageData []byte, options *ProcessOptions) ([]byte, *models.ImageMetadata, error) {
	// Set default options if none provided
	if options == nil {
		options = &ProcessOptions{
			MaxWidth:      1920,
			MaxHeight:     1080,
			Quality:       80,
			PreserveRatio: true,
		}
	}

	// Detect image format
	format, err := p.detectImageFormat(imageData)
	if err != nil {
		return nil, nil, err
	}
	
	// Decode the image
	img, err := p.decodeImage(imageData)
	if err != nil {
		return nil, nil, err
	}
	
	// Get original dimensions and size
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()
	originalSize := int64(len(imageData))
	
	// Calculate new dimensions
	newWidth, newHeight := p.calculateDimensions(
		originalWidth, 
		originalHeight, 
		options.MaxWidth, 
		options.MaxHeight,
	)
	
	// Resize the image
	resizedImg, err := p.resizeImage(img, newWidth, newHeight)
	if err != nil {
		return nil, nil, err
	}
	
	// Encode to WebP
	webpData, err := p.encodeToWebP(resizedImg, options.Quality)
	if err != nil {
		return nil, nil, err
	}
	
	// Create metadata
	sizeReduction := int(100 * (originalSize - int64(len(webpData))) / originalSize)
	metadata := &models.ImageMetadata{
		OriginalWidth:  originalWidth,
		OriginalHeight: originalHeight,
		OriginalFormat: format,
		OriginalSize:   originalSize,
		NewWidth:       newWidth,
		NewHeight:      newHeight,
		NewFormat:      "webp",
		NewSize:        int64(len(webpData)),
		SizeReduction:  sizeReduction,
	}
	
	return webpData, metadata, nil
}

// detectImageFormat determines the image format from the image data
func (p *defaultProcessor) detectImageFormat(data []byte) (string, error) {
	if len(data) < 8 {
		return "", ErrInvalidImage
	}
	
	// Check for JPEG signature (FF D8 FF)
	// JPEG format starts with FF D8 and has an FF marker
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg", nil
	}
	
	// Check for PNG signature
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "png", nil
	}
	
	// Check for BMP signature
	if len(data) >= 2 && data[0] == 0x42 && data[1] == 0x4D {
		return "bmp", nil
	}
	
	// Check for WebP signature (RIFF....WEBP)
	if len(data) >= 12 && bytes.HasPrefix(data, []byte{0x52, 0x49, 0x46, 0x46}) &&
		bytes.Equal(data[8:12], []byte{0x57, 0x45, 0x42, 0x50}) {
		return "webp", nil
	}
	
	return "", ErrInvalidImage
}

// decodeImage decodes image data into an image.Image
func (p *defaultProcessor) decodeImage(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		// Try specific decoders if the standard one fails
		reader.Reset(data)
		
		// Try BMP decoder
		if img, err = bmp.Decode(reader); err == nil {
			return img, nil
		}
		
		// Try WebP decoder
		reader.Reset(data)
		if img, err = webp.Decode(reader); err == nil {
			return img, nil
		}
		
		return nil, ErrInvalidImage
	}
	
	return img, nil
}

// calculateDimensions calculates new dimensions while preserving aspect ratio
func (p *defaultProcessor) calculateDimensions(originalWidth, originalHeight, maxWidth, maxHeight int) (int, int) {
	// If the image is already smaller than the max dimensions, no need to resize
	if originalWidth <= maxWidth && originalHeight <= maxHeight {
		return originalWidth, originalHeight
	}
	
	// Calculate scaling factors for width and height
	widthScale := float64(maxWidth) / float64(originalWidth)
	heightScale := float64(maxHeight) / float64(originalHeight)
	
	// Use the smaller scale to ensure the image fits within both dimensions
	scale := widthScale
	if heightScale < widthScale {
		scale = heightScale
	}
	
	// Calculate new dimensions - round to nearest integer
	newWidth := int(float64(originalWidth) * scale + 0.5)
	newHeight := int(float64(originalHeight) * scale + 0.5)
	
	// Ensure at least 1 pixel in each dimension
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}
	
	return newWidth, newHeight
}

// resizeImage resizes the image to the specified dimensions
func (p *defaultProcessor) resizeImage(img image.Image, width, height int) (image.Image, error) {
	// Use Lanczos resampling for high quality
	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	if resized == nil {
		return nil, ErrResizingFailed
	}
	return resized, nil
}

// encodeToWebP encodes the image to WebP format with the specified quality
func (p *defaultProcessor) encodeToWebP(img image.Image, quality int) ([]byte, error) {
	// Adjust quality to be within valid range (0-100)
	if quality < 0 {
		quality = 0
	} else if quality > 100 {
		quality = 100
	}
	
	// Encode to WebP
	var buf bytes.Buffer
	err := webp.Encode(&buf, img, &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	})
	if err != nil {
		return nil, ErrEncodingFailed
	}
	
	return buf.Bytes(), nil
} 
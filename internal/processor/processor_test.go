package processor

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessorCreation(t *testing.T) {
	processor := New()
	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}
}

func TestCalculateDimensions(t *testing.T) {
	tests := []struct {
		name              string
		originalWidth     int
		originalHeight    int
		maxWidth          int
		maxHeight         int
		expectedNewWidth  int
		expectedNewHeight int
	}{
		{
			name:              "Landscape image within limits",
			originalWidth:     1000,
			originalHeight:    600,
			maxWidth:          1920,
			maxHeight:         1080,
			expectedNewWidth:  1000,
			expectedNewHeight: 600,
		},
		{
			name:              "Landscape image exceeding width",
			originalWidth:     2500,
			originalHeight:    1500,
			maxWidth:          1920,
			maxHeight:         1080,
			expectedNewWidth:  1800,
			expectedNewHeight: 1080,
		},
		{
			name:              "Portrait image exceeding height",
			originalWidth:     1000,
			originalHeight:    2000,
			maxWidth:          1920,
			maxHeight:         1080,
			expectedNewWidth:  540,
			expectedNewHeight: 1080,
		},
		{
			name:              "Tiny image",
			originalWidth:     50,
			originalHeight:    30,
			maxWidth:          1920,
			maxHeight:         1080,
			expectedNewWidth:  50,
			expectedNewHeight: 30,
		},
		{
			name:              "Square image",
			originalWidth:     2000,
			originalHeight:    2000,
			maxWidth:          1920,
			maxHeight:         1080,
			expectedNewWidth:  1080,
			expectedNewHeight: 1080,
		},
	}

	processor := &defaultProcessor{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newWidth, newHeight := processor.calculateDimensions(
				tt.originalWidth,
				tt.originalHeight,
				tt.maxWidth,
				tt.maxHeight,
			)

			if newWidth != tt.expectedNewWidth || newHeight != tt.expectedNewHeight {
				t.Errorf("Expected dimensions %dx%d, got %dx%d",
					tt.expectedNewWidth, tt.expectedNewHeight, newWidth, newHeight)
			}
		})
	}
}

func TestDetectImageFormat(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expectedFormat string
		expectError    bool
	}{
		{
			name:           "JPEG image",
			data:           []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46},
			expectedFormat: "jpeg",
			expectError:    false,
		},
		{
			name:           "PNG image",
			data:           []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expectedFormat: "png",
			expectError:    false,
		},
		{
			name:           "BMP image",
			data:           []byte{0x42, 0x4D, 0x76, 0x38, 0x00, 0x00, 0x00, 0x00},
			expectedFormat: "bmp",
			expectError:    false,
		},
		{
			name:           "WEBP image",
			data:           []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50},
			expectedFormat: "webp",
			expectError:    false,
		},
		{
			name:           "Unsupported format",
			data:           []byte{0x00, 0x01, 0x02, 0x03},
			expectedFormat: "",
			expectError:    true,
		},
		{
			name:           "Empty data",
			data:           []byte{},
			expectedFormat: "",
			expectError:    true,
		},
	}

	processor := &defaultProcessor{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, err := processor.detectImageFormat(tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if format != tt.expectedFormat {
				t.Errorf("Expected format %s, got %s", tt.expectedFormat, format)
			}
		})
	}
}

func TestDecodeImage(t *testing.T) {
	// Create a simple test image
	img := createTestImage(100, 100)
	
	// Encode it to PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
	
	// Get the PNG data
	pngData := buf.Bytes()
	
	processor := &defaultProcessor{}
	
	// Test decoding
	decodedImg, err := processor.decodeImage(pngData)
	if err != nil {
		t.Errorf("Failed to decode valid image: %v", err)
	}
	
	if decodedImg == nil {
		t.Error("Decoded image is nil")
		return
	}
	
	// Check dimensions
	bounds := decodedImg.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Decoded image has wrong dimensions: got %dx%d, want 100x100", 
			bounds.Dx(), bounds.Dy())
	}
	
	// Test invalid data
	_, err = processor.decodeImage([]byte{0, 1, 2, 3})
	if err == nil {
		t.Error("Expected error when decoding invalid image data, got nil")
	}
}

// Helper function to create a test image
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Fill with a gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 255 / width),
				G: uint8(y * 255 / height),
				B: 100,
				A: 255,
			})
		}
	}
	
	return img
}

func TestResizeImage(t *testing.T) {
	// Create a test image of 200x100
	img := createTestImage(200, 100)
	
	processor := &defaultProcessor{}
	
	// Test resize to smaller dimensions
	resized, err := processor.resizeImage(img, 100, 50)
	if err != nil {
		t.Errorf("Failed to resize image: %v", err)
	}
	
	if resized == nil {
		t.Error("Resized image is nil")
		return
	}
	
	bounds := resized.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 50 {
		t.Errorf("Resized image has wrong dimensions: got %dx%d, want 100x50", 
			bounds.Dx(), bounds.Dy())
	}
	
	// Test resize to larger dimensions
	resized, err = processor.resizeImage(img, 400, 200)
	if err != nil {
		t.Errorf("Failed to resize image: %v", err)
	}
	
	if resized == nil {
		t.Error("Resized image is nil")
		return
	}
	
	bounds = resized.Bounds()
	if bounds.Dx() != 400 || bounds.Dy() != 200 {
		t.Errorf("Resized image has wrong dimensions: got %dx%d, want 400x200", 
			bounds.Dx(), bounds.Dy())
	}
}

func TestEncodeToWebP(t *testing.T) {
	// Create a test image
	img := createTestImage(100, 100)
	
	processor := &defaultProcessor{}
	
	// Test with default quality
	webpData, err := processor.encodeToWebP(img, 80)
	if err != nil {
		t.Errorf("Failed to encode to WebP: %v", err)
	}
	
	if len(webpData) == 0 {
		t.Error("WebP data is empty")
	}
	
	// Test with low quality
	lowQualityData, err := processor.encodeToWebP(img, 10)
	if err != nil {
		t.Errorf("Failed to encode to WebP with low quality: %v", err)
	}
	
	// Test with high quality
	highQualityData, err := processor.encodeToWebP(img, 90)
	if err != nil {
		t.Errorf("Failed to encode to WebP with high quality: %v", err)
	}
	
	// Lower quality should generally result in smaller file size
	if len(lowQualityData) >= len(highQualityData) {
		t.Log("Expected lower quality to produce smaller file, but it didn't. This might be OK for very small test images.")
	}
}

func TestProcessFromBytes(t *testing.T) {
	// Create a test image
	img := createTestImage(500, 300)
	
	// Encode it to PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
	
	// Get the PNG data
	pngData := buf.Bytes()
	
	processor := &defaultProcessor{}
	
	// Test processing with default options
	options := &ProcessOptions{
		MaxWidth:      300,
		MaxHeight:     200,
		Quality:       80,
		PreserveRatio: true,
	}
	
	webpData, metadata, err := processor.ProcessFromBytes(pngData, options)
	if err != nil {
		t.Errorf("Failed to process image: %v", err)
	}
	
	if webpData == nil {
		t.Error("WebP data is nil")
		return
	}
	
	// Check metadata
	if metadata.OriginalWidth != 500 || metadata.OriginalHeight != 300 {
		t.Errorf("Wrong original dimensions in metadata: got %dx%d, want 500x300", 
			metadata.OriginalWidth, metadata.OriginalHeight)
	}
	
	if metadata.NewWidth != 300 || metadata.NewHeight != 180 {
		t.Errorf("Wrong new dimensions in metadata: got %dx%d, want 300x180", 
			metadata.NewWidth, metadata.NewHeight)
	}
	
	if metadata.OriginalFormat != "png" {
		t.Errorf("Wrong original format in metadata: got %s, want png", metadata.OriginalFormat)
	}
	
	if metadata.NewFormat != "webp" {
		t.Errorf("Wrong new format in metadata: got %s, want webp", metadata.NewFormat)
	}
}

func TestProcessFromURL(t *testing.T) {
	// Create a test image
	img := createTestImage(500, 300)
	
	// Encode it to PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
	
	// Set up a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(buf.Bytes())
	}))
	defer server.Close()
	
	processor := &defaultProcessor{}
	
	// Test processing with default options
	options := &ProcessOptions{
		MaxWidth:      300,
		MaxHeight:     200,
		Quality:       80,
		PreserveRatio: true,
	}
	
	webpData, metadata, err := processor.ProcessFromURL(server.URL, options)
	if err != nil {
		t.Errorf("Failed to process image from URL: %v", err)
	}
	
	if webpData == nil {
		t.Error("WebP data is nil")
		return
	}
	
	// Check metadata
	if metadata.OriginalWidth != 500 || metadata.OriginalHeight != 300 {
		t.Errorf("Wrong original dimensions in metadata: got %dx%d, want 500x300", 
			metadata.OriginalWidth, metadata.OriginalHeight)
	}
	
	if metadata.NewWidth != 300 || metadata.NewHeight != 180 {
		t.Errorf("Wrong new dimensions in metadata: got %dx%d, want 300x180", 
			metadata.NewWidth, metadata.NewHeight)
	}
	
	if metadata.OriginalFormat != "png" {
		t.Errorf("Wrong original format in metadata: got %s, want png", metadata.OriginalFormat)
	}
	
	if metadata.NewFormat != "webp" {
		t.Errorf("Wrong new format in metadata: got %s, want webp", metadata.NewFormat)
	}
	
	// Test with invalid URL
	_, _, err = processor.ProcessFromURL("http://invalid-domain-that-should-not-exist.xyz", options)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
} 
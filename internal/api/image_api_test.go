package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mark-Life/smart-webp-resize/internal/handler"
	"github.com/Mark-Life/smart-webp-resize/internal/processor"
	"github.com/Mark-Life/smart-webp-resize/pkg/models"
)

// MockImageHandler implements handler.ImageHandler for testing
type MockImageHandler struct {
	GetURLFunc      func(url string) ([]byte, error)
	ValidateURLFunc func(url string) error
	GetUploadFunc   func(r *http.Request, fieldName string) ([]byte, error)
	ValidateTypeFunc func(fileName string) error
}

func (m *MockImageHandler) GetImageFromURL(url string) ([]byte, error) {
	return m.GetURLFunc(url)
}

func (m *MockImageHandler) ValidateURL(url string) error {
	return m.ValidateURLFunc(url)
}

func (m *MockImageHandler) GetImageFromUpload(r *http.Request, fieldName string) ([]byte, error) {
	return m.GetUploadFunc(r, fieldName)
}

func (m *MockImageHandler) ValidateFileType(fileName string) error {
	return m.ValidateTypeFunc(fileName)
}

// MockImageProcessor implements processor.ImageProcessor for testing
type MockImageProcessor struct {
	ProcessURLFunc   func(url string, options *processor.ProcessOptions) ([]byte, *models.ImageMetadata, error)
	ProcessBytesFunc func(imageData []byte, options *processor.ProcessOptions) ([]byte, *models.ImageMetadata, error)
}

func (m *MockImageProcessor) ProcessFromURL(url string, options *processor.ProcessOptions) ([]byte, *models.ImageMetadata, error) {
	return m.ProcessURLFunc(url, options)
}

func (m *MockImageProcessor) ProcessFromBytes(imageData []byte, options *processor.ProcessOptions) ([]byte, *models.ImageMetadata, error) {
	return m.ProcessBytesFunc(imageData, options)
}

func TestHealthEndpoint(t *testing.T) {
	mockHandler := &MockImageHandler{}
	mockProcessor := &MockImageProcessor{}
	api := NewImageAPI(mockHandler, mockProcessor)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	api.Health(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result map[string]string
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["status"] != "OK" {
		t.Errorf("Expected status to be OK, got %v", result["status"])
	}
}

func TestProcessFromURL(t *testing.T) {
	// Mock image data and metadata
	mockImageData := []byte("original image data")
	mockProcessedData := []byte("processed webp data")
	mockMetadata := &models.ImageMetadata{
		OriginalWidth:  500,
		OriginalHeight: 300,
		OriginalFormat: "jpeg",
		OriginalSize:   1000,
		NewWidth:       400,
		NewHeight:      240,
		NewFormat:      "webp",
		NewSize:        500,
		SizeReduction:  50,
	}

	// Set up mock handlers
	mockHandler := &MockImageHandler{
		GetURLFunc: func(url string) ([]byte, error) {
			if url == "http://example.com/image.jpg" {
				return mockImageData, nil
			}
			return nil, handler.ErrInvalidURL
		},
		ValidateURLFunc: func(url string) error {
			if url == "http://example.com/image.jpg" {
				return nil
			}
			return handler.ErrInvalidURL
		},
	}

	mockProcessor := &MockImageProcessor{
		ProcessBytesFunc: func(imageData []byte, options *processor.ProcessOptions) ([]byte, *models.ImageMetadata, error) {
			if bytes.Equal(imageData, mockImageData) {
				return mockProcessedData, mockMetadata, nil
			}
			return nil, nil, processor.ErrProcessingFailed
		},
	}

	api := NewImageAPI(mockHandler, mockProcessor)

	// Test successful image processing
	t.Run("successful image processing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/process?url=http://example.com/image.jpg", nil)
		w := httptest.NewRecorder()

		api.ProcessFromURL(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		if resp.Header.Get("Content-Type") != "image/webp" {
			t.Errorf("Expected Content-Type image/webp, got %v", resp.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if !bytes.Equal(body, mockProcessedData) {
			t.Errorf("Response body does not match expected processed data")
		}
	})

	// Test metadata response
	t.Run("metadata response", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/process?url=http://example.com/image.jpg&metadata=true", nil)
		w := httptest.NewRecorder()

		api.ProcessFromURL(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		if resp.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %v", resp.Header.Get("Content-Type"))
		}

		var metadata models.ImageMetadata
		err := json.NewDecoder(resp.Body).Decode(&metadata)
		if err != nil {
			t.Fatalf("Failed to decode metadata: %v", err)
		}

		if metadata.OriginalWidth != mockMetadata.OriginalWidth ||
			metadata.NewWidth != mockMetadata.NewWidth ||
			metadata.SizeReduction != mockMetadata.SizeReduction {
			t.Errorf("Metadata does not match expected values")
		}
	})

	// Test missing URL parameter
	t.Run("missing URL parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/process", nil)
		w := httptest.NewRecorder()

		api.ProcessFromURL(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status Bad Request, got %v", resp.Status)
		}
	})

	// Test invalid method
	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/process?url=http://example.com/image.jpg", nil)
		w := httptest.NewRecorder()

		api.ProcessFromURL(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status Method Not Allowed, got %v", resp.Status)
		}
	})
}

func TestProcessFromUpload(t *testing.T) {
	// Mock image data and metadata
	mockImageData := []byte("original image data")
	mockProcessedData := []byte("processed webp data")
	mockMetadata := &models.ImageMetadata{
		OriginalWidth:  500,
		OriginalHeight: 300,
		OriginalFormat: "jpeg",
		OriginalSize:   1000,
		NewWidth:       400,
		NewHeight:      240,
		NewFormat:      "webp",
		NewSize:        500,
		SizeReduction:  50,
	}

	// Set up mock handlers
	mockHandler := &MockImageHandler{
		GetUploadFunc: func(r *http.Request, fieldName string) ([]byte, error) {
			if fieldName == "image" {
				return mockImageData, nil
			}
			return nil, handler.ErrNoFile
		},
		ValidateTypeFunc: func(fileName string) error {
			if strings.HasSuffix(fileName, ".jpg") {
				return nil
			}
			return handler.ErrInvalidFileType
		},
	}

	mockProcessor := &MockImageProcessor{
		ProcessBytesFunc: func(imageData []byte, options *processor.ProcessOptions) ([]byte, *models.ImageMetadata, error) {
			if bytes.Equal(imageData, mockImageData) {
				return mockProcessedData, mockMetadata, nil
			}
			return nil, nil, processor.ErrProcessingFailed
		},
	}

	api := NewImageAPI(mockHandler, mockProcessor)

	// Helper function to create a multipart form request
	createMultipartRequest := func(fieldName, fileName string, fileContent []byte) (*http.Request, error) {
		var requestBody bytes.Buffer
		multipartWriter := multipart.NewWriter(&requestBody)
		
		fileWriter, err := multipartWriter.CreateFormFile(fieldName, fileName)
		if err != nil {
			return nil, err
		}
		
		_, err = fileWriter.Write(fileContent)
		if err != nil {
			return nil, err
		}
		
		err = multipartWriter.Close()
		if err != nil {
			return nil, err
		}
		
		req := httptest.NewRequest("POST", "/upload", &requestBody)
		req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
		return req, nil
	}

	// Test successful upload
	t.Run("successful upload", func(t *testing.T) {
		req, err := createMultipartRequest("image", "test.jpg", mockImageData)
		if err != nil {
			t.Fatalf("Failed to create multipart request: %v", err)
		}
		
		w := httptest.NewRecorder()
		api.ProcessFromUpload(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		if resp.Header.Get("Content-Type") != "image/webp" {
			t.Errorf("Expected Content-Type image/webp, got %v", resp.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if !bytes.Equal(body, mockProcessedData) {
			t.Errorf("Response body does not match expected processed data")
		}
	})

	// Test invalid method
	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/upload", nil)
		w := httptest.NewRecorder()

		api.ProcessFromUpload(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status Method Not Allowed, got %v", resp.Status)
		}
	})
} 
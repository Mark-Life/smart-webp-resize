package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetImageFromURL(t *testing.T) {
	// Create a mock server to return test images
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/validimage.jpg":
			// Return a mock image data (just some bytes for the test)
			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("mock image data"))
		case "/errorimage":
			// Return an error status
			w.WriteHeader(http.StatusInternalServerError)
		default:
			// Return not found
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	tests := []struct {
		name          string
		url           string
		wantErr       bool
		expectedBytes []byte
	}{
		{
			name:          "valid image url",
			url:           mockServer.URL + "/validimage.jpg",
			wantErr:       false,
			expectedBytes: []byte("mock image data"),
		},
		{
			name:    "not found url",
			url:     mockServer.URL + "/notfound.jpg",
			wantErr: true,
		},
		{
			name:    "server error",
			url:     mockServer.URL + "/errorimage",
			wantErr: true,
		},
		{
			name:    "invalid url format",
			url:     "http://invalid\\url",
			wantErr: true,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
	}

	handler := NewImageHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := handler.GetImageFromURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(data) != string(tt.expectedBytes) {
				t.Errorf("GetImageFromURL() got unexpected data")
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http url",
			url:     "http://example.com/image.jpg",
			wantErr: false,
		},
		{
			name:    "valid https url",
			url:     "https://example.com/image.jpg",
			wantErr: false,
		},
		{
			name:    "invalid url format",
			url:     "http://invalid\\url",
			wantErr: true,
		},
		{
			name:    "missing scheme",
			url:     "example.com/image.jpg",
			wantErr: true,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
	}

	handler := NewImageHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} 
package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetImageFromUpload(t *testing.T) {
	// Create a test file upload request
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

	tests := []struct {
		name          string
		fieldName     string
		fileName      string
		fileContent   []byte
		wantErr       bool
		expectedBytes []byte
	}{
		{
			name:          "valid image upload",
			fieldName:     "image",
			fileName:      "test.jpg",
			fileContent:   []byte("mock image data"),
			wantErr:       false,
			expectedBytes: []byte("mock image data"),
		},
		{
			name:        "empty file",
			fieldName:   "image",
			fileName:    "empty.jpg",
			fileContent: []byte{},
			wantErr:     true,
		},
	}

	handler := NewImageHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := createMultipartRequest(tt.fieldName, tt.fileName, tt.fileContent)
			if err != nil {
				t.Fatalf("Failed to create multipart request: %v", err)
			}

			data, err := handler.GetImageFromUpload(req, tt.fieldName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageFromUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(data) != string(tt.expectedBytes) {
				t.Errorf("GetImageFromUpload() got unexpected data")
			}
		})
	}
}

func TestValidateFileType(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{
			name:     "valid jpeg",
			fileName: "image.jpg",
			wantErr:  false,
		},
		{
			name:     "valid jpeg uppercase",
			fileName: "IMAGE.JPG",
			wantErr:  false,
		},
		{
			name:     "valid png",
			fileName: "image.png",
			wantErr:  false,
		},
		{
			name:     "valid gif",
			fileName: "image.gif",
			wantErr:  false,
		},
		{
			name:     "invalid extension",
			fileName: "document.pdf",
			wantErr:  true,
		},
		{
			name:     "no extension",
			fileName: "image",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			fileName: "",
			wantErr:  true,
		},
	}

	handler := NewImageHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateFileType(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} 
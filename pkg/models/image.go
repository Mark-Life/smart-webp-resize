package models

// ImageSource represents different ways to provide an image
type ImageSource int

const (
	// ImageSourceUpload is for direct file uploads
	ImageSourceUpload ImageSource = iota
	// ImageSourceURL is for images fetched from a URL
	ImageSourceURL
)

// ResizeRequest contains details about an image resize request
type ResizeRequest struct {
	Source        ImageSource
	URL           string
	MaxWidth      int
	MaxHeight     int
	Quality       int
	PreserveRatio bool
}

// ImageMetadata contains information about processed images
type ImageMetadata struct {
	OriginalWidth  int    `json:"original_width"`
	OriginalHeight int    `json:"original_height"`
	OriginalFormat string `json:"original_format"`
	OriginalSize   int64  `json:"original_size"`
	NewWidth       int    `json:"new_width"`
	NewHeight      int    `json:"new_height"`
	NewFormat      string `json:"new_format"`
	NewSize        int64  `json:"new_size"`
	SizeReduction  int    `json:"size_reduction_percent"`
} 
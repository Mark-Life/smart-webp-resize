# Smart WebP Image Resizer

A serverless function that automatically resizes images to optimal web dimensions and converts them to WebP format for improved performance.

## Features

- Accept images via direct upload or URL
- Smart resizing while preserving aspect ratio
- WebP conversion for optimized file size
- Metadata about original and processed images

## Development

### Prerequisites

- Go 1.24 or higher
- Git

### Setup

1. Clone the repository

   ```
   git clone https://github.com/Mark-Life/smart-webp-resize.git
   cd smart-webp-resize
   ```

2. Start the development server
   ```
   go run cmd/server/main.go
   ```

### API Endpoints

- `GET /health` - Health check endpoint
- `POST /resize` - Resize and convert an image

## License

MIT

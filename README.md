# Smart WebP Image Resizer

A serverless image processing service that automatically resizes images to web-optimal dimensions and converts them to WebP format for improved performance. Supports input via direct upload or URL.

## Features

- ✅ Accept images via direct upload or URL
- ✅ Smart resizing while preserving aspect ratio
- ✅ WebP conversion for optimized file size and quality
- ✅ Detailed metadata about original and processed images
- ✅ Support for JPEG, PNG, and BMP input formats
- ✅ Customizable output dimensions and quality settings
- ✅ Serverless-ready architecture

## API Endpoints

### Health Check

```
GET /health
```

Returns a simple status indicating the service is running.

**Response**:

```json
{
  "status": "OK"
}
```

### Process Image from URL

```
GET/POST /process/url
```

Processes an image from a URL.

**Parameters**:

- `url` (required): URL of the image to process
- `max_width` (optional): Maximum width of the output image (default: 1920)
- `max_height` (optional): Maximum height of the output image (default: 1080)
- `quality` (optional): WebP quality level (1-100, default: 85)
- `preserve_ratio` (optional): Whether to preserve aspect ratio (default: true)
- `metadata` (optional): If set to "true", returns only metadata instead of the image

**Response**:

- If `metadata=true`: JSON metadata about the image processing
- Otherwise: WebP image

### Process Uploaded Image

```
POST /process/upload
```

Processes an uploaded image.

**Parameters**:

- `image` (required): Image file upload (multipart/form-data)
- `max_width` (optional): Maximum width of the output image (default: 1920)
- `max_height` (optional): Maximum height of the output image (default: 1080)
- `quality` (optional): WebP quality level (1-100, default: 85)
- `preserve_ratio` (optional): Whether to preserve aspect ratio (default: true)
- `metadata` (optional): If set to "true", returns only metadata instead of the image

**Response**:

- If `metadata=true`: JSON metadata about the image processing
- Otherwise: WebP image

### Metadata Response Example

```json
{
  "original_width": 2500,
  "original_height": 1500,
  "original_format": "jpeg",
  "original_size": 986543,
  "new_width": 1800,
  "new_height": 1080,
  "new_format": "webp",
  "new_size": 124567,
  "size_reduction_percent": 87
}
```

## Usage Examples

### Using cURL

Process an image from URL:

```bash
curl -X GET "http://localhost:8080/process/url?url=https://example.com/image.jpg&max_width=800&max_height=600&quality=90" --output processed.webp
```

Get only metadata:

```bash
curl -X GET "http://localhost:8080/process/url?url=https://example.com/image.jpg&metadata=true"
```

Upload and process an image:

```bash
curl -X POST -F "image=@/path/to/image.jpg" -F "max_width=800" -F "max_height=600" -F "quality=90" http://localhost:8080/process/upload --output processed.webp
```

### Using the Web Interface

A test web interface is available at:

```
http://localhost:8080/
```

This interface allows you to:

- Upload images directly from your computer
- Process images from URLs
- Customize resizing and quality parameters
- View before and after results
- See detailed metadata

## Development

### Prerequisites

- Go 1.24 or higher
- Git

### Setup

1. Clone the repository

   ```bash
   git clone https://github.com/yourusername/smart-webp-resizer.git
   cd smart-webp-resizer
   ```

2. Install dependencies

   ```bash
   go mod tidy
   ```

3. Run the development server
   ```bash
   ./scripts/run.sh
   ```

### Project Structure

- `cmd/`: Application entry points
- `internal/`: Private application code
  - `api/`: HTTP API implementation
  - `handler/`: Image input handling
  - `processor/`: Core image processing logic
- `pkg/`: Public libraries and models
- `scripts/`: Build and deployment scripts
- `test/`: Test files and test utilities

## Deployment

### Local Deployment

Run the server locally:

```bash
./scripts/run.sh
```

### AWS Lambda Deployment

1. Update AWS Lambda settings in `scripts/deploy-lambda.sh`
2. Run the deployment script:
   ```bash
   ./scripts/deploy-lambda.sh
   ```

## Environment Variables

- `PORT`: HTTP server port (default: 8080)
- `MAX_WIDTH`: Default maximum width (default: 1920)
- `MAX_HEIGHT`: Default maximum height (default: 1080)
- `DEFAULT_QUALITY`: Default WebP quality (default: 85)

## Performance Considerations

- Image processing is memory-intensive, so configure enough memory for your serverless function
- Processing time increases with image size and quality settings
- For very large images, consider setting smaller max dimensions

## License

MIT

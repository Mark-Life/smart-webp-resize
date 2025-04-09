# Smart WebP Image Resizer

A serverless image processing service that automatically resizes images to web-optimal dimensions and converts them to WebP format for improved performance. Supports input via direct upload or URL.

- Backend: Go
- Frontend: Next.js
- Hosting: Google Cloud Run

[webp.mark-life.com](https://webp.mark-life.com/)

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

## React Frontend

The project now includes a modern React frontend built with Next.js that provides an improved user experience for image processing.

### Frontend Features

- Drag-and-drop interface for image uploads
- URL-based image processing
- Customizable settings for image dimensions and quality
- Before/after image comparison
- Detailed metadata display
- Responsive design

### Frontend Development

To work on the frontend separately:

```bash
cd frontend
pnpm install
pnpm dev
```

The development server will start at http://localhost:3000 and will proxy API requests to your Go backend.

## Docker Development

The project includes Docker configurations for both development and production:

### Using Docker Compose for Development

Start the complete environment with frontend and backend:

```bash
# Start all services (frontend, backend, and combined app)
docker-compose up

# Start in detached mode (run in background)
docker-compose up -d
```

This will start:

- The combined application at http://localhost:8080
- The frontend development server at http://localhost:3000
- The backend development server at http://localhost:8081

For individual services:

```bash
# Frontend development only
docker-compose up frontend-dev

# Backend development only
docker-compose up backend-dev
```

### Restart Services After Code Changes

```bash
# Restart and rebuild the frontend after code changes
docker-compose stop frontend-dev
docker-compose up --build frontend-dev

# Restart and rebuild the backend after code changes
docker-compose stop backend-dev
docker-compose up --build backend-dev

# Restart everything and rebuild all containers
docker-compose down
docker-compose up --build
```

### View Logs

```bash
# View logs for a specific service
docker-compose logs frontend-dev
docker-compose logs backend-dev

# Follow logs in real-time
docker-compose logs -f frontend-dev
```

### Building the Production Docker Image

Build and run the production Docker image:

```bash
docker build -t webp-resizer .
docker run -p 8080:8080 webp-resizer
```

## Deployment Options

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

### Google Cloud Run Deployment

1. Update the configuration in `scripts/deploy-cloud-run.sh`
2. Run the deployment script:
   ```bash
   ./scripts/deploy-cloud-run.sh
   ```

This will:

- Build and push the Docker image to Google Container Registry
- Deploy to Google Cloud Run with appropriate resource constraints
- Set up authentication and rate limiting for cost management
- Create a service account for secure access

## Cost Management for Cloud Deployment

The Google Cloud Run deployment includes several cost management features:

1. **Resource Limits**: Configures memory, CPU, and maximum instances
2. **Rate Limiting**: Limits the number of requests per minute
3. **Authentication**: Requires authentication to prevent unauthorized use
4. **Instance Auto-scaling**: Scales down to zero when not in use

You can adjust these settings in the `scripts/deploy-cloud-run.sh` script.

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

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

The MIT License is a permissive license that allows for reuse with very few restrictions. It permits users to:

- Use the code commercially
- Modify the code
- Distribute the code
- Use the code privately
- Sublicense the code

The only requirement is that the original copyright and license notice be included in any copy of the software/source.

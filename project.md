**Project: Smart WebP Image Resizer**

**Description:**

This serverless function takes an image as input (uploaded directly or referenced by URL), automatically resizes it to a sensible standard resolution, converts it to WebP format for optimal web performance, and returns the converted image. The goal is to provide a simple, automated tool that web developers can use to easily optimize images for their websites.

**Functionality:**

1.  **Input:**
    - Accepts images in common formats: JPEG, PNG, (BMP, TIFF if you want to get fancy)
    - Image source:
      - Direct upload (multipart form data)
      - URL (passed as a query parameter or in the request body)
2.  **Image Decoding:**
    - Detect the image format and decode the image data.
3.  **Smart Resizing:**

    - **Determine the Optimal Resolution:** The service will resize the image based on its original dimensions.

      - **Longest Side Constraint:** Resize the image so that its longest side (width or height) is equal to a pre-defined maximum length. For example:

        - `MAX_WIDTH = 1920 pixels`
        - `MAX_HEIGHT = 1080 pixels`

        This means any image thrown in would be resized but its longest side would be never greater than this.

      - **Aspect Ratio Preservation:** Maintain the original aspect ratio of the image during resizing. This prevents distortion.

    - **Consider a "Quality" Parameter (Optional):** Allow the user to specify a desired quality level (e.g., low, medium, high) or a numeric quality value (0-100) that influences both the resizing and WebP encoding.

4.  **WebP Conversion:**
    - Convert the resized image to WebP format.
    - Use a suitable WebP encoding library (e.g., `golang.org/x/image/webp`).
    - Optimize WebP encoding settings for web use (e.g., quality, compression level).
5.  **Output:**
    - Return the WebP image data. Options:
      - Directly as the response body (Content-Type: `image/webp`).
      - Store the WebP image in an object storage service (e.g., AWS S3, Google Cloud Storage) and return the URL to the stored image. This is better for large images, so you don't have to keep it in memory in the service function.
    - Return metadata about the original image and the resized/converted image (e.g., original dimensions, new dimensions, file size reduction).
6.  **Error Handling:**
    - Gracefully handle invalid image formats, corrupted images, and other potential errors.
    - Return informative error messages to the user.

**Resolution Standards and Logic:**

The "smart" resizing logic is key here. Here's a good approach for determining the final resolution:

1.  **Define Maximum Dimensions:** Set maximum width and height values. These represent the largest dimensions that the output image will ever have.

    ```go
    const (
    	MaxWidth  = 1920 // Typical full-screen width
    	MaxHeight = 1080 // Typical full-screen height
    )
    ```

2.  **Determine the Scaling Factor:** Calculate the scaling factor needed to fit the image within the maximum dimensions while preserving the aspect ratio.

    ```go
    func calculateScale(width, height int) float64 {
    	widthScale := float64(MaxWidth) / float64(width)
    	heightScale := float64(MaxHeight) / float64(height)

    	// Use the smaller scale to ensure the image fits within both dimensions
    	scale := math.Min(widthScale, heightScale)
    	return scale
    }
    ```

3.  **Calculate the New Dimensions:** Apply the scaling factor to the original width and height to get the new dimensions.

    ```go
    func calculateNewDimensions(width, height int, scale float64) (int, int) {
    	newWidth := int(math.Round(float64(width) * scale))
    	newHeight := int(math.Round(float64(height) * scale))
    	return newWidth, newHeight
    }
    ```

**Example Scenario:**

- **Input Image:** 3000x2000 JPEG
- **`MaxWidth`:** 1920
- **`MaxHeight`:** 1080
- **Calculation:**
  - `widthScale = 1920 / 3000 = 0.64`
  - `heightScale = 1080 / 2000 = 0.54`
  - `scale = min(0.64, 0.54) = 0.54`
  - `newWidth = 3000 * 0.54 = 1620`
  - `newHeight = 2000 * 0.54 = 1080`
- **Output:** A 1620x1080 WebP image.

**Benefits of This Approach:**

- **Automatic Optimization:** No need for users to manually specify resizing parameters.
- **Consistent Results:** Images are consistently resized to a standard suitable for web display.
- **Web Performance:** WebP format provides excellent compression and quality.

**Next Steps:**

1.  **Choose a Serverless Platform:** AWS Lambda, Google Cloud Functions, or Azure Functions.
2.  **Set up Your Development Environment:** Install Go and any necessary libraries.
3.  **Implement the Core Logic:**
    - Input handling (file upload or URL).
    - Image decoding (JPEG, PNG).
    - Smart resizing logic (as described above).
    - WebP conversion.
    - Output generation (return image data or store in object storage).
4.  **Test Thoroughly:** Test with a variety of images and sizes.
5.  **Deploy:** Package and deploy your function to your chosen serverless platform.

## Implementation Plan (Test-Driven Development Approach)

### Project Implementation Checklist:

- [x] **Project Setup**

  - [x] Initialize Go module (`go.mod`)
  - [x] Set up directory structure (`cmd/`, `internal/`, `pkg/`, `test/`)
  - [x] Configure linting (`.gitignore`)
  - [x] Create a simple API server stub (`cmd/server/main.go`)

- [x] **Define Core Interfaces**

  - [x] Define interfaces for image processing pipeline (`internal/processor/processor.go`)
  - [x] Create domain models for request/response (`pkg/models/image.go`)
  - [x] Document interfaces (`internal/processor/processor.go`, `pkg/models/image.go`)

- [x] **Image Input Handler**

  - [x] Write tests for URL-based image retrieval (`internal/handler/image_handler_test.go`)
  - [x] Implement image retrieval from URL (`internal/handler/image_handler.go`)
  - [x] Write tests for file upload handling (`internal/handler/file_upload_test.go`)
  - [x] Implement file upload handling (`internal/handler/image_handler.go`)
  - [x] API integration (`internal/api/image_api.go`, `cmd/server/main.go`)

- [x] **Image Processing**

  - [x] Write tests for image format detection (`internal/processor/processor_test.go`)
  - [x] Implement image format detection (`internal/processor/processor.go`)
  - [x] Write tests for image decoding (`internal/processor/processor_test.go`)
  - [x] Implement multi-format image decoder (`internal/processor/processor.go`)

- [x] **Smart Resizing**

  - [x] Write tests for dimension calculation logic (`internal/processor/processor_test.go`)
  - [x] Implement dimension calculation (`internal/processor/processor.go`)
  - [x] Write tests for image resizing (`internal/processor/processor_test.go`)
  - [x] Implement image resizing (`internal/processor/processor.go`)

- [x] **WebP Conversion**

  - [x] Write tests for WebP encoding (`internal/processor/processor_test.go`)
  - [x] Implement WebP conversion (`internal/processor/processor.go`)
  - [x] Write tests for quality optimization (`internal/processor/processor_test.go`)
  - [x] Implement quality settings (`internal/processor/processor.go`)

- [x] **Output Generation**

  - [x] Write tests for response formatting (`internal/api/image_api_test.go`)
  - [x] Implement image response handler (`internal/api/image_api.go`)
  - [x] Write tests for metadata generation (`internal/processor/processor_test.go`)
  - [x] Implement metadata generation (`internal/processor/processor.go`)

- [x] **Integration**

  - [x] Write integration tests (`test/testdata/test_image.html`)
  - [x] Connect all components (`cmd/server/main.go`)
  - [x] Create API endpoints (`internal/api/image_api.go`)
  - [x] Implement comprehensive error handling (`internal/processor/processor.go`, `internal/handler/image_handler.go`)

- [ ] **Performance Testing**

  - [ ] Test with various image sizes and formats
  - [ ] Benchmark memory usage
  - [ ] Identify and optimize bottlenecks
  - [ ] Test concurrency handling

- [x] **Deployment**

  - [x] Write serverless configuration (`scripts/deploy-lambda.sh`)
  - [x] Create deployment scripts (`scripts/run.sh`, `scripts/deploy-lambda.sh`)
  - [ ] Test in staging environment
  - [x] Document deployment process (`README.md`)

- [x] **Documentation**

  - [x] Write API documentation (`README.md`)
  - [x] Create usage examples (`README.md`, `test/testdata/test_image.html`)
  - [x] Document configuration options (`pkg/config/config.go`, `README.md`)
  - [x] Complete project README (`README.md`)

- [ ] **React Frontend**

  - [x] **Setup & Structure**
    - [x] Create React application scaffold (`frontend/`)
    - [x] Set up project dependencies (React, TypeScript, etc.) (`frontend/package.json`)
    - [x] Design component hierarchy and page structure (`frontend/app/page.tsx`, `frontend/components/`)
  - [x] **Core Components**
    - [x] Create drag-and-drop file upload zone (`frontend/components/image-uploader.tsx`)
    - [x] Implement URL input field (`frontend/components/image-uploader.tsx`)
    - [x] Build settings panel (image dimensions, quality) (`frontend/components/settings-form.tsx`)
    - [x] Design before/after image comparison component (`frontend/components/results-grid.tsx`)
    - [x] Add download button for processed images (`frontend/components/results-grid.tsx`)
  - [x] **State & API Integration**
    - [x] Implement state management for user settings (`frontend/app/page.tsx`)
    - [x] Connect to backend API endpoints (`frontend/lib/process-images.ts`)
    - [x] Add upload progress indicators (`frontend/components/image-uploader.tsx`)
    - [x] Handle multiple file uploads (queuing) (`frontend/app/page.tsx`)
    - [x] Implement error handling (`frontend/lib/process-images.ts`)
  - [x] **Styling & UX**
    - [x] Create responsive layout (`frontend/app/page.tsx`, `frontend/components/ui/`)
    - [x] Apply consistent styling (`frontend/app/globals.css`, `frontend/tailwind.config.ts`)
    - [x] Add loading states and animations (`frontend/components/image-uploader.tsx`)
    - [x] Implement accessibility features (`frontend/components/ui/`)
  - [x] **Build & Integration**
    - [x] Configure build process (`frontend/next.config.mjs`)
    - [x] Integrate with Go backend (`frontend/lib/process-images.ts`)
    - [x] Update Go server to serve React static files (`cmd/server/main.go`)

- [x] **Docker Containerization**

  - [x] **Development Environment**
    - [x] Create Dockerfile for combined frontend/backend (`Dockerfile`)
    - [x] Write docker-compose.yml for local development (`docker-compose.yml`)
  - [x] **Production Container**
    - [x] Create multi-stage Dockerfile for combined frontend/backend (`Dockerfile`)
    - [x] Optimize container size and build process (`Dockerfile`)
    - [x] Configure environment variables (`docker-compose.yml`)
  - [x] **Documentation**
    - [x] Document Docker setup and commands (`README.md`)
    - [x] Create quick start guide for local development (`README.md`)

- [x] **Cloud Deployment**
  - [x] **Google Cloud Setup**
    - [x] Create Google Cloud account/project
    - [x] Configure Google Cloud credentials
    - [x] Set up Cloud Storage bucket (optional)
  - [x] **Resource Configuration**
    - [x] Configure memory limits and scaling options
    - [x] Implement request throttling/rate limiting
    - [x] Set up monitoring and logging
  - [x] **Deployment Script**
    - [x] Create script for Google Cloud Run deployment (`scripts/deploy-cloud-run.sh`)
    - [x] Add environment variable configuration (`scripts/deploy-cloud-run.sh`)
    - [x] Set up container registry access (`scripts/deploy-cloud-run.sh`)
  - [x] **Documentation**
    - [x] Document cloud deployment process (`README.md`)
    - [x] Create troubleshooting guide (`README.md`)

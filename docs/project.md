**Project: Smart WebP Image Resizer**

**Description:**

This serverless function takes an image as input (uploaded directly or referenced by URL), automatically resizes it to a sensible standard resolution, converts it to WebP format for optimal web performance, and returns the converted image. The goal is to provide a simple, automated tool that web developers can use to easily optimize images for their websites.

**Functionality:**

1.  **Input:**
    - Accepts images in common formats: JPEG, PNG, GIF (BMP, TIFF if you want to get fancy)
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
    - Image decoding (JPEG, PNG, GIF).
    - Smart resizing logic (as described above).
    - WebP conversion.
    - Output generation (return image data or store in object storage).
4.  **Test Thoroughly:** Test with a variety of images and sizes.
5.  **Deploy:** Package and deploy your function to your chosen serverless platform.

## Implementation Plan (Test-Driven Development Approach)

### Project Implementation Checklist:

- [x] **Project Setup**

  - [x] Initialize Go module
  - [x] Set up directory structure (cmd, internal, pkg, test)
  - [x] Configure linting
  - [x] Create a simple API server stub

- [x] **Define Core Interfaces**

  - [x] Define interfaces for image processing pipeline
  - [x] Create domain models for request/response
  - [x] Document interfaces

- [ ] **Image Input Handler**

  - [ ] Write tests for URL-based image retrieval
  - [ ] Implement image retrieval from URL
  - [ ] Write tests for file upload handling
  - [ ] Implement file upload handling

- [ ] **Image Processing**

  - [ ] Write tests for image format detection
  - [ ] Implement image format detection
  - [ ] Write tests for image decoding
  - [ ] Implement multi-format image decoder

- [ ] **Smart Resizing**

  - [ ] Write tests for dimension calculation logic
  - [ ] Implement dimension calculation
  - [ ] Write tests for image resizing
  - [ ] Implement image resizing

- [ ] **WebP Conversion**

  - [ ] Write tests for WebP encoding
  - [ ] Implement WebP conversion
  - [ ] Write tests for quality optimization
  - [ ] Implement quality settings

- [ ] **Output Generation**

  - [ ] Write tests for response formatting
  - [ ] Implement image response handler
  - [ ] Write tests for metadata generation
  - [ ] Implement metadata generation

- [ ] **Integration**

  - [ ] Write integration tests
  - [ ] Connect all components
  - [ ] Create API endpoints
  - [ ] Implement comprehensive error handling

- [ ] **Performance Testing**

  - [ ] Test with various image sizes and formats
  - [ ] Benchmark memory usage
  - [ ] Identify and optimize bottlenecks
  - [ ] Test concurrency handling

- [ ] **Deployment**

  - [ ] Write serverless configuration
  - [ ] Create deployment scripts
  - [ ] Test in staging environment
  - [ ] Document deployment process

- [ ] **Documentation**
  - [ ] Write API documentation
  - [ ] Create usage examples
  - [ ] Document configuration options
  - [ ] Complete project README

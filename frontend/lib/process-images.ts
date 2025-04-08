import type { ImageFile, ImageSettings, ProcessedImage } from "@/app/page";

// This is a mock function that simulates image processing
// In a real application, this would call a backend API
export async function processImages(
  images: ImageFile[],
  settings: ImageSettings
): Promise<ProcessedImage[]> {
  return Promise.all(
    images.map(async (image) => {
      let response;
      let formData;

      // Process based on whether we have a file or URL
      if (image.file) {
        // Handle file upload
        formData = new FormData();
        formData.append("image", image.file);
        formData.append("max_width", settings.maxWidth.toString());
        formData.append("max_height", settings.maxHeight.toString());
        formData.append("quality", settings.quality.toString());
        formData.append("preserve_ratio", "true");
        formData.append("metadata", "true");

        // Call the backend API
        response = await fetch("/process/upload?metadata=true", {
          method: "POST",
          body: formData,
        });
      } else if (image.url) {
        // Handle URL-based image
        const params = new URLSearchParams({
          url: image.url,
          max_width: settings.maxWidth.toString(),
          max_height: settings.maxHeight.toString(),
          quality: settings.quality.toString(),
          preserve_ratio: "true",
          metadata: "true",
        });

        // Call the backend API
        response = await fetch(`/process/url?${params.toString()}`, {
          method: "GET",
        });
      } else {
        throw new Error("Invalid image: neither file nor URL provided");
      }

      if (!response.ok) {
        throw new Error(`API request failed with status ${response.status}`);
      }

      // Parse metadata from the response
      const metadata = await response.json();

      // Create and store the download form data for uploads
      let downloadFormData;

      // Generate URL for the processed image (without metadata parameter but with download and format parameters)
      let downloadUrl;
      if (image.file) {
        downloadFormData = new FormData();
        downloadFormData.append("image", image.file);
        downloadFormData.append("max_width", settings.maxWidth.toString());
        downloadFormData.append("max_height", settings.maxHeight.toString());
        downloadFormData.append("quality", settings.quality.toString());
        downloadFormData.append("preserve_ratio", "true");
        downloadFormData.append("format", "webp");
        downloadFormData.append("download", "true");

        downloadUrl = "/process/upload?format=webp&download=true";
      } else if (image.url) {
        const params = new URLSearchParams({
          url: image.url,
          max_width: settings.maxWidth.toString(),
          max_height: settings.maxHeight.toString(),
          quality: settings.quality.toString(),
          preserve_ratio: "true",
          format: "webp",
          download: "true",
        });

        downloadUrl = `/process/url?${params.toString()}`;
      }

      return {
        id: image.id,
        preview: image.preview,
        name: image.name,
        metadata: metadata,
        downloadUrl: downloadUrl || "",
        formData: downloadFormData,
      };
    })
  );
}

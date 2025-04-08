"use client";

import type { ProcessedImage } from "@/app/page";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Download, ChevronDown, ChevronUp } from "lucide-react";
import { useState } from "react";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { formatFileSize } from "@/lib/utils";

interface ResultsGridProps {
  processedImages: ProcessedImage[];
}

export function ResultsGrid({ processedImages }: ResultsGridProps) {
  const [expandedItems, setExpandedItems] = useState<Record<string, boolean>>(
    {}
  );

  const toggleExpand = (id: string) => {
    setExpandedItems((prev) => ({
      ...prev,
      [id]: !prev[id],
    }));
  };

  const downloadImage = async (image: ProcessedImage) => {
    try {
      let response;

      // Handle the different download URL types
      if (image.downloadUrl.startsWith("/process/upload")) {
        // For file uploads, we need to make a POST request
        if (!image.formData) {
          console.error("Form data not available for upload");
          return;
        }

        response = await fetch(image.downloadUrl, {
          method: "POST",
          body: image.formData,
        });
      } else {
        // For URL-based images, we can make a GET request
        response = await fetch(image.downloadUrl);
      }

      if (!response.ok) {
        throw new Error(`Failed to download: ${response.status}`);
      }

      // Get the image data as a blob
      let blob = await response.blob();

      // Debug: Check content type
      console.log("Downloaded blob content type:", blob.type);

      // If the blob is not webp, use the original API but explicitly request webp
      if (blob.type !== "image/webp") {
        console.log("Not a WebP image. Attempting to force WebP format...");
        let webpResponse;

        if (image.downloadUrl.startsWith("/process/upload")) {
          // For file uploads with WebP explicitly requested
          if (!image.formData) {
            console.error("Form data not available for upload");
            return;
          }

          // Create a new FormData to avoid modifying the original
          const webpFormData = new FormData();
          for (const [key, value] of image.formData.entries()) {
            webpFormData.append(key, value);
          }

          // Force output to webp
          webpFormData.append("format", "webp");

          webpResponse = await fetch(image.downloadUrl, {
            method: "POST",
            body: webpFormData,
          });
        } else {
          // For URL-based images with WebP explicitly requested
          const url = new URL(image.downloadUrl, window.location.origin);
          url.searchParams.append("format", "webp");

          webpResponse = await fetch(url.toString());
        }

        if (!webpResponse.ok) {
          throw new Error(`Failed to download as WebP: ${webpResponse.status}`);
        }

        const webpBlob = await webpResponse.blob();
        console.log("Forced WebP blob content type:", webpBlob.type);

        // Use the WebP blob instead if successful
        if (webpBlob.type === "image/webp") {
          blob = webpBlob;
        } else {
          console.warn("Failed to convert to WebP, using original format");
        }
      }

      // Create a blob URL
      const blobUrl = URL.createObjectURL(blob);

      // Create and trigger download
      const link = document.createElement("a");
      link.href = blobUrl;

      // Ensure the file has a .webp extension
      let fileName = image.name;
      if (!fileName.toLowerCase().endsWith(".webp")) {
        // Remove existing extension if present
        const nameParts = fileName.split(".");
        if (nameParts.length > 1) {
          nameParts.pop(); // Remove extension
          fileName = nameParts.join(".");
        }
        fileName = `${fileName}.webp`;
      }

      link.download = fileName;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      // Clean up the blob URL after download
      setTimeout(() => URL.revokeObjectURL(blobUrl), 100);
    } catch (error) {
      console.error("Error downloading image:", error);
    }
  };

  const downloadAll = async () => {
    for (const image of processedImages) {
      await downloadImage(image);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Processed Images</h2>
        <Button onClick={downloadAll}>
          <Download className="mr-2 h-4 w-4" />
          Download All
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {processedImages.map((image) => (
          <Card key={image.id} className="overflow-hidden">
            <div className="aspect-video relative">
              <img
                src={image.preview || "/placeholder.svg"}
                alt={image.name}
                className="object-cover"
              />
            </div>

            <CardContent className="p-4">
              <div className="flex justify-between items-start mb-2">
                <h3 className="font-medium truncate" title={image.name}>
                  {image.name}
                </h3>
                <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full">
                  {image.metadata.size_reduction_percent}% smaller
                </span>
              </div>

              <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm">
                <div className="text-muted-foreground">Format:</div>
                <div>
                  {image.metadata.original_format} → {image.metadata.new_format}
                </div>

                <div className="text-muted-foreground">Size:</div>
                <div>
                  {formatFileSize(image.metadata.original_size)} →{" "}
                  {formatFileSize(image.metadata.new_size)}
                </div>

                <div className="text-muted-foreground">Dimensions:</div>
                <div>
                  {image.metadata.original_width}×
                  {image.metadata.original_height} → {image.metadata.new_width}×
                  {image.metadata.new_height}
                </div>
              </div>

              <Collapsible
                open={expandedItems[image.id]}
                onOpenChange={() => toggleExpand(image.id)}
                className="mt-2"
              >
                <CollapsibleTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="w-full flex justify-center items-center p-0 h-6"
                  >
                    {expandedItems[image.id] ? (
                      <ChevronUp className="h-4 w-4" />
                    ) : (
                      <ChevronDown className="h-4 w-4" />
                    )}
                    <span className="ml-1 text-xs">
                      {expandedItems[image.id]
                        ? "Less details"
                        : "More details"}
                    </span>
                  </Button>
                </CollapsibleTrigger>
                <CollapsibleContent className="mt-2">
                  <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
                    <div className="text-muted-foreground">Original Width:</div>
                    <div>{image.metadata.original_width}px</div>

                    <div className="text-muted-foreground">
                      Original Height:
                    </div>
                    <div>{image.metadata.original_height}px</div>

                    <div className="text-muted-foreground">
                      Original Format:
                    </div>
                    <div>{image.metadata.original_format}</div>

                    <div className="text-muted-foreground">Original Size:</div>
                    <div>{formatFileSize(image.metadata.original_size)}</div>

                    <div className="text-muted-foreground">New Width:</div>
                    <div>{image.metadata.new_width}px</div>

                    <div className="text-muted-foreground">New Height:</div>
                    <div>{image.metadata.new_height}px</div>

                    <div className="text-muted-foreground">New Format:</div>
                    <div>{image.metadata.new_format}</div>

                    <div className="text-muted-foreground">New Size:</div>
                    <div>{formatFileSize(image.metadata.new_size)}</div>

                    <div className="text-muted-foreground">Size Reduction:</div>
                    <div>{image.metadata.size_reduction_percent}%</div>
                  </div>
                </CollapsibleContent>
              </Collapsible>
            </CardContent>

            <CardFooter className="p-4 pt-0">
              <Button
                variant="outline"
                className="w-full"
                onClick={() => downloadImage(image)}
              >
                <Download className="mr-2 h-4 w-4" />
                Download
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
}

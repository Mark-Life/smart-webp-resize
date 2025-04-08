"use client";

import { useState } from "react";
import { ImageUploader } from "@/components/image-uploader";
import { SettingsForm } from "@/components/settings-form";
import { ResultsGrid } from "@/components/results-grid";
import { processImages } from "@/lib/process-images";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";

export type ImageFile = {
  id: string;
  file?: File;
  url?: string;
  preview: string;
  name: string;
};

export type ProcessedImage = {
  id: string;
  preview: string;
  name: string;
  metadata: {
    original_width: number;
    original_height: number;
    original_format: string;
    original_size: number;
    new_width: number;
    new_height: number;
    new_format: string;
    new_size: number;
    size_reduction_percent: number;
  };
  downloadUrl: string;
  formData?: FormData;
};

export type ImageSettings = {
  maxWidth: number;
  maxHeight: number;
  quality: number;
};

export default function Home() {
  const [images, setImages] = useState<ImageFile[]>([]);
  const [settings, setSettings] = useState<ImageSettings>({
    maxWidth: 1200,
    maxHeight: 1200,
    quality: 80,
  });
  const [processedImages, setProcessedImages] = useState<ProcessedImage[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);

  const handleAddImages = (newImages: ImageFile[]) => {
    setImages((prev) => {
      const combined = [...prev, ...newImages];
      // Limit to 20 images
      return combined.slice(0, 20);
    });
  };

  const handleRemoveImage = (id: string) => {
    setImages((prev) => prev.filter((img) => img.id !== id));
  };

  const handleTransform = async () => {
    if (images.length === 0) return;

    setIsProcessing(true);
    try {
      const results = await processImages(images, settings);
      setProcessedImages(results);
    } catch (error) {
      console.error("Error processing images:", error);
    } finally {
      setIsProcessing(false);
    }
  };

  const handleClearAll = () => {
    setImages([]);
    setProcessedImages([]);
  };

  return (
    <main className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="space-y-8">
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold">WebP Image Resizer</h1>
          <p className="text-muted-foreground">
            Resize and convert your images to WebP format for better web
            performance
          </p>
        </div>

        <div className="grid gap-8 md:grid-cols-[2fr_1fr]">
          <div className="space-y-6">
            <ImageUploader
              images={images}
              onAddImages={handleAddImages}
              onRemoveImage={handleRemoveImage}
              maxImages={20}
            />

            <div className="flex justify-between">
              <Button
                onClick={handleTransform}
                disabled={images.length === 0 || isProcessing}
                size="lg"
              >
                {isProcessing && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Transform Images
              </Button>

              <Button
                variant="outline"
                onClick={handleClearAll}
                disabled={images.length === 0 && processedImages.length === 0}
              >
                Clear All
              </Button>
            </div>
          </div>

          <div>
            <SettingsForm settings={settings} onSettingsChange={setSettings} />
          </div>
        </div>

        {processedImages.length > 0 && (
          <ResultsGrid processedImages={processedImages} />
        )}
      </div>

      <footer className="mt-16 pt-8 border-t text-center text-sm text-muted-foreground">
        <div className="space-y-2">
          <p className="text-xl">
            <a
              href="https://mark-life.com"
              className="text-primary hover:underline"
            >
              Mark Life | Andrey Markin
            </a>
          </p>
          <p>This Go and Next.js app done in two hours.</p>
          <p>
            Bring your business idea – 3 days to Proof of Concept, 2 weeks to
            MVP
          </p>
        </div>
      </footer>
    </main>
  );
}

"use client";

import type React from "react";

import { useState, useRef } from "react";
import type { ImageFile } from "@/app/page";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { X, Upload, Link, ImageIcon } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { v4 as uuidv4 } from "uuid";

interface ImageUploaderProps {
  images: ImageFile[];
  onAddImages: (images: ImageFile[]) => void;
  onRemoveImage: (id: string) => void;
  maxImages: number;
}

export function ImageUploader({
  images,
  onAddImages,
  onRemoveImage,
  maxImages,
}: ImageUploaderProps) {
  const [dragActive, setDragActive] = useState(false);
  const [urlInput, setUrlInput] = useState("");
  const [urlError, setUrlError] = useState("");
  const fileInputRef = useRef<HTMLInputElement>(null);

  const remainingSlots = maxImages - images.length;

  const handleFiles = (files: FileList | null) => {
    if (!files) return;

    const newImages: ImageFile[] = [];
    const filesToProcess = Math.min(files.length, remainingSlots);

    for (let i = 0; i < filesToProcess; i++) {
      const file = files[i];
      if (file.type.startsWith("image/")) {
        newImages.push({
          id: uuidv4(),
          file,
          preview: URL.createObjectURL(file),
          name: file.name,
        });
      }
    }

    onAddImages(newImages);
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    handleFiles(e.dataTransfer.files);
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFiles(e.target.files);
  };

  const handleUrlSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setUrlError("");

    if (!urlInput.trim()) {
      setUrlError("Please enter a URL");
      return;
    }

    try {
      // Basic URL validation
      new URL(urlInput);

      // In a real app, you'd validate if the URL points to an image
      // For this demo, we'll assume it's valid
      onAddImages([
        {
          id: uuidv4(),
          url: urlInput,
          preview: urlInput,
          name: urlInput.split("/").pop() || "image-from-url",
        },
      ]);

      setUrlInput("");
    } catch (error) {
      setUrlError("Please enter a valid URL");
    }
  };

  return (
    <Card className="border-2">
      <Tabs defaultValue="upload">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="upload">
            <Upload className="mr-2 h-4 w-4" />
            Upload Files
          </TabsTrigger>
          <TabsTrigger value="url">
            <Link className="mr-2 h-4 w-4" />
            Image URL
          </TabsTrigger>
        </TabsList>

        <CardContent className="p-6">
          <TabsContent value="upload" className="mt-0">
            <div
              className={`relative border-2 border-dashed rounded-lg p-6 ${
                dragActive
                  ? "border-primary bg-primary/5"
                  : "border-muted-foreground/25"
              }`}
              onDragEnter={handleDrag}
              onDragOver={handleDrag}
              onDragLeave={handleDrag}
              onDrop={handleDrop}
            >
              <input
                ref={fileInputRef}
                type="file"
                multiple
                accept="image/*"
                className="hidden"
                onChange={handleFileInputChange}
              />

              <div className="flex flex-col items-center justify-center space-y-4 py-4">
                <div className="rounded-full bg-primary/10 p-4">
                  <ImageIcon className="h-8 w-8 text-primary" />
                </div>
                <div className="text-center">
                  <p className="text-lg font-medium">Drag & drop images here</p>
                  <p className="text-sm text-muted-foreground">
                    or click to browse files
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {remainingSlots} of {maxImages} slots remaining
                  </p>
                </div>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => fileInputRef.current?.click()}
                  disabled={remainingSlots === 0}
                >
                  Select Files
                </Button>
              </div>
            </div>
          </TabsContent>

          <TabsContent value="url" className="mt-0">
            <form onSubmit={handleUrlSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="image-url">Image URL</Label>
                <div className="flex gap-2">
                  <Input
                    id="image-url"
                    placeholder="https://example.com/image.jpg"
                    value={urlInput}
                    onChange={(e) => setUrlInput(e.target.value)}
                    className={urlError ? "border-red-500" : ""}
                    disabled={remainingSlots === 0}
                  />
                  <Button type="submit" disabled={remainingSlots === 0}>
                    Add
                  </Button>
                </div>
                {urlError && <p className="text-sm text-red-500">{urlError}</p>}
              </div>
            </form>
          </TabsContent>
        </CardContent>
      </Tabs>

      {images.length > 0 && (
        <div className="p-6 pt-0">
          <h3 className="font-medium mb-3">
            Selected Images ({images.length})
          </h3>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
            {images.map((image) => (
              <div key={image.id} className="relative group">
                <div className="aspect-square rounded-md overflow-hidden border bg-muted">
                  <img
                    src={image.preview || "/placeholder.svg"}
                    alt={image.name}
                    width={200}
                    height={200}
                    className="object-cover w-full h-full"
                  />
                </div>
                <Button
                  variant="destructive"
                  size="icon"
                  className="absolute -top-2 -right-2 h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
                  onClick={() => onRemoveImage(image.id)}
                >
                  <X className="h-3 w-3" />
                </Button>
                <p className="text-xs truncate mt-1">{image.name}</p>
              </div>
            ))}
          </div>
        </div>
      )}
    </Card>
  );
}

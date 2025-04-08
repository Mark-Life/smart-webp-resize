#!/bin/bash

# Smart WebP Resizer - Local Run Script

# Set environment variables
export PORT=8080
export MAX_WIDTH=1920
export MAX_HEIGHT=1080
export DEFAULT_QUALITY=85

# Build the application
echo "Building the application..."
go build -o bin/smart-webp-resizer ./cmd/server

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build successful!"

# Run the application
echo "Starting server on port $PORT..."
bin/smart-webp-resizer 
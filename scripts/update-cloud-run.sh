#!/bin/bash
set -e

# Configuration
PROJECT_ID="smart-webp-resize"
REGION="europe-west1"
SERVICE_NAME="webp-resizer"
IMAGE_REPO="europe-west1-docker.pkg.dev/smart-webp-resize/webp-resizer-repo/webp-resizer"
TAG="latest"

# Build new image
echo "Building new container image..."
gcloud builds submit --tag ${IMAGE_REPO}:${TAG}

# Deploy update
echo "Deploying updated service..."
gcloud run deploy ${SERVICE_NAME} \
  --image=${IMAGE_REPO}:${TAG} \
  --region=${REGION} \
  --clear-cloudsql-instances \
  --no-clear-env-vars \
  --no-clear-labels

echo "Update complete!"
echo "Service URL: $(gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format='value(status.url)')"
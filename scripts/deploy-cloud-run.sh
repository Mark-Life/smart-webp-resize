#!/bin/bash
set -e

# Configuration variables
PROJECT_ID="YOUR_GCP_PROJECT_ID"  # Replace with your GCP project ID
REGION="us-central1"              # Replace with your preferred region
SERVICE_NAME="webp-resizer"       # Name of the Cloud Run service
IMAGE_NAME="webp-resizer"         # Name of the container image
MAX_INSTANCES=10                  # Maximum number of instances to scale to
MEMORY="512Mi"                    # Memory allocation per instance
CPU="1"                           # CPU allocation (1 = 1 vCPU)
CONCURRENCY=80                    # Maximum concurrent requests per instance
TIMEOUT="300s"                    # Maximum request timeout
REQUEST_LIMIT="60/minute"         # Rate limiting for cost management

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}Error: gcloud CLI is not installed.${NC}"
    echo "Please install the Google Cloud SDK: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Ensure user is logged in
echo -e "${YELLOW}Checking authentication...${NC}"
gcloud auth print-identity-token &> /dev/null || gcloud auth login

# Set the project
echo -e "${YELLOW}Setting project to ${PROJECT_ID}...${NC}"
gcloud config set project ${PROJECT_ID}

# Enable required services if not already enabled
echo -e "${YELLOW}Enabling required services...${NC}"
gcloud services enable cloudbuild.googleapis.com run.googleapis.com artifactregistry.googleapis.com

# Create Container Registry repository if it doesn't exist
echo -e "${YELLOW}Creating/checking Container Registry...${NC}"
gcloud artifacts repositories describe ${IMAGE_NAME}-repo --location=${REGION} &> /dev/null || \
gcloud artifacts repositories create ${IMAGE_NAME}-repo --repository-format=docker --location=${REGION}

# Build the container image using Cloud Build
echo -e "${YELLOW}Building container image...${NC}"
gcloud builds submit --tag ${REGION}-docker.pkg.dev/${PROJECT_ID}/${IMAGE_NAME}-repo/${IMAGE_NAME}:latest

# Deploy to Cloud Run with resource constraints and rate limiting
echo -e "${YELLOW}Deploying to Cloud Run...${NC}"
gcloud run deploy ${SERVICE_NAME} \
  --image ${REGION}-docker.pkg.dev/${PROJECT_ID}/${IMAGE_NAME}-repo/${IMAGE_NAME}:latest \
  --platform managed \
  --region ${REGION} \
  --memory ${MEMORY} \
  --cpu ${CPU} \
  --max-instances ${MAX_INSTANCES} \
  --concurrency ${CONCURRENCY} \
  --timeout ${TIMEOUT} \
  --set-env-vars MAX_WIDTH=1920,MAX_HEIGHT=1080,DEFAULT_QUALITY=85 \
  --no-allow-unauthenticated

# Configure rate limiting
echo -e "${YELLOW}Configuring rate limiting (${REQUEST_LIMIT})...${NC}"
gcloud run services update ${SERVICE_NAME} \
  --region ${REGION} \
  --platform managed \
  --update-labels run.googleapis.com/request-rate-limit=${REQUEST_LIMIT}

# Create a service account for invoking the service
echo -e "${YELLOW}Creating service account for invocation...${NC}"
SA_NAME="${SERVICE_NAME}-invoker"
SA_EMAIL="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

gcloud iam service-accounts create ${SA_NAME} \
  --display-name "${SERVICE_NAME} Invoker" \
  || echo -e "${YELLOW}Service account already exists${NC}"

# Grant the service account permission to invoke the service
gcloud run services add-iam-policy-binding ${SERVICE_NAME} \
  --region=${REGION} \
  --member="serviceAccount:${SA_EMAIL}" \
  --role="roles/run.invoker"

# Create a service account key (optional - you might want to use Workload Identity instead)
# gcloud iam service-accounts keys create ${SA_NAME}-key.json --iam-account=${SA_EMAIL}

# Display service URL
SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format='value(status.url)')

echo -e "${GREEN}Deployment complete!${NC}"
echo -e "Service URL: ${SERVICE_URL}"
echo -e "${YELLOW}Note: This service requires authentication to access.${NC}"
echo -e "To generate an authentication token, run: gcloud auth print-identity-token"
echo -e "To test the service, run: curl -H \"Authorization: Bearer \$(gcloud auth print-identity-token)\" ${SERVICE_URL}/health"
echo
echo -e "${YELLOW}Cost Management:${NC}"
echo -e "- Max instances: ${MAX_INSTANCES}"
echo -e "- Memory per instance: ${MEMORY}"
echo -e "- Request rate limit: ${REQUEST_LIMIT}"
echo -e "- Authentication required: Yes" 
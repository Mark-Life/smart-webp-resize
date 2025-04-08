#!/bin/bash

# Smart WebP Resizer - AWS Lambda Deployment Script

# Configuration
FUNCTION_NAME="smart-webp-resizer"
REGION="us-east-1"  # Change to your region
ROLE_ARN="arn:aws:iam::YOUR_ACCOUNT_ID:role/lambda-execution-role"  # Change to your role ARN
MEMORY_SIZE=512
TIMEOUT=30
HANDLER="smart-webp-resizer"

# Create bin directory if it doesn't exist
mkdir -p bin

# Build for AWS Lambda (Linux)
echo "Building binary for AWS Lambda..."
GOOS=linux GOARCH=amd64 go build -o bin/$HANDLER ./cmd/server

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Zip the binary
echo "Creating deployment package..."
cd bin
zip -j $FUNCTION_NAME.zip $HANDLER
cd ..

# Check if zip was successful
if [ $? -ne 0 ]; then
    echo "Creating zip file failed!"
    exit 1
fi

# Check if function exists
echo "Checking if function exists..."
aws lambda get-function --function-name $FUNCTION_NAME --region $REGION > /dev/null 2>&1

if [ $? -eq 0 ]; then
    # Update existing function
    echo "Updating existing function..."
    aws lambda update-function-code \
        --function-name $FUNCTION_NAME \
        --zip-file fileb://bin/$FUNCTION_NAME.zip \
        --region $REGION

    # Update configuration
    aws lambda update-function-configuration \
        --function-name $FUNCTION_NAME \
        --timeout $TIMEOUT \
        --memory-size $MEMORY_SIZE \
        --region $REGION \
        --environment "Variables={MAX_WIDTH=1920,MAX_HEIGHT=1080,DEFAULT_QUALITY=85}"
else
    # Create new function
    echo "Creating new function..."
    aws lambda create-function \
        --function-name $FUNCTION_NAME \
        --runtime go1.x \
        --role $ROLE_ARN \
        --handler $HANDLER \
        --zip-file fileb://bin/$FUNCTION_NAME.zip \
        --timeout $TIMEOUT \
        --memory-size $MEMORY_SIZE \
        --region $REGION \
        --environment "Variables={MAX_WIDTH=1920,MAX_HEIGHT=1080,DEFAULT_QUALITY=85}"
fi

if [ $? -eq 0 ]; then
    echo "Deployment completed successfully!"
    echo "To test function, use: aws lambda invoke --function-name $FUNCTION_NAME --payload '{}' output.txt --region $REGION"
else
    echo "Deployment failed!"
    exit 1
fi 
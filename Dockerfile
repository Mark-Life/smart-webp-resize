# Stage 1: Build the React frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Install pnpm
RUN npm install -g pnpm

# Copy frontend package.json and pnpm-lock.yaml
COPY frontend/package.json frontend/pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install

# Copy the rest of the frontend code
COPY frontend ./

# Build the frontend
RUN pnpm build

# Stage 2: Build the Go backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install C compiler and development tools for CGO
RUN apk add --no-cache gcc musl-dev

# Enable CGO
ENV CGO_ENABLED=1

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the Go code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY test/ ./test/

# Build the Go application
RUN go build -o server ./cmd/server

# Stage 3: Final image
FROM alpine:latest

WORKDIR /app

# Install dependencies for running Go binary
RUN apk --no-cache add ca-certificates

# Copy the Go binary from the backend builder
COPY --from=backend-builder /app/server .

# Copy test data directory for html test page
COPY --from=backend-builder /app/test ./test

# Copy the built React app from the frontend builder
COPY --from=frontend-builder /app/frontend/out ./frontend/out

# Create a directory for the frontend static files
RUN mkdir -p /app/static

# Move the frontend build to the static directory
RUN cp -r /app/frontend/out/* /app/static/

# Expose the port the app will run on
EXPOSE 8080

# Define environment variables
ENV PORT=8080

# Command to run the application
CMD ["./server"] 
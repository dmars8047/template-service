# Use the official Go image as the base image
FROM golang:1.20.5-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Run unit tests
RUN go test ./...

# Build the application
RUN go build -o template-service

# Start a new stage with a minimal image
FROM scratch

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/template-service /

# Expose the desired port
EXPOSE 8080

# Set the entry point for the container
ENTRYPOINT ["/template-service"]
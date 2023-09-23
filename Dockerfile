# Use the golden golang alpine image as the base image
FROM dmars8047/golang-alpine-golden:latest-develop AS builder

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
RUN go build -o /app/template-service ./cmd/template-service

# Start a new stage with a minimal image
FROM scratch

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/template-service /

# Expose the desired port
EXPOSE 8080

# Set the entry point for the container
ENTRYPOINT ["/template-service"]

# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Build the Go application
RUN go build -o mediaserver

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/mediaserver .

# Create the media folder inside the container
RUN mkdir -p ./movies

# Expose the HTTP port
EXPOSE 8080

# Run the server
CMD ["./mediaserver"]

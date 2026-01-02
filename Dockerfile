# Use a Go base image
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the WebAssembly binary
RUN mkdir -p web
RUN GOOS=js GOARCH=wasm go build -o web/app.wasm

# Build the server binary
RUN go build -o server

# Expose the port
EXPOSE 8000

# Run the server
CMD ["./server"]

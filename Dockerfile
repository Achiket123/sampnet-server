# Start with a Golang base image
FROM golang:1.23.2-alpine3.19 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/server

# Start a new, clean image with minimal dependencies to run the app
FROM gcr.io/distroless/base-debian10

# Set working directory in container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]

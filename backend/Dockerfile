FROM golang:1.21-alpine AS base

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Development stage
FROM base AS development
RUN go install github.com/cosmtrek/air@latest
EXPOSE 9090
CMD ["go", "run", "main.go"]

# Production stage
FROM base AS production

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 9090

# Run the application
CMD ["./main"]

# Use a minimal alpine image
FROM alpine:latest AS final

WORKDIR /app

# Copy the binary from the build stage
COPY --from=production /app/main .

# Expose port
EXPOSE 9090

# Run the application
CMD ["./main"]
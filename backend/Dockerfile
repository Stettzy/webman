FROM golang:1.21-alpine AS base

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Development stage
FROM base AS development
RUN go install github.com/cosmtrek/air@latest
EXPOSE 8080
CMD ["air"]

# Production stage
FROM base AS production

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
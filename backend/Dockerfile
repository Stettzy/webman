FROM golang:1.22-alpine AS base

WORKDIR /app

# Development stage
FROM base AS development

EXPOSE 9090
CMD ["go", "run", "main.go"]

# Production stage
FROM base AS production

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o main .

# Final stage
FROM alpine:latest AS final

WORKDIR /app

# Copy the binary from the build stage
COPY --from=production /app/main .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 9090

# Run the application
CMD ["./main"]
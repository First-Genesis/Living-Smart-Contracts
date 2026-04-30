# Multi-stage build for Living Smart Contracts
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files
COPY go.mod ./

# Download dependencies (if any)
RUN go mod download

# Copy source code (includes Swagger docs)
COPY . .

# Build the application with Swagger documentation
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o living-contracts ./cmd/server

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/living-contracts .

# Create directories for data and logs
RUN mkdir -p /app/data /app/logs && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment variables
ENV PORT=8080 \
    LOG_LEVEL=info \
    ACTOR_SYSTEM_NAME=living-contracts \
    HTTP_PORT=8080 \
    METRICS_PORT=9090

# Run the application
CMD ["./living-contracts"]

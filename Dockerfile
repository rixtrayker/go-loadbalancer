# Build stage
FROM golang:1.24-alpine AS builder

# Add build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod and sum files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with version info
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG COMMIT=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${COMMIT}" \
    -o loadbalancer .

# Final stage
FROM alpine:latest

# Add runtime dependencies and security updates
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && update-ca-certificates \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app

# Copy only necessary files from builder
COPY --from=builder --chown=appuser:appgroup /build/loadbalancer .
COPY --from=builder --chown=appuser:appgroup /build/config/config.toml ./config/

# Create a non-root user and switch to it
USER appuser

# Set environment variables
ENV GIN_MODE=release

# Expose the port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["./loadbalancer"] 
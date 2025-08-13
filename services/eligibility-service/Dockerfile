# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy source code first (simplified approach)
COPY . .

# Initialize and download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S nphies && \
    adduser -u 1001 -S nphies -G nphies

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create directories for logs and config
RUN mkdir -p /var/log/nphies /etc/nphies/config

# Change ownership to nphies user
RUN chown -R nphies:nphies /root /var/log/nphies /etc/nphies

# Switch to non-root user
USER nphies

# Expose port
EXPOSE 8090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8090/health || exit 1

# Run the application
CMD ["./main"]
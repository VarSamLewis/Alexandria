# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies (required for sqlite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o Alexandria .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

# Create app directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/Alexandria .

# Create directory for database
RUN mkdir -p /root/work/DB/Alexandria

# Set the database path as an environment variable
ENV DB_PATH=/root/work/DB/Alexandria/tickets.db

# Entrypoint
ENTRYPOINT ["./Alexandria"]

# Default command (shows help)
CMD ["--help"]

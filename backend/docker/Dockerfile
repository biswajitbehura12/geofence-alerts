# Backend Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o geofence ./cmd/main.go

# Final stage
FROM alpine:3.18

WORKDIR /root/

# Install postgresql client
RUN apk add --no-cache postgresql-client

# Copy binary from builder
COPY --from=builder /app/geofence .

# Expose port
EXPOSE 8080

# Run application
CMD ["./geofence"]

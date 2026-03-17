# -------- Build Stage --------
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git (required for private/public modules)
RUN apk add --no-cache git

# Copy go module files
COPY go.mod ./
# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# -------- Final Stage --------
FROM alpine:3.18

WORKDIR /app

# Install required packages (optional: postgres client)
RUN apk add --no-cache ca-certificates postgresql-client

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port
EXPOSE 8080

# Run application
CMD ["./app"]
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app .

# Run stage
FROM alpine:3.18
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .

CMD ["./app"]
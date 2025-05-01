
# syntax=docker/dockerfile:1.4

# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install CA certificates (for https requests)
RUN apk add --no-cache ca-certificates

# Copy Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app statically
RUN CGO_ENABLED=0 go build -o main .

# Final stage
FROM alpine:latest

# Copy CA certs and app binary from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /usr/local/bin/main

EXPOSE 8080

ENTRYPOINT ["main"]

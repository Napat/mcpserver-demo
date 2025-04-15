# Stage 1: Build the application
FROM golang:1.23.2-alpine AS builder

# Install git and certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create and set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build mcpserver application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/mcpserver ./cmd/mcpserver

# Stage 2: Create the final image
FROM alpine:3.21

# Add necessary certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /go/bin/mcpserver /app/mcpserver

# Set the entry point
ENTRYPOINT ["/app/mcpserver"] 
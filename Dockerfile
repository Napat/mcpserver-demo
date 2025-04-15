FROM golang:1.23-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git make

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create output directory
RUN mkdir -p bin

# Build the application
RUN make build

# Create final lightweight image
FROM alpine:latest

# Install dependencies
RUN apk --no-cache add ca-certificates tzdata make

# Set the timezone
ENV TZ=Asia/Bangkok

# Set the working directory
WORKDIR /app

# Copy the binary and necessary files
COPY --from=builder /app/bin/mcpserver .
COPY --from=builder /app/public ./public
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/Makefile .
COPY --from=builder /app/cmd/migrate ./cmd/migrate

# Create necessary directories
RUN mkdir -p logs configs/temp

# Expose the port
EXPOSE 8080

# Run the application with config
CMD ["sh", "-c", "cp configs/api/.env configs/temp/.env && ./mcpserver"] 
# Build Stage
FROM golang:1.23-alpine AS builder

# Install necessary packages
RUN apk update && apk add --no-cache git make

# Install Go tools
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Set work directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Generate swagger documentation
RUN swag init -g internal/api/docs.go

# Generate sqlc code
RUN sqlc generate

# Build the application
RUN go build -o /bin/app cmd/api/main.go

# Final Stage
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set work directory
WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /bin/app .

# Expose the server port
EXPOSE 8080

# Command to run the executable
CMD ["./app"]

FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates \
    gcc \
    musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go install --tags musl ./...


# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /go/bin/* /usr/bin/
COPY auth.pem auth.pem

# Change to non-root user
USER appuser

# Command to run the executable
ENTRYPOINT ["teleport-discord-bot"]

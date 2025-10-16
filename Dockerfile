# ============================================
# STEP 1: Build stage
# ============================================
FROM golang:1.25.1 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# ============================================
# STEP 2: Runtime stage
# ============================================
FROM alpine:3.20

WORKDIR /app

# Install SSL certificates (required for Postgres)
RUN apk add --no-cache ca-certificates

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

# Copy .env file if you’re using it
COPY .env .env

# Expose your app port (default Fiber port)
EXPOSE 4000

# Run the app
CMD ["./main"]

# -------- STAGE 1: Build --------
FROM golang:1.23.9 as builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

COPY vendor ./vendor/

# Copy the entire source
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o main ./cmd
# -------- STAGE 2: Minimal Image --------
FROM alpine

# Copy the binary from builder
COPY --from=builder /app/main /main

# Copy config if needed
COPY --from=builder /app/config /config
ENV CONFIG_PATH=/config

COPY --from=builder /app/db/migrations /migrations


# Run the app
ENTRYPOINT ["/main"]
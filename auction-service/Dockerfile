# Stage 1: Build
FROM golang:1.23.9 AS builder

WORKDIR /app

# Copy go files and vendor directory
#COPY go.mod go.sum ./
#RUN go mod download

COPY . ./
#RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o main ./cmd

# Final image
FROM alpine:latest

COPY --from=builder /app/main /main
COPY --from=builder /app/config /config
ENV CONFIG_PATH=/config

COPY --from=builder /app/db/migrations /migrations

ENTRYPOINT ["/main"]
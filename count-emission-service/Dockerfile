# --- Stage 1: Build binary ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a statically linked binary optimized for production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/gRPC/main.go

# --- Stage 2: Final lightweight image ---
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for secure outbound HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy binary from builder stage
COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]
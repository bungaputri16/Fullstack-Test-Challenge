# --- Stage 1: Build ---
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# copy go mod/sum lalu download deps dulu (biar cache efektif)
COPY go.mod go.sum ./
RUN go mod download

# copy semua source code
COPY . .

# build binary
RUN go build -o order-service ./cmd/main.go

# --- Stage 2: Run ---
FROM alpine:3.19

WORKDIR /app

# copy binary dari builder
COPY --from=builder /app/order-service .

# expose port (default 4000 dari .env)
EXPOSE 4000

# jalankan service
CMD ["./order-service"]

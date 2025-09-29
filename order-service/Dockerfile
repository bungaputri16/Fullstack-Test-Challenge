FROM golang:1.25-alpine

WORKDIR /app

# Copy go.mod & go.sum dulu untuk install dependencies
COPY go.mod go.sum ./

# Download semua dependencies (termasuk github.com/google/uuid)
RUN go mod tidy
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary
RUN go build -o order-service ./cmd/main.go

EXPOSE 4000

CMD ["./order-service"]

# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency definition files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api/main.go

# Production stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy compiled binary and configuration example
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.env.example ./config.env

EXPOSE 8080

CMD ["./api"]

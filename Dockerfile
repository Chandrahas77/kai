# Base builder stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Install Goose migration tool
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Ensure the migrations directory is copied
COPY migrations /app/migrations

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations  

# Run migrations before starting the app
CMD ["./main"]
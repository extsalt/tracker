# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o tracker cmd/tracker/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/tracker .
EXPOSE 8000
CMD ["./tracker"]

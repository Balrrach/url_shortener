FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the source code (specifically the app directory)
COPY ./app ./app
COPY ./templates ./templates

# Build the application
WORKDIR /app/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/urlshortener .

# Use a smaller image for the final application
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/urlshortener .
COPY --from=builder /app/templates ./templates

# Set environment variables
ENV GIN_MODE=release

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./urlshortener"]

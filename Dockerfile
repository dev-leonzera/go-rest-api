# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Initialize/Tidy the module and download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN go build -o main .

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./main"]

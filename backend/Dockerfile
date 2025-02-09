# Use a lightweight base image for Go
FROM golang:1.21 as builder
# Set the working directory
WORKDIR /app
# Copy the Go module files
COPY go.mod go.sum ./
RUN go mod download
# Copy the entire backend source code
COPY . .
# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller image for deployment
FROM alpine:latest
# Install dependencies
RUN apk --no-cache add ca-certificates
# Set working directory
WORKDIR /app
# Copy compiled binary from builder
COPY --from=builder /app/main .
# Copy environment file
COPY .env .
# Expose application port
EXPOSE 8080
# Run the binary
CMD ["./main"]

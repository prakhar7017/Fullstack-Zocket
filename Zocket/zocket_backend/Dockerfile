# Use the correct Go version
FROM golang:1.24 AS builder

# Set the working directory
WORKDIR /app

# Copy Go modules files and download dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Use a minimal runtime image
FROM alpine:latest

# Install libc compatibility for Alpine
RUN apk --no-cache add ca-certificates

# Set working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder
COPY --from=builder /app/main .

# Copy the .env file into the container
# COPY .env . 

# Ensure the binary is executable
RUN chmod +x ./main

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]

# Start from the official Golang image to create a build artifact.
FROM golang:1.21 as builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to download dependencies.
COPY go.mod ./
COPY go.sum* ./

# Download dependencies.
RUN go mod download

# Copy the source code.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook_handler .

# Use a minimal alpine image for the final stage.
FROM alpine:latest  

WORKDIR /root/

# Copy the pre-built binary from the builder stage.
COPY --from=builder /app/webhook_handler .

# Expose port 8080.
EXPOSE 8080

# Command to run the executable.
CMD ["./webhook_handler"]

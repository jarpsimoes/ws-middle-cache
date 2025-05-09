# Use the official Golang image as the base image
FROM golang:1.24.3

# Install glibc
RUN apt-get update && apt-get install -y --no-install-recommends \
    libc6 \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install Gin framework
RUN go get -u github.com/gin-gonic/gin

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["go", "run", "cmd/server/main.go"]
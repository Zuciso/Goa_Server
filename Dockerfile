# Start with an official Go image
FROM golang:1.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app (without cross-compilation)
RUN go build -o server .

# Start a new stage from scratch
FROM debian:bullseye-slim

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Ensure the server binary is executable
RUN chmod +x /root/server

# Copy the wait-for-it.sh script into the container
COPY wait-for-it.sh /root/

# Ensure that the script is executable
RUN chmod +x /root/wait-for-it.sh

# Expose port 8080
EXPOSE 8080

# Run the Go app, waiting for MySQL to be ready
CMD ["/root/wait-for-it.sh", "mysql_container:3306", "--", "./server"]


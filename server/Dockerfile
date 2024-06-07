# Use the official Golang image to create a build artifact.
# This is the first stage, called the 'builder'.
FROM golang:1.22-alpine AS builder

# Install necessary packages.
RUN apk add --no-cache git

# Set the current working directory inside the container.
WORKDIR /app

COPY go.mod go.sum client_secret.json youtube-credentials.json ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source from the current directory to the working directory inside the container.
COPY . .

# Build the Go app.
RUN go build -o server .

# Ensure the server binary has executable permissions in the builder stage.
RUN chmod +x server

# Use a minimal image as the base for the final artifact.
# This results in a smaller image size.
FROM alpine:latest

# Install necessary packages.
RUN apk --no-cache add ca-certificates

# Set the current working directory inside the container.
WORKDIR /app

# Copy the binary and other necessary files from the builder stage.
COPY --from=builder /app/server /app/
COPY --from=builder /app/client_secret.json /app/

# Debugging: List the files and their permissions before setting permissions
RUN ls -l /app

# Ensure the server binary has executable permissions in the final stage.
RUN chmod +x /app/server

# Debugging: List the files and their permissions after setting permissions
RUN ls -l /app

# Expose port 3333 to the outside world.
EXPOSE 3333

# Command to run the executable.
CMD ["/app/server"]

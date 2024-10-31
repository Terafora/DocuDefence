# Dockerfile (inside backend folder)
FROM golang:1.20

# Set the working directory within the container
WORKDIR /app

# Copy go.mod and go.sum for dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Expose the port the backend listens on
EXPOSE 8000

# Build the Go application
RUN go build -o main .

# Run the application
CMD ["./main"]

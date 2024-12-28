# Uses an official Go base image
FROM golang:1.23.1

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Create and define the working directory
WORKDIR /app

# Copies the project files to the container
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Exposes the port your application listens on
EXPOSE 3003

# Default command
CMD ["go", "run", "main.go"]

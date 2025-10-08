FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy the rest of the application source code
COPY . .

RUN go mod download

# Build the Go binary
RUN go build .

# Command to run the application
CMD ["go", "run", "."]
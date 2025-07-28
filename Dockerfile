FROM golang:1.24.4

# Install poppler-utils for PDF text extraction
RUN apt-get update && apt-get install -y poppler-utils

# Set working directory to /app
WORKDIR /app

# Copy go.mod (and go.sum if it exists)
COPY go.mod ./
COPY go.sum* ./

# Download dependencies (if any)
RUN go mod download || true

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Define the command to run the application
CMD ["./main"]
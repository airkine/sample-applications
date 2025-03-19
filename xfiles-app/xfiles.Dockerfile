# Use Debian as the base image for building Go
FROM debian:bullseye-slim AS builder

# Install dependencies
RUN apt-get update && apt-get install -y wget tar gcc libc-dev musl-dev

# Set Go version
ENV GO_VERSION=1.23.6

# Download and install Go manually
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set up Go paths
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"

# Create app directory
WORKDIR /app

# Copy Go module files and install dependencies
COPY go.mod go.sum ./
RUN /usr/local/go/bin/go mod tidy

# Copy the entire application source code
COPY . .

# Build a fully static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . && chmod +x main

# Use distroless for minimal runtime
FROM gcr.io/distroless/static:latest

WORKDIR /app

# Copy the built binary, templates, and static files from the builder
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static  
# ✅ Ensure static files are copied

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]

# Use a generic Linux base image
FROM debian:bullseye AS base

# Install required tools
RUN apt-get update && apt-get install -y curl tar

# Install Go manually (version 1.23.2)
RUN curl -sSL https://go.dev/dl/go1.23.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH="/usr/local/go/bin:$PATH"

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the application source code
COPY . .

# Build the application
RUN go build -o azure-footprint azure_footprint.go

# Use a minimal image for runtime
FROM debian:bullseye-slim AS runtime

# Set the working directory
WORKDIR /app

# Copy the built application
COPY --from=base /app/azure-footprint /app/

# Run the application
CMD ["./azure-footprint"]

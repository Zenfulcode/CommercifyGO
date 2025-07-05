FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy the entire codebase first
COPY . .

# Delete the go.work files that reference external modules
RUN rm -f go.work go.work.sum

# Build all three applications
RUN go mod download
RUN go build -o commercify cmd/api/main.go
RUN go build -o commercify-seed cmd/seed/main.go

# Create a minimal final image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata bash

# Copy the binaries from the builder stage
COPY --from=builder /app/commercify /app/commercify
COPY --from=builder /app/commercify-seed /app/commercify-seed
COPY --from=builder /app/templates /app/templates

# Copy .env file if it exists (will be overridden by env_file in docker-compose)
# COPY --from=builder /app/.env /app/

# Set executable permissions for all binaries
RUN chmod +x /app/commercify /app/commercify-seed

# Expose the port
EXPOSE 6091

# Run the API by default
CMD ["/app/commercify"]
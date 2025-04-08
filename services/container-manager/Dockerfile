# ---- Build Stage ----
FROM docker.io/golang:1.24 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the application (static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# ---- Run Stage ----
FROM registry.access.redhat.com/ubi8/ubi-minimal

# Set up a non-root user for security
RUN microdnf install shadow-utils
USER root

# Set working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port (change if necessary)
EXPOSE 8080

# Set the entrypoint to the compiled binary
ENTRYPOINT ["/app/main"]


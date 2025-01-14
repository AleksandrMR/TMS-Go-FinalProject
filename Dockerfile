# Stage 1: Build stage
FROM golang:1.21-bookworm as builder

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.* ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o hashService ./cmd/hashService

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/hashService .

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Set the entrypoint command
ENTRYPOINT ["/app/hashService"]

# Command to run the executable
CMD ["./hashService"]

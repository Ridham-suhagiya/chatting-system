# Build stage
FROM --platform=linux/amd64 golang:1.22-alpine AS build

# Set environment variables to ensure cross-compilation for the correct architecture
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOAMD64=v3

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache make curl


# Copy the source code
COPY . .

# Build the Go application
RUN go build -o server .

# Run stage
FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

# Copy the built Go binary
COPY --from=build /app/server .

# Copy Makefile and migrations folder to final image
COPY --from=build /app/Makefile Makefile
COPY --from=build /app/migrations/ migrations/

# Ensure the binary has execution permission
RUN chmod +x ./server

# Expose the application port
EXPOSE 8080

# Run the application
CMD [ "./server"]

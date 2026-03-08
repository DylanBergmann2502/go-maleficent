# deploy/production/compose/go/Dockerfile
FROM golang:1.26-alpine AS go-build-stage

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Create bin directory and build the application with optimizations
RUN mkdir -p bin && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bin/server ./cmd/server/main.go

# Production stage
FROM alpine:latest

# Install minimal runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    wget \
    bash \
    && rm -rf /var/cache/apk/*

# Install wait-for-it manually
RUN wget -O /usr/local/bin/wait-for-it https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
    && chmod +x /usr/local/bin/wait-for-it

# Create non-root user
RUN addgroup -g 1000 app \
    && adduser -u 1000 -G app -s /bin/sh -D app

WORKDIR /app

# Create bin directory and copy built application from build stage
RUN mkdir -p bin
COPY --from=go-build-stage /app/bin/server ./bin/server
RUN chmod +x ./bin/server

# Copy production scripts
COPY ./deploy/production/compose/go/entrypoint /entrypoint
COPY ./deploy/production/compose/go/start /start

# Remove Windows line endings and make scripts executable
RUN sed -i 's/\r$//g' /entrypoint \
    && sed -i 's/\r$//g' /start \
    && chmod +x /entrypoint /start \
    && chown -R app:app /app

# Switch to non-root user
USER app

EXPOSE 8000

ENTRYPOINT ["/entrypoint"]

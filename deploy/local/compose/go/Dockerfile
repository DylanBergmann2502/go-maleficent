# deploy/local/compose/go/Dockerfile
FROM golang:1.25-alpine AS go-build-stage

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Create bin directory and build the application
RUN mkdir -p bin && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server/main.go

# Production stage
FROM alpine:latest AS go-run-stage

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    bash \
    nano \
    vim \
    sudo \
    postgresql-client \
    && rm -rf /var/cache/apk/*

# Install wait-for-it manually
RUN wget -O /usr/local/bin/wait-for-it https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
    && chmod +x /usr/local/bin/wait-for-it

# Create development user and add to sudoers
RUN addgroup -g 1000 dev-user \
    && adduser -u 1000 -G dev-user -s /bin/bash -D dev-user \
    && echo "dev-user ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/dev-user \
    && chmod 0440 /etc/sudoers.d/dev-user

WORKDIR /app

# Create bin directory and copy built application from build stage
RUN mkdir -p bin
COPY --from=go-build-stage /app/bin/server ./bin/server
RUN chmod +x ./bin/server

# Copy scripts
COPY ./deploy/local/compose/go/entrypoint /entrypoint
COPY ./deploy/local/compose/go/start /start

# Remove Windows line endings and make scripts executable
RUN sed -i 's/\r$//g' /entrypoint \
    && sed -i 's/\r$//g' /start \
    && chmod +x /entrypoint /start

EXPOSE 8000

ENTRYPOINT ["/entrypoint"]

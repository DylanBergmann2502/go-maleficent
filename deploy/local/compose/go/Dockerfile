# deploy/local/compose/go/Dockerfile
FROM golang:1.25-alpine

# Install build dependencies
RUN apk add --no-cache \
    git \
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

# Disable VCS stamping for Docker environment (no git context)
ENV GOFLAGS="-buildvcs=false"

WORKDIR /app

# Copy go mod and go.sum files (if exists)
COPY go.mod ./
COPY go.su[m] ./

# Download dependencies
RUN go mod download

# Install Air globally for hot reload
RUN go install github.com/air-verse/air@latest

# Install goimports for formatting
RUN go install golang.org/x/tools/cmd/goimports@latest

# Install golangci-lint for linting (v2)
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

# Copy source code
COPY . .

# Create tmp directory for Air
RUN mkdir -p tmp

# Copy scripts
COPY ./deploy/local/compose/go/entrypoint /entrypoint
COPY ./deploy/local/compose/go/start /start

# Remove Windows line endings and make scripts executable
RUN sed -i 's/\r$//g' /entrypoint \
    && sed -i 's/\r$//g' /start \
    && chmod +x /entrypoint /start

EXPOSE 8000

ENTRYPOINT ["/entrypoint"]

# Docker Deployment Guide

This guide explains how to build and deploy the Living Smart Contracts application using Docker.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Git

## Quick Start

### 1. Build and Run with Docker Compose

```bash
# Clone the repository (if not already done)
git clone https://github.com/First-Genesis/Living-Smart-Contracts.git
cd Living-Smart-Contracts

# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f living-contracts
```

### 2. Build Docker Image Manually

```bash
# Build the Docker image
docker build -t living-smart-contracts .

# Run the container
docker run -d \
  --name living-contracts \
  -p 8080:8080 \
  -p 9090:9090 \
  -e LOG_LEVEL=info \
  living-smart-contracts
```

## Available Services

After starting with docker-compose, the following services will be available:

- **Living Smart Contracts API**: http://localhost:8080
  - Health check: http://localhost:8080/health
  - API info: http://localhost:8080/api/info
  - Contracts API: http://localhost:8080/api/contracts

- **Prometheus Metrics**: http://localhost:9091
- **Grafana Dashboard**: http://localhost:3000 (admin/admin)

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `ACTOR_SYSTEM_NAME` | `living-contracts` | Actor system identifier |
| `HTTP_PORT` | `8080` | HTTP API port |
| `METRICS_PORT` | `9090` | Metrics endpoint port |

## Docker Image Details

- **Base Image**: Alpine Linux (minimal footprint)
- **Go Version**: 1.21
- **Architecture**: Multi-stage build for optimized size
- **Security**: Runs as non-root user
- **Health Check**: Built-in health monitoring

## Monitoring

The docker-compose setup includes:

1. **Prometheus**: Metrics collection and storage
2. **Grafana**: Visualization and dashboards
3. **Health Checks**: Automatic container health monitoring

## Troubleshooting

### Build Issues

```bash
# Clean build (remove cache)
docker build --no-cache -t living-smart-contracts .

# Check build logs
docker build -t living-smart-contracts . 2>&1 | tee build.log
```

### Runtime Issues

```bash
# Check container logs
docker logs living-contracts

# Access container shell
docker exec -it living-contracts /bin/sh

# Check health status
curl http://localhost:8080/health
```

### Port Conflicts

If ports 8080, 9090, 9091, or 3000 are already in use:

```bash
# Modify docker-compose.yml ports section
# Example: Change "8080:8080" to "8081:8080"
```

## Production Deployment

For production deployment:

1. **Use environment-specific configurations**
2. **Set up proper logging and monitoring**
3. **Configure resource limits**
4. **Use Docker secrets for sensitive data**
5. **Set up backup and recovery procedures**

Example production docker-compose override:

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  living-contracts:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
    environment:
      - LOG_LEVEL=warn
    restart: always
```

Run with: `docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d`

## Development

For development with live reload:

```bash
# Run in development mode
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Or run locally
go run cmd/server/main.go
```

#!/bin/bash

# Build script for Reel Farm Docker image

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
IMAGE_NAME="reel-farm-app"
TAG="latest"
PORT="3000"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -t|--tag)
      TAG="$2"
      shift 2
      ;;
    -p|--port)
      PORT="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [OPTIONS]"
      echo "Options:"
      echo "  -t, --tag TAG     Docker image tag (default: latest)"
      echo "  -p, --port PORT   Port to run container on (default: 3000)"
      echo "  -h, --help        Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option $1"
      exit 1
      ;;
  esac
done

echo -e "${YELLOW}🐳 Building Docker image: ${IMAGE_NAME}:${TAG}${NC}"

# Build the Docker image
docker build -t "${IMAGE_NAME}:${TAG}" .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Docker image built successfully!${NC}"
    echo -e "${YELLOW}📦 Image: ${IMAGE_NAME}:${TAG}${NC}"
    echo ""
    echo -e "${YELLOW}🚀 To run the container:${NC}"
    echo "docker run -p ${PORT}:3000 ${IMAGE_NAME}:${TAG}"
    echo ""
    echo -e "${YELLOW}🧪 To test the health endpoint:${NC}"
    echo "curl http://localhost:${PORT}/health"
    echo ""
    echo -e "${YELLOW}📊 To see container logs:${NC}"
    echo "docker logs <container_id>"
else
    echo -e "${RED}❌ Docker build failed!${NC}"
    exit 1
fi

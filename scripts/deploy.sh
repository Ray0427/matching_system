#!/bin/bash

echo "Deploying Matching System API..."

# Build the application
./scripts/build.sh

if [ $? -ne 0 ]; then
    echo "Build failed, deployment aborted!"
    exit 1
fi

# Build Docker image
docker build -t matching-system .

if [ $? -eq 0 ]; then
    echo "Docker image built successfully!"
    
    # Stop existing containers
    docker-compose down
    
    # Start new containers
    docker-compose up -d
    
    echo "Deployment completed successfully!"
    echo "API is running at http://localhost:8080"
    echo "Health check: http://localhost:8080/health"
    echo "Swagger docs: http://localhost:8080/swagger/index.html"
else
    echo "Docker build failed!"
    exit 1
fi

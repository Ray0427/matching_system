#!/bin/bash

echo "Building Matching System API..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the application
go build -o bin/api cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created at bin/api"
else
    echo "Build failed!"
    exit 1
fi

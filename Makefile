.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Build Docker image
docker-build:
	docker build -t matching-system .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Stop Docker Compose
docker-stop:
	docker-compose down

# Generate Swagger docs
swagger:
	swag init -g cmd/api/main.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

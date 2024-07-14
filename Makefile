# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	
	@go build -o main.go
# @go build -o main main.go

# Run the application
run:
	@go run main.go

# Create DB container
docker-run:
	docker compose build
	docker compose up 


# Shutdown DB container
docker-down:d
	ocker-compose down

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main



.PHONY: all build run test clean

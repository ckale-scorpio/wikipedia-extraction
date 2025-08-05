.PHONY: build run test clean install

# Build the application
build:
	go build -o bin/wikipedia-extraction .

# Run the application
run: build
	./bin/wikipedia-extraction

# Install dependencies
install:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run with example
example: build
	./bin/wikipedia-extraction extract "https://en.wikipedia.org/wiki/Go_(programming_language)" --output example.json --format json

# Build for different platforms
build-all: clean
	GOOS=linux GOARCH=amd64 go build -o bin/wikipedia-extraction-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/wikipedia-extraction-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o bin/wikipedia-extraction-windows-amd64.exe .

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  install    - Install dependencies"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  example    - Run with example Wikipedia page"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  help       - Show this help" 
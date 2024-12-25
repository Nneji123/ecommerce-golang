.PHONY: help vet format test test-cover test-cover-report lint generate-docs all build clean run build-linux build-windows build-macos docker-build docker-run

# Define variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=./bin/api
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_MACOS=$(BINARY_NAME)_macos

default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

vet: ## Runs go vet
	go vet ./...

format: ## Runs go fmt
	gofmt -s -w .

test: ## Runs the unit tests
	go test -v -race -timeout 5m ./...

test-cover: ## Outputs the coverage statistics
	go test -v -race -timeout 5m ./... -coverprofile coverage.out
	go tool cover -func coverage.out
	rm coverage.out

test-cover-report: ## A html report of the coverage statistics
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html
	open coverage.html

lint: ## Runs the linter
	golangci-lint -v run ./...

generate-docs: ## Generates API Documentation using swagger
	## Make sure to install swag first
	swag fmt
	swag init --parseDependency --parseDepth 3 -o ./docs -g ./cmd/api/main.go


all: ## Clean and build binary
	format
	generate-docs
	clean
	build

build: ## Build Binary
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/api/main.go

clean: ## Clean binaries
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(BINARY_UNIX) $(BINARY_WINDOWS) $(BINARY_MACOS)

run: ## Run binaries after build
	go run cmd/api/main.go

run-tasks: ## Run asynchronous task queue with Air
	air --build.cmd "go build -o bin/tasks.exe cmd/tasks/main.go" --build.bin "/bin/tasks.exe"

build-linux: ## Cross compilation for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/api/main.go

build-windows: ## Cross compilation for Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v ./cmd/api/main.go

build-macos: ## Cross compilation for macOS
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MACOS) -v ./cmd/api/main.go

docker-build: ## Build the Docker image
	docker build -t api .

docker-run: ## Run the Docker container
	docker run -p 8080:8080 api

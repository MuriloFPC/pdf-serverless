.PHONY: build deploy local clean test

# Rodar testes
test:
	@echo "Running unit tests..."
	go test -v ./...

# Build para AWS Lambda (Linux ARM64)
build:
	@echo "Building binaries for AWS Lambda..."
ifeq ($(OS),Windows_NT)
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='arm64'; $$env:CGO_ENABLED='0'; go build -o bin/api/bootstrap ./cmd/api"
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='arm64'; $$env:CGO_ENABLED='0'; go build -o bin/worker/bootstrap ./cmd/worker"
else
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/api/bootstrap ./cmd/api
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/worker/bootstrap ./cmd/worker
endif

# Targets para o SAM Build (BuildMethod: makefile)
build-APIFunction:
ifeq ($(OS),Windows_NT)
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='arm64'; $$env:CGO_ENABLED='0'; go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/api"
else
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/api
endif

build-WorkerFunction:
ifeq ($(OS),Windows_NT)
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='arm64'; $$env:CGO_ENABLED='0'; go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/worker"
else
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/worker
endif

# Deploy usando AWS SAM
deploy: test build
	@echo "Deploying to AWS via SAM..."
	sam build
	sam deploy --guided

# Rodar a API localmente (modo memória por padrão)
local:
	@echo "Running API locally in memory mode..."
	USE_MEMORY=true go run cmd/api/main.go

# Limpar binários
clean:
	@echo "Cleaning up..."
	rm -rf bin/

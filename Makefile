.PHONY: build deploy local clean

# Build para AWS Lambda (Linux ARM64)
build:
	@echo "Building binaries for AWS Lambda..."
	GOOS=linux GOARCH=arm64 go build -o bin/api/bootstrap ./cmd/api
	GOOS=linux GOARCH=arm64 go build -o bin/worker/bootstrap ./cmd/worker

# Deploy usando AWS SAM
deploy: build
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

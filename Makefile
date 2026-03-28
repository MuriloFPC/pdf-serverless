.PHONY: build deploy local clean

# Build para AWS Lambda (Linux ARM64)
build:
	@echo "Building binaries for AWS Lambda..."
	@set GOOS=linux& set GOARCH=arm64& set CGO_ENABLED=0& go build -o bin/api/bootstrap ./cmd/api
	@set GOOS=linux& set GOARCH=arm64& set CGO_ENABLED=0& go build -o bin/worker/bootstrap ./cmd/worker

# Targets para o SAM Build (BuildMethod: makefile)
build-APIFunction:
	@set GOOS=linux& set GOARCH=arm64& set CGO_ENABLED=0& go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/api

build-WorkerFunction:
	@set GOOS=linux& set GOARCH=arm64& set CGO_ENABLED=0& go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/worker

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

# Makefile for Reel Farm Server

.PHONY: help pull-schema generate clean generate-api

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

pull-schema: ## Pull schema from database and generate SQLC code
	@echo "🔄 Pulling schema from database..."
	@./scripts/pull-schema.sh

generate: ## Generate SQLC code only (requires schema to exist)
	@echo "🔄 Generating SQLC code..."
	@sqlc generate
	@echo "✅ SQLC code generated"

generate-api: ## Generate OpenAPI code
	@echo "🔄 Generating OpenAPI code..."
	@mkdir -p internal/api
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config ../openapi/config.yaml ../openapi/api.yaml
	@echo "✅ OpenAPI code generated"

clean: ## Clean generated files
	@echo "🧹 Cleaning generated files..."
	@rm -rf db/
	@rm -f sql/schema/schema.sql
	@rm -f internal/api/generated.go
	@echo "✅ Cleaned"

dev: ## Run development server
	@echo "🚀 Starting development server..."
	@go run cmd/main.go

build: ## Build the application
	@echo "🔨 Building application..."
	@go build -o bin/server cmd/main.go
	@echo "✅ Built successfully"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test ./...

install-deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

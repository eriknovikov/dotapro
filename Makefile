.PHONY: build deploy clean help run-api run-ui build-ui-local

# Lambda function names
API_FUNCTION_NAME := dotapro-lambda-api
SCRAPER_FUNCTION_NAME := dotapro-lambda-scraper
DB_URL := postgres://postgres:admin@172.17.0.1:15432/dotapro?sslmode=disable
# Build directory
BUILD_DIR := .build

# Help target
help:
	@echo "Available targets:"
	@echo "  make build      - Build and zip both Lambda functions"
	@echo "  make deploy     - Deploy both Lambda functions to AWS"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make run-api    - Run API locally (requires .env.local)"
	@echo "  make run-ui     - Run UI locally (requires ui/.env.local)"
	@echo "  make build-ui-local - Build UI for local testing (uses ui/.env.production.local)"
	@echo ""
	@echo "Individual targets:"
	@echo "  make build-api      - Build and zip API function"
	@echo "  make build-scraper  - Build and zip Scraper function"
	@echo "  make deploy-api     - Deploy API function"
	@echo "  make deploy-scraper - Deploy Scraper function"

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Build API function
build-api: $(BUILD_DIR)
	@echo "🔨 Building API function..."
	@cd api && \
		GOOS=linux GOARCH=arm64 go build -o bootstrap . && \
		cd .. && \
		zip -j $(BUILD_DIR)/api.zip api/bootstrap && \
		rm api/bootstrap && \
		echo "✅ API function built and zipped"

# Build Scraper function
build-scraper: $(BUILD_DIR)
	@echo "🔨 Building Scraper function..."
	@cd scraper && \
		GOOS=linux GOARCH=arm64 go build -o bootstrap . && \
		cd .. && \
		zip -j $(BUILD_DIR)/scraper.zip scraper/bootstrap && \
		rm scraper/bootstrap && \
		echo "✅ Scraper function built and zipped"

# Build both functions
build: build-api build-scraper
	@echo "🎉 All functions built successfully"

# Deploy API function
deploy-api: build-api
	@echo "🚀 Deploying API function..."
	@aws lambda update-function-code \
		--function-name $(API_FUNCTION_NAME) \
		--zip-file fileb://$(BUILD_DIR)/api.zip > /dev/null
	@echo "✅ API function deployed successfully at $$(date +%H:%M.%S)"

# Deploy Scraper function
deploy-scraper: build-scraper
	@echo "🚀 Deploying Scraper function..."
	@aws lambda update-function-code \
		--function-name $(SCRAPER_FUNCTION_NAME) \
		--zip-file fileb://$(BUILD_DIR)/scraper.zip > /dev/null
	@echo "✅ Scraper function deployed successfully at $$(date +%H:%M.%S)"

# Deploy both functions
deploy: deploy-api deploy-scraper
	@echo "🎉 All functions deployed successfully"


migrate-up: 
	migrate -path database/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" down

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run API locally
run-api:
	@echo "🚀 Running API locally..."
	@cd api && $(MAKE) run

# Run UI locally
run-ui:
	@echo "🚀 Running UI locally..."
	@cd ui && pnpm dev

# Build UI for local testing (uses localhost API)
build-ui-local:
	@echo "🔨 Building UI for local testing..."
	@cd ui && pnpm build
	@echo "✅ UI built successfully (using localhost API from .env.production.local)"

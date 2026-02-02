.PHONY: build deploy clean help

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
	@cd api/src && \
		GOOS=linux GOARCH=amd64 go build -o bootstrap . && \
		cd ../../ && \
		zip -j $(BUILD_DIR)/api.zip api/src/bootstrap && \
		rm api/src/bootstrap && \
		echo "✅ API function built and zipped"

# Build Scraper function
build-scraper: $(BUILD_DIR)
	@echo "🔨 Building Scraper function..."
	@cd scraper/src && \
		GOOS=linux GOARCH=amd64 go build -o bootstrap . && \
		cd ../../ && \
		zip -j $(BUILD_DIR)/scraper.zip scraper/src/bootstrap && \
		rm scraper/src/bootstrap && \
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

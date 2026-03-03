.PHONY: build deploy clean help deploy-ui

#VARS
API_FUNCTION_NAME := dotapro-lambda-api
SCRAPER_FUNCTION_NAME := dotapro-lambda-scraper
CF_DISTRIBUTION_ID := E1WWIEAKQHKZ85
S3_BUCKET_NAME := dotapro-ui
BUILD_DIR := .build

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


# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

deploy-ui:
	@echo "--- Building UI ---"
	cd ui && pnpm build
	
	@echo "--- Uploading Assets (Cache Forever) ---"
	# Syncs everything EXCEPT index.html with a 1-year cache
	aws s3 sync ui/dist s3://$(S3_BUCKET_NAME) \
		--delete \
		--exclude "index.html" \
		--cache-control "max-age=31536000, public, immutable"
	
	@echo "--- Uploading Entry Point (No Cache) ---"
	# Uploads index.html specifically with headers to force re-validation
	aws s3 cp ui/dist/index.html s3://$(S3_BUCKET_NAME)/index.html \
		--cache-control "no-cache, no-store, must-revalidate"
	
	@echo "--- Invalidating CloudFront Cache ---"
	# Tells CloudFront to clear its edge cache immediately
	aws cloudfront create-invalidation \
		--distribution-id $(CF_DISTRIBUTION_ID) \
		--paths "/*"
	
	@echo "--- Deployment Complete! ---"
.PHONY: tidy test test-verbose test-cover test-series test-filters

tidy:
	go mod tidy

test:
	go test ./... -short

test-verbose:
	go test ./... -v -short

test-cover:
	go test ./... -cover -short

test-series:
	go test ./series/... -v -short

test-filters:
	go test ./filtersmetadata/... -v -short
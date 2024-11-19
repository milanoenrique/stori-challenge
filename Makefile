BINARY=job
.PHONY: all
all: deps test test_local test_local_coverage clean build


.PHONY: deps
deps:
	@echo "Ensuring deps"
	@go mod tidy

.PHONY: build
build:
	@echo "Building"
	go build -o /build/${BINARY} cmd/job/*.go

.PHONY: clean
clean:
	go clean -testcache 

.PHONY: test_local
test_local:
	go test -cover ./...

.PHONY: test_local_coverage
test_local_coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

.PHONY: test
test:
	@echo "Running tests"
	go test -json > report.json -cover -coverprofile=coverage.out -race ./...

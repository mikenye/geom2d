# Default number of workers for fuzzing (can be overridden via `make FUZZ_WORKERS=X`)
FUZZ_WORKERS ?= 4

# Define help message (lines starting with ## will be parsed for help output)
.PHONY: help
help:  ## Show available make commands
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[1;32m%-25s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test:  ## Run all unit tests 100 times
	go test -count=100 ./...

.PHONY: test_short
test_short:  ## Run short unit tests
	go test -short ./...

.PHONY: fuzz_findintersections
fuzz_findintersections:  ## Run fuzz testing for linesegment intersection functions (long running). Fuzzing is performed for 2, 3 and 4 line segments, with 1M executions each.
	go test -run=FuzzFindIntersections_Int_2Segments -fuzz=FuzzFindIntersections_Int_2Segments -fuzztime=1000000x -parallel=$(FUZZ_WORKERS) ./linesegment/...
	go test -run=FuzzFindIntersections_Int_3Segments -fuzz=FuzzFindIntersections_Int_3Segments -fuzztime=1000000x -parallel=$(FUZZ_WORKERS) ./linesegment/...
	go test -run=FuzzFindIntersections_Int_4Segments -fuzz=FuzzFindIntersections_Int_4Segments -fuzztime=1000000x -parallel=$(FUZZ_WORKERS) ./linesegment/...

.PHONY: lint
lint:  ## Run golangci-lint, staticcheck and go vet
	go tool github.com/golangci/golangci-lint/cmd/golangci-lint run
	go tool honnef.co/go/tools/cmd/staticcheck ./...
	go vet ./...

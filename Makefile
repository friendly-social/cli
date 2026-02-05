# general variables
OUT_DIR := ./out
COVER_FILE := $(OUT_DIR)/coverage.out

# testing parameters
VERBOSE ?= 0
COVERAGE ?= 0

# recipes list
.PHONY: test lint fmt clean

# source files for tracking changes
SRC := $(shell find . -type f -name '*.go')

# public recipe for testing (with or without coverage)
test: $(OUT_DIR)/test.cache.$(COVERAGE)

# public recipe for formatting
fmt: $(OUT_DIR)/fmt.cache

# public recipe for linting
lint: $(OUT_DIR)/lint.cache

# cleaning up garbage
clean:
	@echo ">> Cleaning up..."
	@rm -rf $(BIN_DIR) $(OUT_DIR)
	@echo ">> Cleaned."

# creating output directory if it's missing
$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

# formatting source code
$(OUT_DIR)/fmt.cache: $(SRC) | $(OUT_DIR)
	@echo ">> Formatting..."
	@gofmt -s -w .
	@echo ">> Formatted."
	@touch $@

# linting source code
$(OUT_DIR)/lint.cache: $(SRC) | $(OUT_DIR)
	@echo ">> Linting..."
	@golangci-lint run ./... --timeout=5m > $(OUT_DIR)/lint.log || { \
		cat $(OUT_DIR)/lint.log; \
		exit 1; \
	};
	@echo ">> Linted."
	@touch $@

# testing source code (with or without coverage)
$(OUT_DIR)/test.cache.$(COVERAGE): $(SRC) | $(OUT_DIR)
	@set -e; \
	if [ "$(COVERAGE)" = "1" ]; then \
		echo ">> Testing with coverage..."; \
		go test ./... -coverprofile=$(COVER_FILE) -covermode=atomic > $(OUT_DIR)/test.log || { \
			cat $(OUT_DIR)/test.log; \
			exit 1; \
		}; \
		grep -v "/dummy" $(COVER_FILE) > $(COVER_FILE).tmp && mv $(COVER_FILE).tmp $(COVER_FILE); \
		COVERAGE_OUTPUT=$$(go tool cover -func=$(COVER_FILE)); \
		if [ "$(VERBOSE)" = "1" ]; then \
			echo "$$COVERAGE_OUTPUT" | grep -v "total:"; \
		fi; \
		echo "$$COVERAGE_OUTPUT" | grep "total:" | awk '{print ">> Test coverage:", $$3}'; \
	else \
		echo ">> Testing..."; \
		go test ./... > $(OUT_DIR)/test.log || { \
			cat $(OUT_DIR)/test.log; \
			exit 1; \
		}; \
		if [ "$(VERBOSE)" = "1" ]; then \
			cat $(OUT_DIR)/test.log; \
		fi; \
		echo ">> All tests passed."; \
	fi;
	@touch $@

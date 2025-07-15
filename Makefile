APP_NAME := asaas-cli

MAIN_PACKAGE := ./cmd/main.go

BIN_DIR := bin

BINARY := $(BIN_DIR)/$(APP_NAME)

GO_BUILD_FLAGS := -ldflags="-s -w" -tags=netgo
GO_TEST_FLAGS := -v -race -cover

INSTALL_PATH := $(GOBIN)
ifeq ($(GOBIN),)
    INSTALL_PATH := /usr/local/bin
endif

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build $(GO_BUILD_FLAGS) -o $(BINARY) $(MAIN_PACKAGE) && \
	chmod +x $(BINARY) && \
	echo "Successfully built $(APP_NAME) to $(BINARY)"

.PHONY: test
test:
	@echo "Running tests..."
	@go test $(GO_TEST_FLAGS) ./...

.PHONY: install
install: build
	@echo "Installing $(APP_NAME) to $(INSTALL_PATH)..."
	@if [ -w "$(INSTALL_PATH)" ]; then \
		cp $(BINARY) $(INSTALL_PATH); \
	else \
		echo "Insufficient permissions. Trying with sudo..."; \
		sudo cp $(BINARY) $(INSTALL_PATH); \
	fi
	@echo "Successfully installed $(APP_NAME) to $(INSTALL_PATH)/$(APP_NAME)"

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@go clean -modcache
	@echo "Cleanup complete."

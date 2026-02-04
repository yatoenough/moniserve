APP_NAME=moniserve
TARGET_DIR=bin

default: run

build-debug:
	@go build -o $(TARGET_DIR)/debug/$(APP_NAME) cmd/main.go

build-release:
	@go build -ldflags "-s -w" -o $(TARGET_DIR)/release/$(APP_NAME) cmd/main.go

clean:
	@rm -rf $(TARGET_DIR)

run: build-debug
	@./$(TARGET_DIR)/debug/$(APP_NAME)

watch:
	@air

test:
	@go test -v ./internal/...

.PHONY: default build-debug build-release clean run watch test

default: lint

.PHONY: lint
lint:
	go mod tidy
	go fmt

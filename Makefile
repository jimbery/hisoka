.PHONY: all serverless deps docker docker-cgo clean docs test test-race test-integration fmt lint install deploy-docs

install:

deps:
	@go mod tidy

lint:
	golangci-lint run

fmtfix:
	gofmt --w .

dev:
	go run cmd/main.go --local

test:
	go test ./...

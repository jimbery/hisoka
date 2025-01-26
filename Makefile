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
	# Run the Docker container in the background and set up the trap in the same shell
	@sh -c 'docker run --name test-postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres && \
	trap "docker stop test-postgres; docker rm test-postgres" EXIT && \
	sleep 5 && \
	go test ./...'

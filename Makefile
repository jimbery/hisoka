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
	# Run the Docker container in the background
	docker run --name test-postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres

	# Set a trap to clean up the Docker container if the script exits
	trap 'docker stop test-postgres; docker rm test-postgres' EXIT

	until docker exec test-postgres pg_isready -U postgres; do
	  echo "Waiting for PostgreSQL to be ready..."
	  sleep 1
	done

	# Run the Go tests
	go test ./...

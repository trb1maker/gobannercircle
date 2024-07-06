generate-api:
	@go generate ./...

lint-install:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.59.0

lint:
	@./bin/golangci-lint run ./...

test:
	@go test -v -race -count 10 ./...

build-app:
	@go build -o bin/service main.go

test-logic:
	@go test -tags=logic -timeout=5m ./internal/app/logic_test.go
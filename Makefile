include .env

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

build-container:
	@docker build -t service:dev -f build/Dockerfile .

run: build-app
	@STORAGE_HOST=${STORAGE_HOST} \
	STORAGE_PORT=${STORAGE_PORT} \
	STORAGE_DBNAME=${STORAGE_DBNAME} \
	STORAGE_USER=${STORAGE_USER} \
	STORAGE_PASSWORD=${STORAGE_PASSWORD} \
	./bin/service migrate

	@STORAGE_HOST=${STORAGE_HOST} \
	STORAGE_PORT=${STORAGE_PORT} \
	STORAGE_DB=${STORAGE_DB} \
	STORAGE_USER=${STORAGE_USER} \
	STORAGE_PASSWORD=${STORAGE_PASSWORD} \
	NOTIFY_HOST=${NOTIFY_HOST} \
	NOTIFY_PORT=${NOTIFY_PORT} \
	NOTIFY_TOPIC=${NOTIFY_TOPIC} \
	NOTIFY_PARTITION=${NOTIFY_PARTITION} \
	APP_HOST=${APP_HOST} \
	APP_PORT=${APP_PORT} \
	LOGGER_LEVEL=${LOGGER_LEVEL} \
	./bin/service start

up:
	@docker compose -f deploy/docker-compose.yaml up -d

down:
	@docker compose -f deploy/docker-compose.yaml down
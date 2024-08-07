.PHONY: default dev run clear install test docs build db_generate migrations_up migrations_down migrations_create lint

include .env

default: dev

rpc_dev:
	@air -c .rpc.air.toml
rest_dev:
	@air -c .rest.air.toml
run:
	@go run ./cmd/restapi
clear:
	@rm ./tmp/main
install:
	@go mod download && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && go install github.com/swaggo/swag/cmd/swag@latest && go install github.com/pressly/goose/v3/cmd/goose@latest && go install github.com/air-verse/air@latest
test:
	@ENV_FILEPATH=$(ENV_FILEPATH) go test ./internal/domain/usecase
coverage:
	@ENV_FILEPATH=$(ENV_FILEPATH) go test ./internal/domain/usecase -coverprofile ./tmp/test_coverage.out && go tool cover -html=tmp/test_coverage.out
docs:
	@swag init -g ./cmd/restapi/main.go -o ./docs
build:
	@GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ./tmp/restapi ./cmd/restapi
db_generate:
	@sqlc generate
grpc_generate:
	@protoc --proto_path=proto proto/*.proto  --go_out=. --go-grpc_out=.
migrations_up:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up
migrations_down:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down
migrations_create:
	@goose create $(NAME) sql
lint:
	@golangci-lint run && nilaway ./...

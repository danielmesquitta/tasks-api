.PHONY: default rest_dev rpc_dev install test coverage docs build db_gen grpc_gen migrations_up migrations_down migrations_create lint

include .env

default: rest_dev

rest_dev:
	@air -c .rest.air.toml
rpc_dev:
	@air -c .rpc.air.toml
clear:
	@find ./tmp -mindepth 1 ! -name '.gitkeep' -delete
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
db_gen:
	@sqlc generate
grpc_gen:
	@protoc --proto_path=proto proto/*.proto  --go_out=. --go-grpc_out=.
migrations_up:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up
migrations_down:
	@goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down
migrations_create:
	@goose create $(NAME) sql
lint:
	@golangci-lint run && nilaway ./...

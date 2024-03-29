.DEFAULT_GOAL := build

.PHONY: check_migrate add_migration check_lint check_sqlc lint build run install reset gen local_build dev

check_migrate:
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

add_migration: check_migrate
	 migrate create -dir=internal/db/migrations/ -seq -ext sql $(name)

check_lint:
	which staticcheck || go install honnef.co/go/tools/cmd/staticcheck@latest

check_sqlc:
	which sqlc || go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.14.0

lint: check_lint check_sqlc
	go vet ./...
	staticcheck ./...
	cd internal/db && sqlc compile

build:
	@sh -c './scripts/build.sh'

run: build
	./blog

install:
	./bin/blog install

reset:
	./bin/blog reset

gen:
	go generate ./...


check_air:
	which air || go install github.com/cosmtrek/air@v1.40.3

local_build: check_sqlc gen
	@sh -c './scripts/build.sh'

# Hot reload for development
dev: check_air
	air

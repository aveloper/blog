.DEFAULT_GOAL := build

check_migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

add_migration: check_migrate
	 migrate create -dir=migrations/ -seq -ext sql $(name)

build:
	go build -o build/

run: build
	./blog
.DEFAULT_GOAL := build

check_migrate:
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

add_migration: check_migrate
	 migrate create -dir=internal/db/migrations/ -seq -ext sql $(name)

check_lint:
	which staticcheck || go install honnef.co/go/tools/cmd/staticcheck@latest

lint: check_lint
	go vet ./...
	staticcheck ./...

build:
	@sh -c './scripts/build.sh'

run: build
	./blog

install:
	./output/blog install

reset:
	./output/blog reset
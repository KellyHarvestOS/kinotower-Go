include .env
export

run:
	go run ./cmd/api

build:
	go build ./...

seed:
	set AUTO_SEED=true&& set SEED_ONLY=true&& go run ./cmd/api

migrate:
	set AUTO_MIGRATE=true&& set AUTO_SEED=false&& set MIGRATE_ONLY=true&& go run ./cmd/api

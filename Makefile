.PHONY: migration sqlc live

## migration: Create a new migration file with the provided name
migration:
ifndef name
	$(error name is required, e.g. `make migration name=create_documents`)
endif
	goose -dir "internal/database/migrations" create $(name) sql

## sqlc: Generate Go code from SQL queries and migrations
sqlc:
	go tool sqlc generate

## live: Run the API with air live reloading
live:
	air \
		--build.cmd "go build -o tmp/api ./cmd/api" \
		--build.entrypoint "tmp/api" \
		--build.exclude_dir "tmp,data" \
		--build.include_ext "go,sql,env"

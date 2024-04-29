dev: 
	@APP_ENV=development go run main.go

gen:
	@sqlc generate

migrate-up:
	@cd database/schema && goose postgres $(DATABASE_URL) up

migrate-down:
	@cd database/schema && goose postgres $(DATABASE_URL) down

tidy:
	@go mod vendor && go mod tidy

test:
	@go test -v ./...

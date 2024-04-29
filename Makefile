dev: 
	@APP_ENV=development go run main.go

test:
	@go test -v ./...

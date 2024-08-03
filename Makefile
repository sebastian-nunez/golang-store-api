build:
	@go build -o bin/golang-store-api cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/golang-store-api
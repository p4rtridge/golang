include .env
export

.PHONY: start
start:
	docker-compose up -d --build

.PHONY: stop
stop:
	docker-compose rm -vfs && docker rmi gofiber:latest

.PHONY: build
build:
	go build -o ./bin/application ./cmd/app/main.go

.PHONY: run
run: build
	./bin/application

.PHONY: test
test:
	go test -count=1 -coverprofile ./bin/cover.out ./services/...

.PHONY: cover
cover: test
	go tool cover -html ./bin/cover.out -o cover.html && brave cover.html

.PHONY: migrate
migrate:
	 go run ./migrations/migrate.go

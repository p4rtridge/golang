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
	go test -count=1 ./services/user/test/... ./services/product/test/... ./services/order/test/... ./services/auth/test/...

.PHONY: migrate
migrate:
	 go run ./migrations/migrate.go

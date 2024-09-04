include .env
export

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down --remove-orphans

run:
	go mod tidy && go mod download && \
	go run ./cmd/app/main.go

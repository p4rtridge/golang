.PHONY: start
start:
	docker-compose up --build --detach && docker image prune -f

.PHONY: stop
stop:
	docker-compose down --rmi local


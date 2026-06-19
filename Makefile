.PHONY: docker-up docker-down test

docker-up:
	docker compose up --build

docker-down:
	docker compose down

test:
	go test ./...

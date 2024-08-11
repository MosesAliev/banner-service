build:
	docker build --tag my-app:v1 .

run:
	docker compose up -d

test: run
	docker compose exec web go test ./internal/tests/...
	docker compose down

server: run
	docker compose exec web go run ./...
build:
	docker build --tag my-app:v1 .



test:
	docker compose down
	docker compose up -d
	docker compose exec web go test ./internal/tests/...
	docker compose down

server:
	docker compose down
	docker compose up -d
	docker compose exec web go run ./...


kill:
	docker compose down
run:
	go run backend/main.go
lint:
	cd backend && golangci-lint run
build:
	docker compose build
up:
	docker compose up --build
down:
	docker compose down -v
logs:
	docker compose logs -f
clean:
	docker volume prune -f
run:
	go run backend/main.go
lint:
	cd backend && golangci-lint run
run docker:
	docker compose up -d
.PHONY: run lint build up down logs clean

run:
	@echo "Запуск backend в режиме разработки..."
	go run backend/main.go

lint:
	@echo "Запуск golangci-lint для backend..."
	cd backend && golangci-lint run

build:
	@echo "Сборка Docker-образов..."
	docker-compose build

up:
	@echo "Запуск всех сервисов через Docker Compose..."
	docker-compose up --build

down:
	@echo "Остановка и удаление контейнеров и томов..."
	docker-compose down -v

logs:
	@echo "Показ логов всех сервисов..."
	docker-compose logs -f

clean:
	@echo "Очистка неиспользуемых томов Docker..."
	docker volume prune -f

.PHONY: help db-up db-down run test build clean

# Показать помощь
help:
	@echo "Доступные команды:"
	@echo "  db-up        - Запустить PostgreSQL в Docker"
	@echo "  db-down      - Остановить PostgreSQL Docker контейнер"
	@echo "  run          - Запустить приложение локально"
	@echo "  test         - Запустить тесты"
	@echo "  build        - Собрать приложение"
	@echo "  clean        - Очистить сгенерированные файлы"

# Запустить PostgreSQL через Docker
db-up:
	@echo "Запуск PostgreSQL..."
	docker compose up postgres -d
	@echo "PostgreSQL доступен на localhost:5432"

# Остановить PostgreSQL Docker
db-down:
	@echo "Остановка PostgreSQL..."
	docker compose down

# Запустить приложение локально
run:
	@echo "Запуск приложения..."
	go run cmd/server/main.go

# Запустить тесты
test:
	@echo "Запуск тестов..."
	go test ./... -v

# Собрать приложение
build:
	@echo "Сборка приложения..."
	go build -o bin/orders cmd/server/main.go
	@echo "Приложение собрано в bin/orders"

# Очистить сгенерированные файлы
clean:
	@echo "Очистка..."
	rm -rf bin/
	rm -f server
	docker compose down -v

.PHONY: help run stop test clean build docker-up docker-down

# Показать помощь
help:
	@echo "Доступные команды:"
	@echo "  run          - Запустить все сервисы через Docker Compose"
	@echo "  stop  		  - Остановить все сервисы Docker"
	@echo "  test         - Запустить тесты"
	@echo "  build        - Собрать приложение"
	@echo "  clean        - Очистить сгенерированные файлы"

# Запустить через Docker Compose
run:
	@echo "Запуск приложения через Docker Compose..."
	docker compose up -d
	@echo "Приложение доступно на http://localhost:8080"
	@echo "Проверка здоровья: curl http://localhost:8080/health"

# Остановить Docker сервисы
stop:
	@echo "Остановка Docker сервисов..."
	docker compose down

# Запустить тесты
test:
	@echo "Запуск тестов..."
	go test ./... -v

# Собрать приложение
build:
	@echo "Сборка приложения..."
	go build -o bin/orders cmd/server/main.go
	@echo "Приложение собрано в bin/orders"

# Генерировать моки
generate:
	@echo "Генерация моков..."
	go generate ./internal/domain/repositories/

# Очистить сгенерированные файлы
clean:
	@echo "Очистка..."
	rm -rf bin/
	rm -f server
	docker compose down -v

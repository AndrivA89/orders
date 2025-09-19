# Orders Service

Микросервис для управления заказами.

## Архитектура

1. **Domain Layer** (`internal/domain/`) - Бизнес-логика и сущности
   - `entities/` - Доменные модели (User, Product, Order)
   - `repositories/` - Интерфейсы репозиториев
   - `services/` - Интерфейсы бизнес-сервисов

2. **Application Layer** (`internal/application/`) - Прикладная логика
   - `services/` - Реализация бизнес-логики

3. **Infrastructure Layer** (`internal/infrastructure/`) - Внешние зависимости
   - `database/` - Подключение к БД и модели
   - `repositories/` - Реализация репозиториев
   - `config/` - Конфигурация приложения

4. **Transport Layer** (`internal/transport/`) - Внешние интерфейсы
   - `http/` - REST API handlers, middleware, DTO

## Основные сущности

### User (Пользователь)
- `id` - UUID
- `firstname` - Имя
- `lastname` - Фамилия
- `fullname` - Полное имя (firstname + lastname)
- `age` - Возраст (не младше 18 лет)
- `is_married` - Семейное положение
- `password` - Пароль (не меньше 8 символов, хешируется)

### Product (Товар)
- `id` - UUID
- `description` - Описание
- `tags` - Теги (JSON array)
- `quantity` - Количество на складе
- `price` - Цена в копейках

### Order (Заказ)
- `id` - UUID
- `user_id` - ID пользователя
- `status` - Статус (pending, confirmed, cancelled, completed)
- `total` - Общая сумма
- `items` - Позиции заказа с историчностью цен

## Функциональность

### Основные возможности
- Регистрация пользователя с валидацией возраста (18+) и пароля (8+ символов)
- Создание и управление товарами с тегами и количеством
- Создание заказов с проверкой наличия товара на складе
- Историчность заказов - ProductSnapshot сохраняет цены на момент заказа
- Автоматическое резервирование товара при создании заказа
- Подтверждение и отмена заказов с обновлением остатков

## API Endpoints

### Пользователи
- `POST /api/v1/users` - Регистрация пользователя
- `GET /api/v1/users/{id}` - Получить пользователя
- `GET /api/v1/users/{user_id}/orders` - Заказы пользователя

### Товары
- `POST /api/v1/products` - Создать товар
- `GET /api/v1/products` - Список товаров
- `GET /api/v1/products/{id}` - Получить товар
- `PUT /api/v1/products/{id}/quantity` - Обновить количество

### Заказы
- `POST /api/v1/orders` - Создать заказ
- `GET /api/v1/orders/{id}` - Получить заказ
- `POST /api/v1/orders/{id}/confirm` - Подтвердить заказ
- `POST /api/v1/orders/{id}/cancel` - Отменить заказ

### Служебные
- `GET /health` - Проверка здоровья сервиса

## Запуск проекта

### Быстрый старт (Docker Compose)

```bash
# Клонируем репозиторий
git clone https://github.com/AndrivA89/orders.git
cd orders

# Запускаем все сервисы одной командой
make run
# или
docker compose up -d

# Ждем ~10 секунд и проверяем что сервис работает
curl http://localhost:8080/health
```

**Быстрый тест API:**
```bash
# Создать пользователя
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Тест","last_name":"Пользователь","age":25,"password":"test123"}'

# Создать товар  
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"description":"Тестовый товар","quantity":10,"price":1000}'
```

### Локальный запуск для разработки

```bash
# Устанавливаем зависимости
go mod download

# Запускаем только PostgreSQL
docker compose up -d postgres

# Запускаем сервис локально
make build && ./bin/orders
# или
go run cmd/server/main.go
```

### Доступные команды

```bash
make help              # Показать все команды
make run               # Запустить через Docker
make stop              # Остановить Docker сервисы  
make test              # Запустить тесты
make build             # Собрать приложение
make generate          # Генерировать моки
```

## Тестирование

### Запуск тестов
```bash
# Все тесты
go test ./...

# С подробным выводом
go test ./... -v

# С покрытием кода
go test -cover ./...

# Конкретный пакет
go test ./internal/domain/entities/ -v
```

### Генерация моков
```bash
# Генерация всех моков
make generate

# Ручная генерация
go generate ./internal/domain/repositories/
```

## Примеры использования

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe", 
    "age": 25,
    "is_married": false,
    "password": "password123"
  }'
```

### Создание товара

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "description": "iPhone 15 Pro",
    "quantity": 10,
    "price": 99999
  }'
```

### Создание заказа

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-uuid-here",
    "items": [
      {
        "product_id": "product-uuid-here", 
        "quantity": 2
      }
    ]
  }'
```

## Технологический стек

- **Go 1.24** - Основной язык
- **Gin** - HTTP framework с middleware
- **GORM** - ORM для работы с PostgreSQL  
- **PostgreSQL 15** - Реляционная база данных
- **Docker & Docker Compose** - Контейнеризация
- **OpenTelemetry** - Observability и трассировка
- **Testify + Gomock** - Unit тестирование
- **Logrus** - Структурированное логирование
- **bcrypt** - Хеширование паролей


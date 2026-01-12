# Руководство по работе с шаблоном FinOps Microservice

## Введение

Это шаблон микросервиса для обработки финансовых операций. Вам предоставлена готовая архитектура и каркас кода — ваша задача реализовать недостающие методы.

## Структура проекта

```
cmd/server/main.go              ← Точка входа (заполнить)
internal/
  api/handlers.go               ← HTTP обработчики (заполнить)
  config/config.go              ← Конфигурация (готово)
  db/db.go                      ← Подключение БД (заполнить)
  middleware/middleware.go       ← Логирование (заполнить)
  models/models.go              ← Структуры данных (готово)
  processor/processor.go        ← Асинхронная обработка (заполнить)
  repositories/                 ← Доступ к БД (заполнить)
  services/                     ← Бизнес-логика (заполнить)
migrations/init.sql             ← Схема БД (готово)
tests/                          ← Тесты (заполнить)
```

## Этап 1: REST API и база данных

### 1.1 Подготовка
```bash
go mod download
docker-compose up -d db
```

### 1.2 Реализация слоя доступа к данным

**internal/db/db.go** — функция `NewPool(dsn)`
- Создать пул подключений pgxpool
- Проверить подключение через Ping()
- Вернуть пул или ошибку

**internal/repositories/user_repo.go** — методы интерфейса
- `GetBalance(ctx, userID)` — SELECT баланс пользователя
- `UpdateBalance(ctx, userID, amount, tx)` — UPDATE баланс (используйте переданную транзакцию tx!)

**internal/repositories/transaction_repo.go** — методы интерфейса
- `CreateTransaction(ctx, tx, dbTx)` — INSERT новую транзакцию, вернуть ID
- `GetTransaction(ctx, id)` — SELECT транзакцию по ID
- `UpdateTransaction(ctx, id, newTx)` — UPDATE поля транзакции
- `DeleteTransaction(ctx, id)` — DELETE транзакцию

### 1.3 Реализация бизнес-логики

**internal/services/user_service.go**
- `GetBalance()` — делегировать в репозиторий

**internal/services/transaction_service.go** (КРИТИЧНО!)
- `CreateTransaction()` — самый важный метод:
  1. Валидировать amount > 0 и type ∈ {deposit, withdraw}
  2. Начать БД транзакцию (`pool.Begin()`)
  3. Убедиться, что пользователь существует
  4. Для withdraw проверить баланс пользователя
  5. Вызвать `txRepo.CreateTransaction()` и `userRepo.UpdateBalance()` с переданной транзакцией
  6. Коммитить транзакцию
  
- `GetTransaction()` — получение транзакции
- `UpdateTransaction()` — обновление (проверить processed == false)
- `DeleteTransaction()` — удаление с реверсом баланса если processed == true

### 1.4 Реализация HTTP API

**internal/api/handlers.go** — реализовать все обработчики:
- `CreateTransactionHandler()` — POST /transactions (201 Created)
- `GetTransactionHandler()` — GET /transactions/{id} (200 или 404)
- `UpdateTransactionHandler()` — PUT /transactions/{id} (204 No Content)
- `DeleteTransactionHandler()` — DELETE /transactions/{id} (204 No Content)
- `GetUserBalanceHandler()` — GET /users/{user_id}/balance (200 OK)
- `HealthCheckHandler()` — GET /health (200 OK)

### 1.5 Middleware и главная функция

**internal/middleware/middleware.go**
- `LoggingMiddleware()` — логировать метод, путь, время выполнения
- `RecoveryMiddleware()` — перехватить panic, логировать, вернуть 500

**cmd/server/main.go**
- Загрузить конфигурацию
- Создать пул БД
- Инициализировать репозитории, сервисы
- Создать HTTP маршруты
- Запустить сервер с обработкой сигналов завершения

### Проверка этапа 1
```bash
go run cmd/server/main.go
curl http://localhost:8080/health          # Должен вернуть 200
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"amount":"100","type":"deposit"}'
go test ./tests -v
```

## Этап 2: Конкурентная обработка

### 2.1 Реализация Processor

**internal/processor/processor.go**
- `NewProcessor(pool, numWorkers)` — создать процессор с worker горутинами
  - Инициализировать buffered канал jobs (размер 100)
  - Запустить numWorkers горутин через worker()
  
- `Submit(tx)` — отправить транзакцию в очередь обработки

- `worker()` (приватный) — цикл обработки из канала jobs

- `process(tx)` (приватный) — обработка одной транзакции:
  1. Получить или создать мьютекс для пользователя (защитить RWMutex)
  2. Заблокировать мьютекс пользователя
  3. Начать БД транзакцию
  4. Получить баланс, проверить его для withdraw
  5. Обновить баланс в БД
  6. Установить processed = true
  7. Коммитить транзакцию

- `Close()` — закрыть канал для graceful shutdown

### 2.2 Интеграция в приложение

В `cmd/server/main.go`:
- Создать Processor: `proc := processor.NewProcessor(pool, 5)`
- Передавать `proc` в обработчики
- Вызывать `proc.Submit()` после создания транзакции

### 2.3 Написание тестов

**tests/api_test.go, tests/concurrency_test.go, tests/integration_test.go**
- Unit тесты JSON сериализации и валидации
- Тесты конкурентности и race conditions
- Интеграционные тесты эндпоинтов

### Проверка этапа 2
```bash
go test ./tests -v
go test -race ./tests -v              # КРИТИЧНО: не должно быть data races!
```

## Этап 3: Оптимизация и контейнеризация

### 3.1 Создание Dockerfile

На основе `Dockerfile.template`:
- Build stage: собрать Go приложение
- Runtime stage: минимальный alpine образ

### 3.2 Настройка docker-compose

В `docker-compose.yml` раскомментировать и заполнить сервис `app`:
```yaml
app:
  build: .
  ports:
    - "8080:8080"
  environment:
    DATABASE_URL: postgres://user:pass@db:5432/finops?sslmode=disable
    PORT: 8080
  depends_on:
    db:
      condition: service_healthy
```

### 3.3 Профилирование (опционально)

Добавить pprof в main.go для анализа производительности.

### Проверка этапа 3
```bash
docker-compose up --build
curl http://localhost:8080/health
docker-compose down
```

## Критерии успешного выполнения

**ОР1: REST API и БД**
- ✅ Приложение запускается без ошибок
- ✅ Все CRUD операции работают
- ✅ Health check возвращает 200
- ✅ Валидация входных данных работает

**ОР2: Конкурентная обработка**
- ✅ `go test ./tests -v` проходит все тесты
- ✅ `go test -race ./tests -v` — ноль race conditions
- ✅ Транзакции обрабатываются асинхронно (processed флаг меняется)

**ОР3: Контейнеризация**
- ✅ `docker-compose up` запускает оба сервиса
- ✅ Приложение работает в контейнере
- ✅ API доступен на localhost:8080

## Полезные команды

```bash
# Запустить БД
docker-compose up -d db

# Запустить приложение
go run cmd/server/main.go

# Запустить тесты
go test ./tests -v

# Запустить с проверкой race conditions
go test -race ./tests -v

# Запустить приложение и БД
docker-compose up

# Остановить все
docker-compose down

# Просмотреть логи
docker-compose logs -f db
docker-compose logs -f app
```

## Важные замечания

1. **ACID транзакции** — используйте `pool.Begin()` для гарантии консистентности при параллельных операциях
2. **Parameterized queries** — всегда используйте `$1, $2` вместо конкатенации строк
3. **Синхронизация** — используйте мьютексы для защиты баланса каждого пользователя
4. **Тестирование** — `go test -race` обязателен для проверки гонок данных
5. **Context** — передавайте context.Context везде для отмены и таймаутов

## Где найти помощь

- Комментарии `// TODO` в коде указывают, что нужно реализовать
- Интерфейсы в `internal/repositories/interfaces.go` показывают контракт методов

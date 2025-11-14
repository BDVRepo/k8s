# gRPC Песочница

Учебный проект, показывающий, как два микросервиса на Go (`orders` и `payments`) обмениваются данными по gRPC и выстроены по принципам чистой архитектуры.

## Структура проекта

```
├── cmd/                # Точки входа сервисов
├── docs/               # Учебные материалы и заметки по архитектуре
├── internal/           # Слои приложения, домена и инфраструктуры
│   ├── models/         # Доменные модели (Order, Payment, Money)
│   ├── methods/handlers/{domain}/
│   │   ├── adapter/   # Репозитории (инфраструктурный слой)
│   │   ├── usecase/    # Бизнес-логика (application слой)
│   │   └── grpc/       # gRPC handlers (транспортный слой)
│   ├── infrastructure/ # Инфраструктурные адаптеры (gRPC клиенты)
│   └── gen/            # Сгенерированный код из proto
├── pkg/                # Общие вспомогательные пакеты (конфиги, генератор id и т.д.)
├── proto/              # Протобуф-файлы (контракты gRPC)
└── Makefile            # Утилиты для сборки, генерации и тестов
```

## Что такое gRPC и как это работает

### Основные понятия

**gRPC** — фреймворк для межсервисного взаимодействия:
- Использует **Protocol Buffers** (protobuf) для сериализации данных (бинарный формат, быстрее JSON)
- Работает поверх **HTTP/2** (поддержка стриминга, мультиплексирования)
- Типобезопасные контракты через `.proto` файлы
- Автоматическая генерация клиент/сервер кода

### Сущности gRPC

1. **`.proto` файлы** (`proto/`) — описание контрактов:
   - `message` — структуры данных (CreateOrderRequest, OrderSummary)
   - `service` — интерфейсы RPC-методов (CreateOrder, GetOrder)

2. **Сгенерированный код** (`internal/gen/`) — создается командой `make proto`:
   - `*_pb.go` — Go структуры для сообщений
   - `*_grpc.pb.go` — интерфейсы серверов/клиентов

3. **gRPC Server** — реализует интерфейс из `*_grpc.pb.go`, обрабатывает RPC-вызовы

4. **gRPC Client** — вызывает методы удаленного сервиса по сети

### Архитектура проекта (Clean Architecture)

**Поток запроса при создании заказа:**

```
Клиент (grpcurl)
  ↓ gRPC вызов
gRPC Handler (internal/methods/handlers/orders/grpc/)
  ↓ преобразует protobuf → domain command
Use Case (internal/methods/handlers/orders/usecase/)
  ↓ создает Order, вызывает payments.Authorize()
gRPC Client (internal/infrastructure/grpc/)
  ↓ gRPC вызов к payments service
Payments Service
  ↓ обрабатывает, возвращает paymentID
Use Case (продолжение)
  ↓ сохраняет Order через Repository
Repository (internal/methods/handlers/orders/adapter/)
  ↓ сохраняет в память/БД
gRPC Handler (возврат)
  ↓ преобразует domain → protobuf
Клиент получает ответ
```

**Слои архитектуры:**

- **Transport (grpc/)** — принимает gRPC запросы, валидирует, преобразует protobuf ↔ domain
- **Application (usecase/)** — бизнес-логика, оркестрация, зависит только от интерфейсов
- **Domain (models/)** — бизнес-сущности (Order, Payment), методы домена
- **Infrastructure (adapter/, infrastructure/)** — реализации репозиториев, gRPC клиенты

**Преимущества:**
- Тестируемость: use case тестируется через моки интерфейсов
- Независимость: domain не зависит от транспорта/БД
- Расширяемость: можно добавить REST/GraphQL без изменения use case

### Откуда появляются сущности

1. **Proto контракты** → пишутся вручную в `proto/`
2. **Сгенерированный код** → `make proto` создает `internal/gen/`
3. **Domain модели** → пишутся вручную в `internal/models/`
4. **Use cases** → пишутся вручную в `internal/methods/handlers/*/usecase/`
5. **Repositories** → пишутся вручную в `internal/methods/handlers/*/adapter/`
6. **gRPC handlers** → пишутся вручную в `internal/methods/handlers/*/grpc/`
7. **Infrastructure clients** → пишутся вручную в `internal/infrastructure/`

## Требования

- Go 1.24+ (Go toolchain сам подтянет нужную версию)
- `protoc` v3.21+ с плагинами `protoc-gen-go` и `protoc-gen-go-grpc`
- Docker (по желанию, для контейнерного сценария)

Установка плагинов Go:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Генерация gRPC-кода

```bash
make proto
```

Команда запускает `protoc`, обновляет сгенерированные заглушки в `internal/gen/orderspb` и `internal/gen/paymentspb`, затем очищает временные директории.

## Локальный запуск

Через терминал:

```bash
make run-payments   # запускает сервис payments на порту :50052
make run-orders     # запускает сервис orders на порту :50051 (делает вызовы в payments)
```

Через Visual Studio Code:

1. Откройте рабочую папку `grpc`
2. Перейдите на вкладку Run and Debug
3. Выберите составную конфигурацию `Payments + Orders`, чтобы поднять оба сервиса одновременно

## Пример взаимодействия

Проверьте через [`grpcurl`](https://github.com/fullstorydev/grpcurl):

```bash
grpcurl -plaintext localhost:50051 orders.v1.OrdersService/CreateOrder \
  -d '{"user_id":"u1","product_id":"p1","price":{"currency":"USD","amount":42}}'
```

```bash
grpcurl -plaintext localhost:50051 orders.v1.OrdersService/GetOrder \
  -d '{"order_id":"<value from create response>"}'
```

## Тестирование

### 1. Автотесты Go

```bash
make test
```

- Прогоняет `go test ./...`, гарантируя, что проект собирается и не содержит регрессий.  
- Рекомендуется запускать после любых изменений в коде или протобуфах.  
- При необходимости можно добавить собственные unit-тесты в `internal/methods/handlers/*/usecase/` или `internal/models/...`.

### 2. Интерактивные gRPC-запросы

Используйте `grpcurl` из корневой машины (или внутри контейнера):

```bash
# Создание заказа (orders → payments)
grpcurl -plaintext localhost:50051 orders.v1.OrdersService/CreateOrder \
  -d '{"user_id":"u1","product_id":"p1","price":{"currency":"USD","amount":42}}'

# Получение заказа по ID
grpcurl -plaintext localhost:50051 orders.v1.OrdersService/GetOrder \
  -d '{"order_id":"<ID из ответа CreateOrder>"}'

# Проверка статуса платежа напрямую
grpcurl -plaintext localhost:50052 payments.v1.PaymentsService/GetPaymentStatus \
  -d '{"payment_id":"<ID из ответа CreateOrder>"}'
```

### 3. Логи и наблюдаемость

- Просмотр потоков логов:

  ```bash
  docker-compose logs -f orders
  docker-compose logs -f payments
  ```

  Позволяет убедиться, что цепочка `orders → payments` отрабатывает корректно.

- При локальном запуске через `make run-*` логи выводятся прямо в терминал.

### 4. Расширенные сценарии

- Добавьте gRPC-интерсепторы для логирования/метрик и проверяйте их работу через те же команды.  
- Реализуйте интеграционные тесты, поднимая оба сервиса в тестовой среде (например, `docker compose up` внутри CI) и выполняя `grpcurl`/e2e-скрипты автоматически.  
- При переходе на стриминговые RPC дополните документацию примерами `grpcurl` или клиентскими скриптами.

## Docker

Сборка и запуск через Compose:

```bash
docker-compose up --build
```

В результате будут доступны:

- Сервис orders на `localhost:50051`
- Сервис payments на `localhost:50052`

## Дальнейшие шаги

- Подключить постоянное хранилище, реализовав новые репозитории
- Добавить interceptors для наблюдаемости и аутентификации
- Освоить стриминговые RPC по рекомендациям из `docs/grpc-learning.md`

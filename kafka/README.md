# Kafka Learning Project

Учебный проект для изучения Kafka, event-driven архитектуры и основных паттернов с рабочими примерами на Go.

## Содержание

- [Быстрый старт](#быстрый-старт)
- [Основные понятия Kafka](#основные-понятия-kafka)
- [Паттерны Event-Driven](#паттерны-event-driven)
  - [Event Sourcing](#event-sourcing)
  - [CQRS](#cqrs)
  - [Saga Pattern](#saga-pattern)
  - [Outbox Pattern](#outbox-pattern)
- [Интеграция с gRPC](#интеграция-с-grpc)
- [Примеры использования](#примеры-использования)

## Быстрый старт

### Запуск инфраструктуры

```bash
# Запустить Kafka, Zookeeper и Kafka UI
make up

# Kafka UI доступен на http://localhost:8090
```

### Базовые примеры

```bash
# В одном терминале - запустить consumer
go run ./cmd/consumer

# В другом терминале - запустить producer
go run ./cmd/producer
```

## Основные понятия Kafka

### Топики (Topics)

Топик - это категория/канал для сообщений. Производители (producers) публикуют сообщения в топики, потребители (consumers) читают из топиков.

### Партиции (Partitions)

Топик разделен на партиции для параллельной обработки. Сообщения с одинаковым ключом попадают в одну партицию (гарантия порядка).

### Consumer Groups

Группа потребителей распределяет нагрузку между собой. Каждое сообщение обрабатывается только одним consumer в группе.

### Offset

Позиция чтения в партиции. Consumer отслеживает свой offset для каждого топика/партиции.

## Паттерны Event-Driven

### Event Sourcing

**Суть паттерна:** Храним все события, которые произошли с сущностью, вместо хранения текущего состояния.

**Преимущества:**
- Полная история изменений
- Возможность восстановления состояния на любой момент времени
- Аудит и отладка
- Возможность "переиграть" события

**Когда использовать:**
- Когда важна полная история изменений
- Когда нужно восстановление состояния
- Для аудита и compliance

**Пример:**

```go
// Сохраняем событие
eventStore.SaveEvent(ctx, OrderCreatedEvent{...})
eventStore.SaveEvent(ctx, OrderPaidEvent{...})

// Восстанавливаем состояние
state := eventStore.RebuildState(ctx, orderID)
```

**Как хранится история:** в примере `InMemoryEventStore` использует `map[string][]Event`, где ключ — `order_id`, а значение — упорядоченный список событий. Каждое новое событие аппендится в этот slice и одновременно публикуется в Kafka. `RebuildState(orderID)` просто итерируется по всему списку и применяет события одно за другим (Created → Paid → Shipped), поэтому можно восстановить любое прошлое состояние или посмотреть всю историю заказа.

**Запуск примера:**

```bash
go run ./cmd/orders-producer
```

**Поток данных:**

```
CreateOrder → SaveEvent(OrderCreated) → Kafka → RebuildState
ProcessPayment → SaveEvent(OrderPaid) → Kafka → RebuildState
ShipOrder → SaveEvent(OrderShipped) → Kafka → RebuildState
```

### CQRS

**Суть паттерна:** Разделение операций чтения (Read) и записи (Write) на отдельные модели.

**Write Side (cmd/cqrs-write):**
- Принимает команды `CreateOrder/ProcessPayment`
- Валидирует, генерирует ID и публикует события `OrderCreated` (topic `orders`) и `OrderPaid` (topic `payments`)
- Не хранит read-модель, только фиксирует факт, что «что-то произошло»

**Read Side (cmd/analytics):**
- Подписывается на `orders` и `payments`
- Строит собственную read model `AnalyticsReadModel`, в которой:
  - `map[userID]count`, `map[productID]count` для быстрых агрегаций
  - `recentOrders` (ring buffer) для последних 100 заказов
- Периодически (каждые 10 секунд) показывает `GetStats()` и `GetRecentOrders()`

**Преимущества:**
- Масштабируемость (read и write масштабируются независимо)
- Оптимизация для разных use cases
- Разделение ответственности

**Когда использовать:**
- Когда read и write нагрузки сильно различаются
- Когда нужны разные представления данных
- Для аналитики и отчетности

**Пример:**

```go
// Write Side - публикует события
writeSide.CreateOrder(ctx, CreateOrderCommand{...})

// Read Side - строит read model из событий
readModel.HandleOrderCreated(ctx, event)
stats := readModel.GetStats()
```

**Запуск примера:**

```bash
# Terminal 1: Write Side
go run ./cmd/cqrs-write

# Terminal 2: Read Side (Analytics)
go run ./cmd/analytics
```

**Поток данных:**

```
Write Side:
  CreateOrder → Publish(OrderCreated) → Kafka
  ProcessPayment → Publish(OrderPaid) → Kafka

Read Side:
  Subscribe(Kafka orders) → HandleOrderCreated → increment totals/recent orders
  Subscribe(Kafka payments) → HandleOrderPaid → обновить статус заказа → GetStats()
```

### Saga Pattern

**Суть паттерна:** Оркестратор управляет длинным бизнес-процессом («шаги») и фиксирует, какие локальные транзакции прошли. Если на шаге N происходит ошибка, он запускает компенсирующие действия для шагов 0…N‑1 (rollback).

**Типы Saga:**
- **Choreography:** Каждый сервис знает, что делать дальше (децентрализовано)
- **Orchestration:** Центральный оркестратор управляет шагами (централизовано)

**Компенсация:**
При ошибке на любом шаге выполняются компенсирующие действия для отката предыдущих шагов.

**Преимущества:**
- Работа с распределенными транзакциями
- Нет необходимости в 2PC (Two-Phase Commit)
- Масштабируемость

**Когда использовать:**
- Длинные бизнес-процессы
- Распределенные транзакции между сервисами
- Когда нужна компенсация при ошибках

**Пример (cmd/saga-orchestrator):**

```go
orchestrator := saga.NewSagaOrchestrator(producer, sagaID, orderID)
orchestrator.AddStep("CreateOrder", "order-commands", "order-compensate")   // Step #1, с командным топиком и топиком компенсации
orchestrator.AddStep("ProcessPayment", "payment-commands", "payment-compensate")
orchestrator.AddStep("ShipOrder", "shipping-commands", "shipping-compensate")

// Успешное выполнение
orchestrator.Execute(ctx)

// С ошибкой и компенсацией
orchestrator.ExecuteWithFailure(ctx, failAtStep)
```

**Запуск примера:**

```bash
go run ./cmd/saga-orchestrator
```

**Поток данных:**

```
Orchestrator:
  Step1(CreateOrder cmd) → публикуем команду в topic order-commands → сервис исполняет → оркестратор получает SagaStepCompleted → продолжает
  Step2(ProcessPayment) → ...
  Step3(ShipOrder) → ...
  После финального шага публикуется SagaCompleted event

При ошибке:
  Step1(CreateOrder) → Success
  Step2(ProcessPayment) → Error → оркестратор запускает compensate-цепочку
  Compensate(Step1) → публикуем сообщение в topic order-compensate → сервис откатывает действие
```

### Outbox Pattern

**Суть паттерна:** Гарантия доставки событий из БД транзакций через промежуточную таблицу (outbox).

**Проблема:**
- Нужно сохранить данные в БД и опубликовать событие в Kafka
- Если Kafka недоступен, что делать?
- Как гарантировать exactly-once delivery?

**Решение:**
1. Сохраняем событие в outbox таблицу в той же БД транзакции
2. Отдельный процесс (publisher) читает из outbox и публикует в Kafka
3. После успешной публикации помечаем событие как опубликованное

В нашем примере (`internal/patterns/outbox`):
- `TransactionalService.CreateOrderWithOutbox` моделирует транзакцию: заказ + запись в `InMemoryOutboxStore` (map[eventID]OutboxEvent) происходят под одним контекстом.
- `OutboxPublisher.Start` каждые N секунд читает `GetUnpublishedEvents`, публикует через общий producer и вызывает `MarkAsPublished`.
- Если Kafka недоступна, событие остаётся в store и будет доставлено при следующей попытке.

**Преимущества:**
- Гарантия доставки (событие не потеряется)
- Транзакционная согласованность
- Exactly-once delivery

**Когда использовать:**
- Когда критична гарантия доставки событий
- При работе с БД транзакциями
- Для критичных бизнес-событий

**Пример:**

```go
// В транзакции БД
BEGIN TRANSACTION
  INSERT INTO orders ...
  INSERT INTO outbox (topic, key, value) ...
COMMIT

// Отдельный процесс публикует из outbox
publisher.Start(ctx) // читает из outbox и публикует в Kafka
```

**Запуск примера:**

```bash
go run ./cmd/outbox-example
```

**Поток данных:**

```
Application:
  CreateOrder → Save to DB + Save to Outbox (в одной транзакции) → Commit

Publisher (background):
  Poll Outbox → Get Unpublished Events → Publish to Kafka → Mark as Published
  (при неудаче событие остаётся в outbox до успешной публикации)
```

## Интеграция с gRPC

Проект интегрирован с существующим gRPC проектом (`/home/bdv/projects/k8s/grpc`).

### Как это работает

1. **gRPC Orders Service** создает заказ через gRPC
2. После успешного создания заказ публикует событие `OrderCreated` в Kafka
3. **Kafka Consumer** подписывается на события и обрабатывает их

### Запуск интеграции

```bash
# Terminal 1: Запустить Kafka
cd /home/bdv/projects/k8s/kafka
make up

# Terminal 2: Запустить gRPC Orders Service
cd /home/bdv/projects/k8s/grpc
go run ./cmd/orders

# Terminal 3: Запустить Kafka Consumer для событий от gRPC
cd /home/bdv/projects/k8s/kafka
go run ./cmd/grpc-events-consumer

# Terminal 4: Создать заказ через gRPC
grpcurl -plaintext -d '{"user_id":"user-123","product_id":"prod-456","price":{"currency":"USD","amount":99.99}}' \
  localhost:50051 orders.v1.OrdersService/CreateOrder
```

### Архитектура

```
gRPC Client
    ↓
gRPC Orders Service
    ↓ CreateOrder
Use Case (CreateOrder)
    ↓ Save Order
Repository
    ↓ Publish Event
Kafka Producer
    ↓
Kafka Topic (orders)
    ↓
Kafka Consumer (grpc-events-consumer)
    ↓ Process Event
Analytics / Notifications / Search Index
```

## Примеры использования

### 1. Базовый Producer/Consumer

```bash
# Consumer
go run ./cmd/consumer

# Producer
go run ./cmd/producer
```

### 2. Event Sourcing

```bash
go run ./cmd/orders-producer
```

Показывает:
- Сохранение событий
- Восстановление состояния из событий
- Публикацию событий в Kafka

### 3. CQRS

```bash
# Write Side
go run ./cmd/cqrs-write

# Read Side (Analytics)
go run ./cmd/analytics
```

### 4. Saga Pattern

```bash
go run ./cmd/saga-orchestrator
```

Показывает:
- Успешное выполнение Saga
- Обработку ошибок и компенсацию

### 5. Outbox Pattern

```bash
go run ./cmd/outbox-example
```

Показывает:
- Сохранение событий в outbox
- Автоматическую публикацию из outbox в Kafka

## Структура проекта

```
kafka/
├── cmd/
│   ├── producer/              # Базовый producer
│   ├── consumer/              # Базовый consumer
│   ├── orders-producer/       # Event Sourcing пример
│   ├── analytics/             # CQRS Read Side
│   ├── cqrs-write/            # CQRS Write Side
│   ├── saga-orchestrator/     # Saga Pattern
│   ├── outbox-example/        # Outbox Pattern
│   └── grpc-events-consumer/  # Consumer для gRPC событий
├── internal/
│   ├── models/                # Event модели
│   ├── kafka/                 # Kafka клиенты
│   ├── patterns/
│   │   ├── eventsourcing/     # Event Sourcing
│   │   ├── cqrs/              # CQRS
│   │   ├── saga/              # Saga Pattern
│   │   └── outbox/            # Outbox Pattern
│   └── handlers/              # Обработчики событий
├── pkg/
│   └── config/                # Конфигурация
├── docker-compose.yml         # Kafka инфраструктура
├── Makefile                   # Команды управления
└── README.md                  # Документация
```

## Сравнение паттернов

| Паттерн | Когда использовать | Преимущества | Недостатки |
|---------|-------------------|--------------|------------|
| **Event Sourcing** | Нужна полная история, аудит | История изменений, восстановление состояния | Сложность, размер хранилища |
| **CQRS** | Разные read/write нагрузки | Масштабируемость, оптимизация | Сложность синхронизации |
| **Saga** | Распределенные транзакции | Нет 2PC, масштабируемость | Сложность компенсации |
| **Outbox** | Критична гарантия доставки | Exactly-once, транзакционность | Дополнительная таблица |

## Полезные команды

```bash
# Управление инфраструктурой
make up          # Запустить Kafka
make down        # Остановить Kafka
make logs        # Логи Kafka
make kafka-ui    # Открыть Kafka UI (http://localhost:8090)

# Разработка
go mod tidy      # Обновить зависимости
go test ./...    # Запустить тесты
```

## Kafka UI

После запуска `make up`, Kafka UI доступен на http://localhost:8090

В UI можно:
- Просматривать топики и сообщения
- Мониторить consumer groups
- Просматривать партиции и offsets
- Публиковать тестовые сообщения

## Дальнейшее изучение

1. **Kafka Streams** - обработка потоков данных
2. **Kafka Connect** - интеграция с внешними системами
3. **Schema Registry** - управление схемами сообщений
4. **Exactly-once semantics** - гарантии доставки
5. **Kafka в production** - настройка, мониторинг, тюнинг

## Ресурсы

- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [kafka-go Library](https://github.com/segmentio/kafka-go)
- [Event-Driven Architecture Patterns](https://martinfowler.com/articles/201701-event-driven.html)


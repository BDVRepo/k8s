# Docker — Теория и практика

> Аналогия: Docker — это стандартный **грузовой контейнер** из логистики.
> Неважно, что внутри — он одинаково грузится на любой корабль, поезд или грузовик.
> Так и Docker-контейнер: работает одинаково на ноутбуке, CI-сервере и в production.

---

## 1. Базовые команды

```bash
# Сборка
docker build . -t myapp:1.0.0          # собрать образ с тегом
docker build . -t myapp:1.0.0 --no-cache  # без кеша

# Образы
docker images                           # список образов
docker rmi <id>                         # удалить образ
docker tag myapp:1.0.0 user/myapp:1.0.0
docker push user/myapp:1.0.0
docker login

# Запуск
docker run myapp:1.0.0                  # запустить контейнер
docker run -d myapp:1.0.0              # в фоне (detached)
docker run -p 8080:8080 myapp:1.0.0    # проброс порта HOST:CONTAINER
docker run --network host myapp:1.0.0  # использовать сеть хоста
docker run -v /host/path:/container/path myapp:1.0.0  # bind mount

# Контейнеры
docker ps                               # запущенные контейнеры
docker ps -a                            # все (включая остановленные)
docker logs <id>                        # логи
docker logs <id> -f                     # логи в реальном времени
docker stop <id>
docker exec -it <id> bash               # войти в контейнер

# Очистка
docker system prune                     # удалить всё неиспользуемое
docker system prune -a                  # + все неиспользуемые образы (агрессивнее)
docker system prune --volumes           # + тома (осторожно: потеря данных!)
docker image ls && docker history <image>  # проверить размер и слои
```

---

## 2. Dockerfile — ключевые инструкции

> Аналогия: Dockerfile — это **рецепт** для блюда.
> Каждый шаг (инструкция) создаёт новый слой, как слои торта.
> Если ингредиент не изменился — Docker берёт его из кеша.

| Инструкция | Синтаксис | Откуда → Куда / Что делает |
|---|---|---|
| `FROM` | `FROM <image>:<tag>` | Базовый образ. Пример: `FROM node:20-alpine` |
| `COPY` | `COPY <src> <dest>` | С хоста → в образ. Пример: `COPY . /app` |
| `ADD` | `ADD <src> <dest>` | Как COPY, но умеет распаковывать архивы и скачивать URL. Пример: `ADD app.tar.gz /app/` |
| `RUN` | `RUN <cmd> <arg1>` | Выполнить при сборке (создаёт слой). Пример: `RUN npm ci --only=production` |
| `CMD` | `CMD ["<exe>", "<arg>"]` | Команда по умолчанию при старте. Пример: `CMD ["node", "server.js"]` |
| `ENTRYPOINT` | `ENTRYPOINT ["<exe>"]` | Фиксированный исполняемый файл. Пример: `ENTRYPOINT ["python", "main.py"]` |
| `EXPOSE` | `EXPOSE <port>` | Документирует порт (не открывает! нужен `-p`). Пример: `EXPOSE 8080` |
| `ENV` | `ENV <key>=<value>` | Переменная окружения в образе. Пример: `ENV NODE_ENV=production` |
| `ARG` | `ARG <name>=<default>` | Аргумент только при сборке. Пример: `ARG APP_VERSION=1.0.0` |
| `WORKDIR` | `WORKDIR <path>` | Рабочая директория внутри образа. Пример: `WORKDIR /app` |
| `USER` | `USER <user\|uid>` | Под каким пользователем запускать. Пример: `USER node` |

---

## 3. CMD vs ENTRYPOINT

> Аналогия:
> - `ENTRYPOINT` — это **профессия** человека (всегда программист).
> - `CMD` — это **задача на сегодня** (по умолчанию "пиши код", но можно поменять на "пиши тесты").

```dockerfile
# Вариант 1: только CMD (легко заменить целиком)
CMD ["./app", "--port=8080"]
# docker run img                  -> ./app --port=8080
# docker run img echo hello       -> echo hello  (CMD заменён полностью)

# Вариант 2: ENTRYPOINT + CMD (частый production-паттерн)
ENTRYPOINT ["./app"]
CMD ["--port=8080"]
# docker run img                  -> ./app --port=8080
# docker run img --port=9090      -> ./app --port=9090  (заменили только аргумент)
# docker run --entrypoint sh img  -> полностью сменили ENTRYPOINT
```

### Сравнение CMD vs ENTRYPOINT

| | CMD | ENTRYPOINT |
|---|---|---|
| Можно заменить через `docker run` | Да, полностью | Только через `--entrypoint` |
| Используется как | Команда по умолчанию | Фиксированный бинарник |
| В связке | Полная команда | Получает CMD как аргументы |
| Рекомендация | Для гибких команд | Для "образ = одно приложение" |

### Shell-форма vs Exec-форма

| | Shell-форма | Exec-форма |
|---|---|---|
| Синтаксис | `ENTRYPOINT ./app` | `ENTRYPOINT ["./app"]` |
| PID 1 | `/bin/sh` (не получает сигналы) | само приложение ✅ |
| Graceful shutdown | Не работает | Работает |
| Рекомендация | Избегать | **Всегда использовать** |

---

## 4. Multi-stage builds

> Аналогия: Мы строим дом — нужны рабочие (компилятор) и стройматериалы (исходники).
> Но жить в готовом доме будет только семья (бинарник).
> Рабочих и строительный мусор в финальный дом не тащим.

```dockerfile
# Стадия 1: сборка (большой образ)
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

# Стадия 2: финальный образ (маленький)
FROM alpine:3.18
COPY --from=builder /app/app .
CMD ["./app"]
```

### С полностью пустой базой `scratch` — минимальный размер

```dockerfile
# Стадия 1: статический бинарник
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

# Стадия 2: пустой образ (нет вообще ничего)
FROM scratch
COPY --from=builder /app /app
ENTRYPOINT ["/app"]
```

### Облегченные базовые образы

| Образ | Размер | Плюсы | Минусы / Особенности |
|---|---|---|---|
| `ubuntu` / `debian` | ~70-100 MB | Полный Linux, всё есть | Большой, много лишнего |
| `alpine` | ~5 MB | Маленький, удобный | musl vs glibc — иногда несовместимость |
| `distroless` | ~2-20 MB | Почти без лишнего, безопасен | Нет shell для отладки |
| `scratch` | 0 MB | Абсолютный минимум | Только статический бинарник |

**Важно про `scratch`:**
- Нет `sh`, `bash`, `apk`, `apt`, сертификатов.
- Нужен статический бинарник (`CGO_ENABLED=0`), иначе приложение не стартует.
- Для HTTPS нужно дополнительно копировать CA-сертификаты.

---

## 5. Практичные фишки для уменьшения образа

- Используйте `.dockerignore` — чтобы не копировать `node_modules`, `.git`, тесты.
- Держите кэш-зависимости отдельно: сначала `COPY package*.json .`, потом `RUN npm ci`, потом `COPY . .`.
- Только production-зависимости: `npm ci --omit=dev`, `pip install --no-cache-dir`.
- Объединяйте установку/очистку в один `RUN`, чтобы не плодить слои:
```dockerfile
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
```
- Не храните секреты в образе — передавайте через env/secret на этапе запуска.
- Проверяйте размер: `docker image ls` и `docker history <image>`.

---

## 6. Сети Docker

> Аналогия: сети — это как **коридоры в офисном здании**.
> - `bridge` — коридор на этаже (контейнеры одного здания общаются).
> - `host` — ты сидишь не в офисе, а прямо на улице (хост-сеть, нет изоляции).
> - `none` — комната без дверей (нет сетевого доступа).

| Тип | Описание | Когда использовать |
|---|---|---|
| `bridge` | По умолчанию. Изолированная виртуальная сеть | Локальная разработка |
| `host` | Контейнер использует сетевой стек хоста напрямую | Когда важна производительность сети |
| `none` | Нет сетевого доступа | Максимальная изоляция |
| custom bridge | `docker network create mynet` — контейнеры резолвят друг друга по имени | docker-compose, связанные контейнеры |

```bash
docker network create mynet
docker run --network mynet --name db postgres
docker run --network mynet --name app myapp   # app может обратиться к db по имени "db"
```

---

## 7. Volumes vs Bind Mounts

> Аналогия:
> - **Volume** — это сейф в банке (Docker управляет, надёжно, переносимо).
> - **Bind Mount** — это папка у тебя на столе (ты сам указываешь путь, быстро и удобно для разработки).

| | Volume | Bind Mount |
|---|---|---|
| Управление | Docker (`docker volume create`) | Путь на хосте |
| Переносимость | Высокая | Зависит от структуры хоста |
| Использование | Данные БД, постоянные данные | Разработка, монтирование конфигов |
| Синтаксис | `-v myvolume:/data` | `-v /host/path:/data` |
| Где хранится | `/var/lib/docker/volumes/` | Указанный путь хоста |

---

## 8. docker-compose vs Kubernetes

| | docker-compose | Kubernetes |
|---|---|---|
| Где запускать | Один хост (локально) | Много нод (кластер) |
| Масштабирование | Ограничено | Горизонтальное, автоматическое |
| Самовосстановление | Нет (при падении хоста) | Да |
| Сложность | Простой | Сложнее, но мощнее |
| Типичное применение | Local dev, тесты | Production |

# Теория для собеседования: Docker + Kubernetes (Junior)

> Конспект по материалам курса. Покрывает все темы треков 02_Docker, 03_Kubernetes, 05_Lifecycle, 07_Lifehacks, 08_CICD.

---

## 1. Docker

### 1.1 Базовые команды

```bash
docker build . -t myapp:1.0.0          # собрать образ с тегом
docker build . -t myapp:1.0.0 --no-cache
docker images                           # список образов
docker run myapp:1.0.0                  # запустить контейнер
docker run -d myapp:1.0.0              # в фоне (detached)
docker run -p 8080:8080 myapp:1.0.0    # проброс порта host:container
docker run --network host myapp:1.0.0  # использовать сеть хоста
docker run -v /host/path:/container/path myapp:1.0.0  # bind mount
docker ps                               # запущенные контейнеры
docker ps -a                            # все контейнеры
docker logs <id>                        # логи
docker logs <id> -f                     # логи в реальном времени
docker stop <id>
docker rmi <id>                         # удалить образ
docker system prune                     # удалить всё неиспользуемое
docker tag myapp:1.0.0 user/myapp:1.0.0
docker push user/myapp:1.0.0
docker login
```

### 1.2 Dockerfile — ключевые инструкции (синтаксис: откуда -> куда)

| Инструкция | Синтаксис (форма) | Откуда -> Куда / Что означает |
|---|---|---|
| `FROM` | `FROM <image>:<tag>` | Базовый образ, например: `FROM node:20-alpine` |
| `COPY` | `COPY <src> <dest>` | `<src>` с хоста -> `<dest>` в образе, пример: `COPY . /app` |
| `ADD` | `ADD <src> <dest>` | `<src>` (локальный файл/архив/URL) -> `<dest>` в образе, пример: `ADD app.tar.gz /app/` |
| `RUN` | `RUN <command> <arg1> <arg2>` | Выполнить команду при сборке, пример: `RUN npm ci --only=production` |
| `CMD` | `CMD ["<exe>", "<arg1>", "<arg2>"]` | Команда по умолчанию при старте, пример: `CMD ["node", "server.js"]` |
| `ENTRYPOINT` | `ENTRYPOINT ["<exe>", "<arg1>"]` | Фиксированный исполняемый файл, пример: `ENTRYPOINT ["python", "main.py"]` |
| `EXPOSE` | `EXPOSE <port>` | Документирует порт в контейнере, пример: `EXPOSE 8080` |
| `ENV` | `ENV <key>=<value>` | Задать переменную окружения, пример: `ENV NODE_ENV=production` |
| `ARG` | `ARG <name>=<default>` | Аргумент только для сборки, пример: `ARG APP_VERSION=1.0.0` |
| `WORKDIR` | `WORKDIR <path>` | Рабочая директория внутри образа, пример: `WORKDIR /app` |
| `USER` | `USER <username|uid>` | Под каким пользователем запускать процесс, пример: `USER node` |

**CMD vs ENTRYPOINT:**
```dockerfile
# Вариант 1: только CMD (легко заменить целиком)
CMD ["./app", "--port=8080"]
# docker run img                  -> ./app --port=8080
# docker run img echo hello       -> echo hello  (CMD заменен полностью)

# Вариант 2: ENTRYPOINT + CMD (частый production-паттерн)
ENTRYPOINT ["./app"]
CMD ["--port=8080"]
# docker run img                  -> ./app --port=8080
# docker run img --port=9090      -> ./app --port=9090  (заменили только CMD-аргумент)
# docker run --entrypoint sh img  -> полностью сменили ENTRYPOINT
```

Когда что выбирать:
- `CMD` — если нужно дать команду по умолчанию, которую легко заменить.
- `ENTRYPOINT` — если контейнер должен всегда запускать конкретный бинарник.
- Рекомендация: используйте exec-форму (`["cmd","arg"]`), а не shell-форму (`cmd arg`).

### 1.3 Multi-stage builds

Зачем: уменьшить размер финального образа (не тащить компилятор в production).

```dockerfile
# Стадия 1: сборка
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

# Стадия 2: финальный образ
FROM alpine:3.18
COPY --from=builder /app/app .
CMD ["./app"]
```

Финальный образ содержит только бинарник, без Go-тулчейна.

Пример с полностью пустой базой (`scratch`) — минимальный размер:

```dockerfile
# Стадия 1: сборка статического бинарника
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

# Стадия 2: пустой образ (вообще без shell и пакетов)
FROM scratch
COPY --from=builder /app /app
ENTRYPOINT ["/app"]
```

Важно про `scratch`:
- В `scratch` нет `sh`, `bash`, `apk`, `apt`, сертификатов и прочих утилит.
- Нужен статический бинарник (часто `CGO_ENABLED=0`), иначе приложение может не стартовать.
- Для HTTPS иногда нужно дополнительно копировать CA-сертификаты из builder-стадии.

Облегченные базовые образы (что чаще использовать):
- `alpine` — маленький и удобный, но иногда бывают нюансы совместимости (musl vs glibc).
- `distroless` — почти без лишнего софта, хороший баланс для production-безопасности.
- `scratch` — самый минимальный, но подходит не всем приложениям.

Практичные фишки для уменьшения образа:
- Используйте `.dockerignore`, чтобы не копировать `node_modules`, `.git`, тесты, артефакты.
- Держите кэш-зависимости отдельно: сначала `COPY package*.json .`, потом `RUN npm ci`, потом `COPY . .`.
- Ставьте только production-зависимости (`npm ci --omit=dev`, `pip install --no-cache-dir`).
- Объединяйте команды установки/очистки в один `RUN`, чтобы не плодить лишние слои.
- Не храните секреты в образе: передавайте через env/secret на этапе запуска.
- Регулярно проверяйте размер: `docker image ls` и `docker history <image>`.

### 1.4 Сети Docker

| Тип     | Описание |
|---------|----------|
| `bridge`| По умолчанию. Контейнеры изолированы, общаются через виртуальный bridge |
| `host`  | Контейнер использует сетевой стек хоста напрямую (нет изоляции) |
| `none`  | Нет сетевого доступа |
| custom  | `docker network create mynet` — именованная bridge-сеть, контейнеры резолвят друг друга по имени |

### 1.5 Volumes vs Bind Mounts

| | Volume | Bind Mount |
|---|--------|-----------|
| Управление | Docker (`docker volume create`) | Путь на хосте |
| Переносимость | Высокая | Зависит от хоста |
| Использование | Данные БД, постоянные данные | Разработка, конфиги |
| Синтаксис | `-v myvolume:/data` | `-v /host/path:/data` |

---

## 2. Архитектура Kubernetes

### 2.1 Компоненты Control Plane (Master)

| Компонент | Роль |
|-----------|------|
| `kube-apiserver` | Единая точка входа. Принимает запросы kubectl/REST, валидирует и записывает в etcd |
| `etcd` | Распределённое key-value хранилище. Хранит всё состояние кластера |
| `kube-scheduler` | Выбирает, на какую Node поставить Pod (смотрит на ресурсы, affinity, taints) |
| `kube-controller-manager` | Запускает контроллеры (Deployment, ReplicaSet, Node, Job...) в бесконечном цикле reconcile |

### 2.2 Компоненты Worker Node

| Компонент | Роль |
|-----------|------|
| `kubelet` | Агент на каждой ноде. Получает PodSpec, запускает/останавливает контейнеры через CRI |
| `kube-proxy` | Управляет сетевыми правилами (iptables/ipvs) для работы Services |
| Container Runtime | Запускает контейнеры. Обычно `containerd` или `CRI-O` |

### 2.3 Как работает `kubectl apply`

```
kubectl apply -f pod.yaml
        ↓
kube-apiserver (валидация, авторизация RBAC)
        ↓
etcd (сохранение desired state)
        ↓
kube-scheduler (выбирает Node для Pod)
        ↓
kubelet на выбранной Node (скачивает образ, запускает контейнер)
        ↓
Pod Running
```

### 2.4 Reconcile loop (принцип работы контроллеров)

Контроллер постоянно сравнивает **desired state** (то, что в etcd) с **actual state** (то, что реально запущено). При расхождении — корректирует.

---

## 3. Базовые объекты Kubernetes

### 3.1 Namespace

Логическая изоляция ресурсов внутри кластера. Ресурсы в разных NS не видят друг друга по короткому имени.

```bash
kubectl get namespaces
kubectl get ns
kubectl apply -f namespace.yaml
kubectl create ns my-namespace
kubectl -n my-namespace get pods
```

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: home-dev
```

### 3.2 Pod

Минимальная развёртываемая единица. Содержит один или несколько контейнеров, которые:
- делят один сетевой namespace (один IP)
- делят volumes
- всегда запускаются на одной Node

**Жизненный цикл Pod:**
```
Pending → Running → Succeeded
                  → Failed
                  → Unknown
```

- `Pending` — Pod принят, ждёт планирования или загрузки образа
- `Running` — хотя бы один контейнер работает
- `Succeeded` — все контейнеры завершились с кодом 0
- `Failed` — хотя бы один контейнер завершился с ненулевым кодом

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: simple-web
  namespace: home-dev
spec:
  containers:
    - name: web
      image: nginx:1.25
      ports:
        - containerPort: 8080
```

```bash
kubectl apply -f pod.yaml
kubectl -n home-dev get pods
kubectl -n home-dev describe pod simple-web
kubectl -n home-dev logs simple-web
kubectl -n home-dev port-forward pod/simple-web 8080:8080
kubectl delete -f pod.yaml
```

### 3.3 ReplicaSet

Гарантирует, что всегда запущено N копий Pod. Если Pod падает — создаёт новый.

**Важно:** ReplicaSet не умеет обновлять уже запущенные Pods при изменении образа. Нужно удалять Pods вручную.

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: goapp-rs
  namespace: home-dev
spec:
  replicas: 3
  selector:
    matchLabels:
      app: goapp       # ReplicaSet ищет Pods с этим label
  template:
    metadata:
      labels:
        app: goapp     # Pods получают этот label
    spec:
      containers:
        - name: web
          image: myapp:1.0.0
```

### 3.4 Deployment

Управляет ReplicaSet. Добавляет:
- **Rolling update** — обновление без простоя
- **Rollback** — откат на предыдущую версию
- **История** — хранит старые ReplicaSets

```
Deployment → ReplicaSet (v2, 3 pods running)
           → ReplicaSet (v1, 0 pods, хранится для rollback)
```

```bash
kubectl -n home-dev apply -f deployment.yaml
kubectl -n home-dev rollout status deployment/go-http-server
kubectl -n home-dev rollout history deployment/go-http-server
kubectl -n home-dev rollout undo deployment/go-http-server
kubectl -n home-dev rollout restart deployment/go-http-server
kubectl -n home-dev scale deployment/go-http-server --replicas=5
kubectl diff -f deployment.yaml   # что изменится перед apply
```

### 3.5 DaemonSet

Гарантирует, что Pod запущен на **каждой** Node кластера (или на выбранных через nodeSelector).

**Типичное применение:**
- Сбор логов (Fluentd, Filebeat)
- Мониторинг (node-exporter)
- Сетевые плагины (Calico, Flannel)

При добавлении новой Node Pod автоматически создаётся на ней.

### 3.6 Job

Одноразовая задача. Pod запускается, выполняет работу и завершается.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi-job
spec:
  backoffLimit: 4       # сколько раз повторить при ошибке
  template:
    spec:
      containers:
        - name: pi
          image: python:3.9
          command: ["python", "-c", "from math import pi; print(pi)"]
      restartPolicy: Never   # Never или OnFailure
```

```bash
kubectl apply -f job.yaml
kubectl -n home-dev get jobs
kubectl -n home-dev get po
kubectl -n home-dev logs <pod-name>
kubectl explain job.spec.backoffLimit
```

### 3.7 CronJob

Периодически запускает Job по расписанию (cron-синтаксис).

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/2 * * * *"   # каждые 2 минуты
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: hello
              image: busybox
              command: ["echo", "hello"]
          restartPolicy: Never
```

**Cron формат:** `минуты часы день_месяца месяц день_недели`

```bash
kubectl -n home-dev get cronjob
kubectl -n home-dev get jobs
kubectl -n home-dev describe cronjob hello
# Запустить вручную из CronJob:
kubectl -n home-dev create job --from=cronjob/hello hello-manual
```

### 3.8 Сравнение workload-объектов

| Объект | Когда использовать |
|--------|-------------------|
| `Pod` | Отладка, тесты. Никогда в production напрямую |
| `ReplicaSet` | Почти никогда напрямую — используйте Deployment |
| `Deployment` | Stateless приложения (веб-сервисы, API) |
| `DaemonSet` | Один Pod на каждой ноде (мониторинг, логи) |
| `StatefulSet` | Stateful приложения (БД, Kafka) — с сохранением имён и порядка |
| `Job` | Одноразовые задачи (миграции, расчёты) |
| `CronJob` | Периодические задачи (бэкапы, отчёты) |

---

## 4. Сеть и доступ

### 4.1 Service

Service — стабильный сетевой эндпоинт для набора Pods (через selector по labels).

**Типы Services:**

| Тип | Доступ | Когда использовать |
|-----|--------|--------------------|
| `ClusterIP` | Только внутри кластера | Внутренние сервисы, микросервисы |
| `NodePort` | Снаружи через порт ноды (30000–32767) | Dev/тесты без облачного LB |
| `LoadBalancer` | Внешний облачный балансировщик | Production в облаке |
| `ExternalName` | Редирект на внешнее DNS-имя | Интеграция с внешними сервисами |

**ClusterIP:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: goweb-svc
  namespace: demo
spec:
  selector:
    app: goapp          # находит Pods с этим label
  ports:
    - port: 80          # порт Service
      targetPort: 8070  # порт Pod
```

**NodePort:**
```yaml
spec:
  type: NodePort
  selector:
    app: goapp
  ports:
    - port: 80
      targetPort: 8070
      nodePort: 30050   # фиксированный порт на ноде (30000-32767)
```

**LoadBalancer:**
```yaml
spec:
  type: LoadBalancer
  selector:
    app: goapp
  ports:
    - port: 80
      targetPort: 8080
```

### 4.2 Service Discovery

Kubernetes DNS автоматически создаёт записи для каждого Service.

**Формат DNS:**
```
<service-name>.<namespace>.svc.<cluster-domain>
# Пример:
goweb-svc.demo.svc.cluster.local
# Из того же namespace достаточно:
goweb-svc
```

Проверка из Pod:
```bash
kubectl run -it --image=ubuntu --restart=Never shell -- \
  sh -c 'apt-get install -y dnsutils && nslookup goweb-svc.demo'
```

### 4.3 Ingress

Ingress управляет HTTP/HTTPS маршрутизацией извне кластера к Services.

**Нужен IngressController** (nginx, traefik, istio и т.д.):
```bash
# Установка nginx ingress controller:
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
```

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  namespace: demo-ingress
spec:
  ingressClassName: nginx
  rules:
    - host: hello-world.example       # маршрут по hostname
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: goweb-svc
                port:
                  number: 80
```

**Ingress vs Service LoadBalancer:**
- LoadBalancer — один IP на один Service
- Ingress — один IP, маршрутизация по hostname/path на много Services

---

## 5. Lifecycle: Probes, Limits, ConfigMap, Secret

### 5.1 Probes (Пробы)

Kubernetes проверяет здоровье контейнеров через три типа проб:

| Проба | Действие при провале | Назначение |
|-------|---------------------|------------|
| `livenessProbe` | Перезапустить контейнер | Контейнер завис, не отвечает |
| `readinessProbe` | Убрать Pod из балансировки Service | Pod не готов принимать трафик |
| `startupProbe` | Перезапустить контейнер | Защищает от срабатывания liveness при долгом старте |

**Типы проверок:**
```yaml
# HTTP GET
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 10   # подождать перед первой проверкой
  periodSeconds: 5          # интервал между проверками
  failureThreshold: 3       # сколько провалов до действия

# Выполнить команду
readinessProbe:
  exec:
    command: ["cat", "/tmp/ready"]

# TCP соединение
livenessProbe:
  tcpSocket:
    port: 5432
```

**Важно:** readinessProbe не убивает Pod — только убирает из Endpoints Service.

### 5.2 Resource Requests и Limits

```yaml
resources:
  requests:        # минимум для планировщика (scheduler)
    memory: "64Mi"
    cpu: "250m"    # 250 millicores = 0.25 CPU
  limits:          # максимум (container runtime enforces)
    memory: "128Mi"
    cpu: "500m"
```

| | Requests | Limits |
|---|----------|--------|
| Роль | Scheduler использует для выбора Node | Container Runtime ограничивает потребление |
| Превышение CPU | Throttling (замедление) | — |
| Превышение Memory | — | OOMKilled (Pod перезапускается) |

**QoS классы:**
- `Guaranteed` — requests == limits (лучшая предсказуемость)
- `Burstable` — requests < limits
- `BestEffort` — не указаны ни requests, ни limits (вытесняется первым)

### 5.3 ConfigMap

Хранение несекретной конфигурации.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: demo
data:
  LOG_LEVEL: "debug"
  DATABASE_URL: "postgres://db:5432/mydb"
  config.yaml: |
    log_level: debug
    database: postgres://db:5432/mydb
```

**Способы подключения:**

```yaml
# 1. Все ключи как переменные окружения
envFrom:
  - configMapRef:
      name: app-config

# 2. Один ключ как переменная
env:
  - name: LOG_LEVEL
    valueFrom:
      configMapKeyRef:
        name: app-config
        key: LOG_LEVEL

# 3. Монтирование как файлы в volume
volumes:
  - name: config-vol
    configMap:
      name: app-config
containers:
  - volumeMounts:
    - name: config-vol
      mountPath: /etc/config
```

### 5.4 Secret

Хранение чувствительных данных. Значения хранятся в base64 (не шифрование!).

**Типы:**
| Тип | Назначение |
|-----|------------|
| `Opaque` | Произвольные данные (пароли, токены) |
| `kubernetes.io/tls` | TLS-сертификат и ключ |
| `kubernetes.io/dockerconfigjson` | Credentials для Docker Registry |
| `kubernetes.io/service-account-token` | Токен для ServiceAccount |

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
data:
  password: cGFzc3dvcmQ=   # base64("password")
  username: YWRtaW4=        # base64("admin")
```

```bash
# Создать из литерала
kubectl create secret generic my-secret --from-literal=password=mypass
# Создать из файла
kubectl create secret generic tls-secret --from-file=tls.crt --from-file=tls.key
```

**TLS Secret для Ingress:**
```bash
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config cert.conf
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
kubectl create secret tls my-tls --cert=server.crt --key=server.key
```

**Способы подключения Secret — те же, что у ConfigMap** (`envFrom`, `env.valueFrom`, `volumeMounts`).

### 5.5 RollingUpdate — стратегия обновления

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1   # макс. недоступных Pods во время обновления
      maxSurge: 1         # макс. дополнительных Pods сверх replicas
```

**Процесс:**
1. Создаётся новый ReplicaSet с новым образом
2. Постепенно поднимаются новые Pods (не превышая `replicas + maxSurge`)
3. Постепенно убиваются старые Pods (не опускаясь ниже `replicas - maxUnavailable`)

**Другая стратегия — Recreate:**
```yaml
strategy:
  type: Recreate   # убить все старые Pods, потом создать новые (кратковременный downtime)
```

---

## 6. Debug: сценарии и диагностика

### 6.1 Алгоритм диагностики

```bash
kubectl get pods -n <namespace>          # смотрим статус
kubectl describe pod <pod> -n <ns>       # события и детали
kubectl logs <pod> -n <ns>               # логи
kubectl logs <pod> -n <ns> --previous    # логи предыдущего контейнера (после crash)
kubectl get events -n <ns> --sort-by='.lastTimestamp'
```

### 6.2 Статусы и их причины

| Статус | Причина | Что проверить |
|--------|---------|---------------|
| `ImagePullBackOff` / `ErrImagePull` | Образ не найден или нет доступа к registry | Имя образа, тег, `imagePullSecrets` |
| `ContainerCreating` (зависло) | PVC не смонтирован, Secret/ConfigMap не найден | `kubectl describe pod`, events |
| `Pending` | Нет ресурсов на нодах, не прошёл scheduling | `kubectl describe pod`, секция Events — `Insufficient CPU/memory`, taints |
| `READY 0/1` | readinessProbe не проходит | Лог приложения, ответ на health endpoint |
| `CrashLoopBackOff` | Контейнер падает сразу после старта | `kubectl logs --previous`, ошибка в CMD/ENTRYPOINT |
| Pods нет (0 ready в Deployment) | Неправильный selector в Deployment/RS | `kubectl describe deployment`, проверить `selector.matchLabels` и labels в template |
| Networking issue | Неправильный selector в Service или неверный port | `kubectl describe svc`, проверить `selector` и `targetPort` |

### 6.3 kubectl debug — ephemeral containers

Добавляет временный отладочный контейнер в запущенный Pod (без перезапуска):

```bash
# Добавить ubuntu-контейнер в Pod для отладки
kubectl debug -it <pod-name> --image=ubuntu -- bash

# Создать копию Pod с другим образом для отладки
kubectl debug <pod-name> -it \
  --copy-to=debug-pod \
  --image=ubuntu \
  --container=debug-container -- sh
```

### 6.4 nsenter — отладка на уровне Node

Вход в namespace контейнера напрямую с ноды (когда `kubectl exec` недоступен):

```bash
# На ноде: найти PID контейнера
sudo crictl ps | grep <pod-name>
sudo crictl inspect --output go-template --template '{{.info.pid}}' <container-id>

# Войти в network namespace контейнера
sudo nsenter -t <pid> -n netstat -tulpn
sudo nsenter -t <pid> -n ip link show
sudo nsenter -t <pid> -n dig google.com

# Войти в mount namespace
sudo nsenter -t <pid> -m cat /etc/resolv.conf
```

### 6.5 Проверка сети между Pod'ами

```bash
# Запустить утилиту curl внутри кластера
kubectl run curl --image=radial/busyboxplus:curl -i --tty --rm
# Внутри:
curl http://goweb-svc.demo.svc.cluster.local:80/health

# Проверить DNS
kubectl run -it --image=busybox --restart=Never dns-test -- nslookup kubernetes.default

# Проверить TCP-порт
echo "yes" > /dev/tcp/10.0.11.3/5432 && echo "open" || echo "close"
```

---

## 7. Инструменты и lifehacks

### 7.1 Алиасы kubectl

```bash
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgpo='kubectl get pods -o wide'
alias kgs='kubectl get secret'
alias kdp='kubectl describe pod'
alias kl='kubectl logs'

# Добавить в ~/.bashrc + source ~/.bashrc
```

### 7.2 Автодополнение

```bash
sudo apt install bash-completion
echo 'source <(kubectl completion bash)' >> ~/.bashrc
echo 'alias k=kubectl' >> ~/.bashrc
echo 'complete -o default -F __start_kubectl k' >> ~/.bashrc
source ~/.bashrc
```

### 7.3 Полезные команды kubectl

```bash
# Получить YAML существующего ресурса
kubectl get pod nginx -o yaml
kubectl get pod nginx -o yaml > nginx.yaml

# Сухой прогон (не применять, только вывести YAML)
kubectl run nginx --image=nginx --dry-run=client -o yaml

# Смотреть ресурсы в реальном времени
watch kubectl -n home-dev get pods

# Посмотреть дерево ресурсов
kubectl krew install tree
kubectl tree deployment my-deploy

# Capacity ресурсов
kubectl krew install resource-capacity
kubectl resource-capacity --pods
```

### 7.4 CoreDNS — добавить внутреннее DNS-имя

```bash
kubectl -n kube-system edit cm coredns
# Добавить в NodeHosts:
# 10.0.2.15 my-internal-service.example.com

# Перезапустить CoreDNS
kubectl -n kube-system delete pod -l k8s-app=kube-dns
```

---

## 8. CI/CD: GitLab + Kubernetes

### 8.1 GitLab Runner в Kubernetes

```bash
# Добавить helm-репозиторий
helm repo add gitlab https://charts.gitlab.io
helm repo update gitlab

# Создать namespace
kubectl create ns gitlab-runner

# Установить Runner (токен берём из GitLab → Settings → CI/CD → Runners)
helm upgrade --install -n gitlab-runner gitlab-runner gitlab/gitlab-runner -f values.yaml
```

### 8.2 Пример .gitlab-ci.yml

```yaml
stages:
  - build
  - test
  - deploy

build:
  stage: build
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

test:
  stage: test
  script:
    - go test ./...

deploy:
  stage: deploy
  script:
    - helm upgrade --install myapp ./chart
        --set image.tag=$CI_COMMIT_SHA
        -n production
```

### 8.3 Kaniko — сборка образов без Docker daemon

Kaniko строит образы внутри Pod без привилегий Docker socket:
```yaml
# В .gitlab-ci.yml
build:
  image:
    name: gcr.io/kaniko-project/executor:v1.9.0
    entrypoint: [""]
  script:
    - /kaniko/executor
        --context $CI_PROJECT_DIR
        --dockerfile $CI_PROJECT_DIR/Dockerfile
        --destination $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
```

### 8.4 Helm — основные команды

```bash
helm repo add myrepo https://charts.example.com
helm repo update
helm search repo myrepo

helm install myapp myrepo/myapp -f values.yaml -n production
helm upgrade --install myapp myrepo/myapp -f values.yaml -n production
helm rollback myapp 1 -n production
helm uninstall myapp -n production
helm list -n production
helm get values myapp -n production
```

---

## 9. Типичные вопросы на собеседовании

### Docker

**Q: Чем отличается CMD от ENTRYPOINT?**
ENTRYPOINT — неизменяемая точка входа, CMD — аргументы по умолчанию к ней. CMD легко переопределить в `docker run`, ENTRYPOINT — только через `--entrypoint`.

**Q: Что такое multi-stage build и зачем?**
Позволяет использовать большой образ для сборки (golang, node) и маленький для финального образа. Итоговый образ не содержит компилятор/исходники — меньше размер, меньше уязвимостей.

**Q: Чем volume отличается от bind mount?**
Volume управляется Docker, хранится в `/var/lib/docker/volumes`. Bind mount — монтирование конкретного пути хоста. Volume переносимее, bind mount удобнее при разработке.

### Kubernetes

**Q: Из чего состоит Control Plane?**
API Server (точка входа), etcd (хранилище состояния), Scheduler (выбирает Node для Pod), Controller Manager (reconcile loop).

**Q: Чем Deployment отличается от ReplicaSet?**
Deployment управляет ReplicaSet и добавляет: rolling update без downtime, rollback (`rollout undo`), историю версий (старые RS хранятся).

**Q: Разница между liveness и readiness probe?**
liveness — если падает, контейнер перезапускается. readiness — если падает, Pod убирается из балансировки Service, но не перезапускается.

**Q: Чем requests отличается от limits?**
requests — сколько ресурсов scheduler резервирует при размещении Pod на Node. limits — максимум, который container runtime позволит использовать.

**Q: Как работает RollingUpdate?**
Создаёт новый ReplicaSet с новым образом, постепенно поднимает новые Pods и убивает старые. `maxUnavailable` — сколько Pods может быть недоступно, `maxSurge` — сколько дополнительных Pods можно создать.

**Q: Чем ConfigMap отличается от Secret?**
ConfigMap — несекретная конфигурация (plain text). Secret — чувствительные данные (base64, не шифрование!). Секреты можно настроить на encryption at rest в etcd.

**Q: Что такое Ingress и зачем нужен IngressController?**
Ingress — объект Kubernetes, описывающий правила HTTP-маршрутизации. Сам по себе ничего не делает — нужен IngressController (nginx, traefik), который читает эти правила и настраивает балансировщик.

**Q: Как Pod'ы находят друг друга?**
Через Service DNS: `<service>.<namespace>.svc.cluster.local`. CoreDNS создаёт записи для каждого Service автоматически.

**Q: Что такое DaemonSet?**
Гарантирует запуск одного Pod на каждой Node. Используется для агентов мониторинга, сбора логов, сетевых плагинов.

**Q: CrashLoopBackOff — что это и что делать?**
Pod постоянно падает и перезапускается. Делать: `kubectl logs <pod> --previous` — посмотреть логи упавшего контейнера, найти причину. Обычно: ошибка в CMD, нет нужного файла/переменной, приложение не может подключиться к БД.

**Q: Как попасть внутрь контейнера для отладки?**
`kubectl exec -it <pod> -- bash` — если bash есть. `kubectl debug -it <pod> --image=ubuntu` — если bash нет (добавляет ephemeral container). `nsenter` — если нужен доступ на уровне ноды.

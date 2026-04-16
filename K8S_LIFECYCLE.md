# Kubernetes — Lifecycle: Probes, Resources, ConfigMap, Secret

---

## 1. Probes — проверки здоровья

> Аналогия:
> - **livenessProbe** — врач, который проверяет "живой ли пациент?". Если нет — реанимация (перезапуск).
> - **readinessProbe** — HR, который проверяет "готов ли сотрудник принимать задачи?". Если нет — задачи не дают, но не увольняют.
> - **startupProbe** — таймер для новичка: "не трогать первые 30 секунд, он ещё стартует".

### Сравнение проб

| | livenessProbe | readinessProbe | startupProbe |
|---|---|---|---|
| При провале | Перезапуск контейнера | Убрать из балансировки Service | Перезапуск контейнера |
| Цель | Обнаружить зависший контейнер | Не слать трафик на неготовый Pod | Защитить медленно стартующий контейнер |
| Убивает Pod? | Да (перезапускает) | Нет | Да (перезапускает) |
| Типы проверок | httpGet, exec, tcpSocket | httpGet, exec, tcpSocket | httpGet, exec, tcpSocket |

**Важно:** readinessProbe не убивает Pod — только убирает из Endpoints Service.

### YAML-примеры проб

```yaml
# HTTP GET — самый распространённый
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
  periodSeconds: 5

# TCP соединение (для БД, Redis)
livenessProbe:
  tcpSocket:
    port: 5432
  initialDelaySeconds: 15

# startupProbe — пока не пройдёт, liveness не проверяется
startupProbe:
  httpGet:
    path: /healthz
    port: 8080
  failureThreshold: 30    # 30 попыток × 10 сек = 5 минут на старт
  periodSeconds: 10
```

### Практичный паттерн (все три вместе)

```yaml
containers:
  - name: app
    image: myapp:1.0.0
    startupProbe:
      httpGet:
        path: /healthz
        port: 8080
      failureThreshold: 30
      periodSeconds: 10
    livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
      periodSeconds: 10
    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
      periodSeconds: 5
```

---

## 2. Resource Requests и Limits

> Аналогия:
> - **requests** — это минимальная зарплата, которую вы гарантируете сотруднику (scheduler использует для выбора ноды).
> - **limits** — это максимум, который он может потратить на командировки (container runtime ограничивает).

```yaml
resources:
  requests:        # минимум для планировщика (scheduler)
    memory: "64Mi"
    cpu: "250m"    # 250 millicores = 0.25 CPU
  limits:          # максимум (container runtime enforces)
    memory: "128Mi"
    cpu: "500m"
```

### Что происходит при превышении

| Ресурс | Превышение requests | Превышение limits |
|---|---|---|
| CPU | Разрешено (буст до limits) | Throttling (замедление) |
| Memory | Разрешено | OOMKilled (Pod перезапускается) |

### QoS классы

> Аналогия: "Гарантированный" — бизнес-класс в самолёте. "BestEffort" — standby без места.

| Класс | Условие | Поведение при нехватке памяти |
|---|---|---|
| `Guaranteed` | requests == limits | Вытесняется последним |
| `Burstable` | requests < limits | Вытесняется вторым |
| `BestEffort` | requests/limits не указаны | Вытесняется первым |

**Рекомендация:** всегда указывайте хотя бы requests. Без них scheduler не знает куда ставить Pod.

### HPA и VPA

| | HPA | VPA |
|---|---|---|
| Расшифровка | Horizontal Pod Autoscaler | Vertical Pod Autoscaler |
| Что масштабирует | Количество реплик | requests/limits контейнера |
| По чему | CPU, memory, custom metrics | Историческое потребление |
| Требует | metrics-server | VPA operator |
| Аналогия | Нанять больше сотрудников | Повысить зарплату одному сотруднику |

---

## 3. ConfigMap

> Аналогия: ConfigMap — это **инструкция по эксплуатации** (открытый документ, не секрет).
> Передаём приложению, как ему работать: уровень логов, адрес БД, порт.

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

### Три способа подключить ConfigMap к Pod

```yaml
# 1. Все ключи как переменные окружения
envFrom:
  - configMapRef:
      name: app-config

# 2. Один конкретный ключ как переменная
env:
  - name: LOG_LEVEL
    valueFrom:
      configMapKeyRef:
        name: app-config
        key: LOG_LEVEL

# 3. Монтировать как файлы (config.yaml появится в /etc/config/config.yaml)
volumes:
  - name: config-vol
    configMap:
      name: app-config
containers:
  - volumeMounts:
    - name: config-vol
      mountPath: /etc/config
```

**Обновление:**
- При монтировании как volume — обновляется автоматически (задержка ~1 минута).
- При использовании как env-переменные — нужен перезапуск Pod.

---

## 4. Secret

> Аналогия: Secret — это **запечатанный конверт** с конфиденциальными данными.
> Base64 — это не замок, а просто конверт. Настоящая защита — RBAC + encryption at rest.

**Важно: base64 — это кодирование, не шифрование!**

### Типы Secrets

| Тип | Назначение |
|---|---|
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

**imagePullSecret** — для доступа к приватному registry:
```yaml
spec:
  imagePullSecrets:
    - name: registry-secret
```

### ConfigMap vs Secret

| | ConfigMap | Secret |
|---|---|---|
| Данные | Несекретные | Чувствительные |
| Хранение в etcd | Plaintext | Base64 (можно включить encryption at rest) |
| Типы | - | Opaque, tls, dockerconfigjson, ... |
| Подключение | envFrom, env, volume | envFrom, env, volume |
| Обновление в Pod | Через restart (env) / авто (volume) | Через restart (env) / авто (volume) |

---

## 5. RollingUpdate — стратегия обновления

> Аналогия: RollingUpdate — это как **обновление самолёта прямо в полёте**.
> Меняем двигатели один за другим, не приземляясь (no downtime).

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1   # макс. недоступных Pods во время обновления
      maxSurge: 1         # макс. дополнительных Pods сверх replicas
```

**Процесс:**
1. Создаётся новый ReplicaSet с новым образом.
2. Поднимаются новые Pods (не превышая `replicas + maxSurge`).
3. Убиваются старые Pods (не опускаясь ниже `replicas - maxUnavailable`).

**Recreate (альтернатива):**
```yaml
strategy:
  type: Recreate   # убить всё → создать новое (есть downtime)
```

| | RollingUpdate | Recreate |
|---|---|---|
| Downtime | Нет | Да (кратковременный) |
| Сложность | Выше (нужна совместимость версий) | Проще |
| Применение | Production API | Когда нельзя иметь 2 версии одновременно |

---

## 6. Жизненный цикл Pod и контейнера

### Init Containers

> Аналогия: Init Containers — это **подготовительная работа перед открытием магазина**.
> Сначала завоз товара (ждём БД), потом открытие (основное приложение).

```yaml
initContainers:
  - name: wait-db
    image: busybox
    command: ['sh', '-c', 'until nc -z db-service 5432; do sleep 2; done']
  - name: run-migrations
    image: myapp-migrator
    command: ['./migrate']
```

### Sidecar контейнер

> Аналогия: Sidecar — это **помощник в той же машине**.
> Основное приложение едет, sidecar рядом: собирает логи, проксирует трафик, обновляет секреты.

Примеры: Fluentd (логи), Envoy/Istio proxy (service mesh), Vault Agent (секреты).

### Graceful Shutdown

```
kubectl delete pod → SIGTERM → ждём terminationGracePeriodSeconds (30 сек) → SIGKILL
```

**preStop hook** — дать время на graceful drain:
```yaml
lifecycle:
  preStop:
    exec:
      command: ["/bin/sh", "-c", "sleep 5"]
```

### PodDisruptionBudget (PDB)

> Аналогия: PDB — это **минимальный штат на дежурстве**.
> При техобслуживании нод K8s выселяет Pods — PDB гарантирует, что хотя бы N реплик останутся живыми.

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: myapp-pdb
spec:
  minAvailable: 2      # или maxUnavailable: 1
  selector:
    matchLabels:
      app: myapp
```

---

## 7. RBAC — контроль доступа

> Аналогия: RBAC — это **пропускная система в офисе**.
> - Role — список что можно делать в конкретном этаже (namespace).
> - ClusterRole — список что можно делать во всём здании.
> - RoleBinding — выдать пропуск конкретному человеку.
> - ServiceAccount — удостоверение личности для процессов (Pod).

| Объект | Область | Назначение |
|---|---|---|
| `Role` | Namespace | Права в пределах одного NS |
| `ClusterRole` | Весь кластер | Права на уровне кластера |
| `RoleBinding` | Namespace | Привязать Role/ClusterRole к субъекту в NS |
| `ClusterRoleBinding` | Весь кластер | Привязать ClusterRole к субъекту глобально |
| `ServiceAccount` | Namespace | Идентификатор для Pod |

---

## 8. Безопасность

### SecurityContext

```yaml
securityContext:
  runAsUser: 1000         # не root
  runAsNonRoot: true
  readOnlyRootFilesystem: true   # нет записи в FS контейнера
  capabilities:
    drop: ["ALL"]
```

### Pod Security Standards (PSS)

| Уровень | Ограничения |
|---|---|
| `privileged` | Без ограничений |
| `baseline` | Базовая защита (нет privileged, hostNetwork) |
| `restricted` | Строгие ограничения (runAsNonRoot, read-only FS) |

```bash
# Применить на namespace
kubectl label namespace production pod-security.kubernetes.io/enforce=restricted
```

# K8S FULL GUIDE — pgo-dashboards-api

## Словарь терминов

- **K8S** — Kubernetes, платформа для запуска контейнеров.
- **CI** (Continuous Integration) — автоматическая сборка и проверка кода.
- **CD** (Continuous Delivery) — автоматическая доставка в среду.
- **Deploy** — выкат новой версии в кластер.
- **Namespace** — логическая "папка" ресурсов внутри кластера.
- **Manifest** — YAML-файл с описанием ресурса Kubernetes.
- **Job** — одноразовая задача (запускается, отрабатывает, завершается).
- **Deployment** — описание долгоживущего приложения (pod-ов).
- **Pod** — запущенный контейнер в кластере.
- **Service** — стабильный внутренний адрес для pod-ов.
- **Ingress** — правила входящего HTTP/HTTPS-трафика снаружи.
- **DNS** — система перевода имен в адреса.
- **DB** — база данных.
- **Overlay** — набор изменений поверх общего шаблона.
- **Patch** — точечное изменение YAML-файла.
- **Rollout** — процесс обновления pod-ов.
- **SOPS** — инструмент шифрования секретов (env-файлов).
- **Registry** — хранилище Docker-образов.
- **Image tag** — версия образа (например `backend-http:1.2.3`).
- **ExternalName** — тип сервиса, который является DNS-алиасом на внешний адрес.
- **Init-контейнер** — вспомогательный контейнер, который выполняется до старта основного.

---

## Используемые паттерны

1. **Base + Overlays**
   Общий шаблон (`base`) + изменения под каждую среду (`overlays`).
   Один и тот же каркас, разные параметры для `test`, `go`, `bat` и т.д.

2. **Immutable Artifact (неизменяемый артефакт)**
   CI собирает образ один раз и присваивает ему тег версии (`BUILD_VERSION`).
   CD выкатывает именно этот конкретный образ — никакой пересборки при деплое.

3. **Migration-first**
   Сначала запускаются миграции БД, только потом стартует приложение.
   Защищает от ситуации, когда новый код работает со старой схемой БД.

4. **Init-gate**
   Deployment-ы не стартуют, пока не завершится migration Job.
   Реализуется через init-контейнер `k8s-wait-for`.

5. **External dependency via Service alias**
   Подключение к БД идет через `ExternalName` (DNS-алиас).
   Приложение не знает реального адреса БД — его подставляет overlay.

6. **Deterministic deploy order**
   В TeamCity шаги deploy идут в строгом порядке:
   версия -> секреты -> migration -> apply всего -> rollout restart.

---

## 0) Как всё связано (общая схема)

Git-репозиторий хранит **и код приложения, и описание инфраструктуры**.  
TeamCity забирает весь репозиторий и делает две вещи: собирает образы и выкатывает их в кластер.

```
Git (код приложения + kubernetes/ манифесты)
    │
    ▼
TeamCity забирает репозиторий целиком (checkout)
    │
    ├── 1. Читает cmd/ → собирает Go-бинарник
    ├── 2. Читает kubernetes/build/Dockerfile → собирает Docker-образы
    ├── 3. Пушит образы в приватный registry (registry.proffitgo.com)
    │
    └── 4. Читает kubernetes/deploy/overlays/<env>/
            ├── Подставляет BUILD_VERSION в манифесты
            ├── Расшифровывает секреты через SOPS
            └── kubectl apply -k ./ → применяет в Kubernetes кластер
```

Почему `kubernetes/` лежит в том же Git-репозитории:
- изменения в манифестах видны в истории Git рядом с кодом,
- код и его деплой-конфигурация всегда соответствуют друг другу в одной ветке,
- это классический паттерн **GitOps**: инфраструктура описана в Git и применяется оттуда.

---

## 1) Главное (30 секунд)

- Деплой через `Kustomize`: общий шаблон (`base`) + изменения под среду (`overlay`).
- Порядок старта: `Job миграции` → `Deployment API/sync` → `Service` → `Ingress`.
- API и sync не стартуют, пока не завершится `dashboards-api-migration`.
- Версия образов берется из `BUILD_VERSION` и подставляется автоматически.
- БД подключается через `ExternalName` сервис — DNS-алиас на балансировщик.

---

## 1) Жизненный цикл запуска

```
TeamCity
   └── build образов + push в registry
   └── записывает BUILD_VERSION
   └── расшифровывает секреты (SOPS)
   └── kubectl apply -k ./overlays/go
          └── Kustomize собирает: base + generators + patches
          └── Job dashboards-api-migration
                 └── init: ждет БД (pg_isready)
                 └── main: запускает migrate
          └── Deployment dashboards-api
                 └── init: ждет job migration
                 └── main: стартует API
          └── Deployment dashboards-data-sync
                 └── init: ждет job migration
                 └── main: стартует sync
          └── Service -> получает endpoints API
          └── Ingress -> начинает пропускать внешний трафик
```

---

## 2) Что делает overlay `go`

Файл: `kubernetes/deploy/overlays/go/kustomization.yaml`

1. Ставит namespace: `go-pgo-dashboards-api-golang-release`.
2. Обновляет `ConfigMap build-version` из файла `ConfigMapValues/build-version`.
3. Обновляет `Secret dashboards-api` из расшифрованного `SecretValuesDecrypted/.env`.
4. Подставляет `BUILD_VERSION` как tag образа в Deployment-ы и Job.
5. Накладывает патчи:
   - количество реплик,
   - image repo для API/sync/job,
   - домен/TLS для Ingress,
   - endpoint/порты балансировщика БД.

---

## 3) Главные файлы

- `kubernetes/deploy/overlays/go/kustomization.yaml` — точка входа для окружения `go`.
- `kubernetes/deploy/base/kustomization.yaml` — список всех базовых ресурсов.
- `kubernetes/deploy/base/Job/dashboards-api-migration.yaml` — задача миграций.
- `kubernetes/deploy/base/Deployment/dashboards-api.yaml` — HTTP API.
- `kubernetes/deploy/base/Deployment/dashboards-data-sync.yaml` — фоновый sync.
- `kubernetes/deploy/base/Service/dashboards-api.yaml` — внутренний адрес API.
- `kubernetes/deploy/base/Ingress/dashboards-api.yaml` — входящий трафик.
- `kubernetes/deploy/base/Service/external-db-load-balancer.yaml` — DNS-алиас до БД.

---

## 4) TeamCity: CI build + push

Этот шаг отвечает за создание артефактов (образов) и их публикацию.

```bash
#!/usr/bin/env bash
set -euo pipefail
# -e: стоп при любой ошибке
# -u: стоп при обращении к несуществующей переменной
# pipefail: ошибка в любом месте пайпа = ошибка шага

# Версия берется из TeamCity build number или commit SHA
BUILD_VERSION="${BUILD_VERSION}"

# Логин в приватный registry
echo "${REGISTRY_PASSWORD}" | docker login "${REGISTRY}" \
  --username "${REGISTRY_USER}" --password-stdin

# Сборка HTTP образа (из Dockerfile, target = application_http)
docker build \
  -f kubernetes/build/Dockerfile \
  --target application_http \
  --build-arg IMAGE_VERSION=1.24-alpine \
  --build-arg UID=1001 \
  --build-arg USER=appuser \
  --build-arg MAIN_PACKAGE_HTTP=cmd/backend \
  -t "${REGISTRY}/backend-http:${BUILD_VERSION}" .

# Сборка SYNC образа (тот же Dockerfile, другой target)
docker build \
  -f kubernetes/build/Dockerfile \
  --target application_sync \
  --build-arg IMAGE_VERSION=1.24-alpine \
  --build-arg UID=1001 \
  --build-arg USER=appuser \
  --build-arg MAIN_PACKAGE_SYNC=cmd/async \
  -t "${REGISTRY}/backend-sync:${BUILD_VERSION}" .

# Публикация обоих образов в registry
docker push "${REGISTRY}/backend-http:${BUILD_VERSION}"
docker push "${REGISTRY}/backend-sync:${BUILD_VERSION}"

# Передача BUILD_VERSION в следующий TeamCity step
echo "##teamcity[setParameter name='system.BUILD_VERSION' value='${BUILD_VERSION}']"
```

---

## 5) TeamCity: CD deploy (реальный прод-вариант)

Этот шаг берет готовые образы и выкатывает их в прод-кластер.

```bash
#!/usr/bin/env bash
set -euo pipefail

# Проверка ручного запуска:
# если запуск руками и не выбран ни один HELM_* флаг,
# запускаются все проекты;
# если шаг GO и HELM_GO не выбран, этот шаг пропускается.
if [[ -n "%teamcity.build.triggeredBy.username%" ]] \
   && [[ -z "%HELM_GO%" ]] && [[ -z "%HELM_BWG%" ]] && [[ -z "%HELM_FAW%" ]] && [[ -z "%HELM_IGETIS%" ]] && [[ -z "%HELM_BAT%" ]]; then
  echo "Ручной запуск без кастомизации. Запускаются все проекты."
else
  if [[ -n "%teamcity.build.triggeredBy.username%" ]] && [[ -z "%HELM_GO%" ]]; then
    echo "Пропуск шага из-за ручного запуска и HELM_GO  не выбран."
    exit 0
  fi
fi

export SOPS_AGE_KEY_FILE="%system.SOPS_AGE_KEY_FILE%"
CTX="ci-agent@proffit-prod-cluster-spb-ru-9"
NS="go-pgo-dashboards-api-golang-release"

BUILD_VERSION=%dep.ProffitBackend_PgoDashboardsApiGolang_Release_Build.system.BUILD_VERSION%
echo "BUILD_VERSION=${BUILD_VERSION}" | tee ./ConfigMapValues/build-version

# Без этой проверки сборка может падать,
# если директория уже есть на агенте
[ -d ./SecretValuesDecrypted ] || mkdir ./SecretValuesDecrypted
/scripts/sops_decrypt.sh --verbose --directory=./SecretValuesEncrypted --output-dir=./SecretValuesDecrypted

source /scripts/fillup_kubeconfig_var.sh

# Удаляем старый migration Job, чтобы миграции могли выполниться заново
kubectl --context $CTX -n $NS delete job -l app.kubernetes.io/component=migration --ignore-not-found

# Сначала применяем только migration-сущности
kubectl --context $CTX apply -k ./ -l app.kubernetes.io/component=migration

# Затем применяем всё остальное
kubectl --context $CTX apply -k ./

kubectl --context $CTX -n $NS rollout restart deployment
```

Ключевые отличия от тестового варианта:

- другой контекст кластера: `proffit-prod-cluster-spb-ru-9`;
- другой namespace: `go-pgo-dashboards-api-golang-release`;
- версия берется из release build chain:
  `%dep.ProffitBackend_PgoDashboardsApiGolang_Release_Build.system.BUILD_VERSION%`;
- добавлена логика ручного запуска с `HELM_*` флагами (выбор, какие проекты деплоить).

---

## 6) Подключение к БД (`ExternalName`)

`ExternalName` — это сервис, который не ведет на pod-ы, а является DNS-алиасом.

В `base` — общий placeholder:
```yaml
type: ExternalName
externalName: db-load-balancer.db-load-balancer.svc.cluster.local
```

В overlay `go` — патч заменяет адрес и порты:
```yaml
externalName: go-db-lb.db-load-balancer.svc.cluster.local
ports:
  - name: db-master
    port: 5040
  - name: db-replica
    port: 5041
  ...
```

Смысл: приложение всегда обращается к `external-db-load-balancer`,
а реальный endpoint и порты задаются в overlay под каждую среду.

---

## 7) Диагностика проблем

### A) Pod завис в статусе `Init`

Причина: init-контейнер `wait-for-migration` ждет migration Job.

```bash
kubectl get pods -n go-pgo-dashboards-api-golang-release
kubectl describe pod <pod-name> -n go-pgo-dashboards-api-golang-release
kubectl get job -n go-pgo-dashboards-api-golang-release
kubectl logs job/dashboards-api-migration -n go-pgo-dashboards-api-golang-release
```

### B) Migration Job в статусе `Failed`

Проверить: доступность БД, корректность env из секрета, image tag, аргумент `migrate`.

```bash
kubectl describe job dashboards-api-migration -n go-pgo-dashboards-api-golang-release
kubectl logs job/dashboards-api-migration -n go-pgo-dashboards-api-golang-release
```

### C) Ошибки 404/502 снаружи

Проверять по цепочке: Ingress → Service → Endpoints → Pod readiness.

---

## 8) Чеклист перед релизом

```bash
# Посмотреть итоговые манифесты после сборки Kustomize
kubectl kustomize kubernetes/deploy/overlays/go

# Проверить что apply пройдет без ошибок (без реального применения)
kubectl apply -k kubernetes/deploy/overlays/go --dry-run=client

# Посмотреть что изменится в кластере
kubectl diff -k kubernetes/deploy/overlays/go

# Проверить статус migration Job
kubectl get job -n go-pgo-dashboards-api-golang-release

# Посмотреть логи миграций
kubectl logs job/dashboards-api-migration -n go-pgo-dashboards-api-golang-release
```

---

## 9) Готовый ответ на собеседовании

"У нас деплой построен на Kustomize по схеме `base + overlays`.
Есть общий шаблон ресурсов и окруженческие изменения под каждую среду.
Поднимаются два deployment — API и sync — и migration Job.
Deployment-ы через init-контейнер ждут завершения миграций, чтобы не стартовать на старой схеме БД.
Версия образов подставляется автоматически через `BUILD_VERSION`.
БД подключается через `ExternalName` сервис — DNS-алиас, который меняется по среде.
В TeamCity сборка делает образы и пушит их, а deploy расшифровывает секреты, запускает миграции и выкатывает всё в кластер."

---

## 10) Шпаргалка (1 экран)

```
base        = общий шаблон ресурсов
overlay     = изменения под среду (namespace, image, host, db ports)
Job         = одноразовая задача (миграции)
Deployment  = долгоживущее приложение
Service     = внутренний адрес к pod-ам
Ingress     = вход внешнего трафика
ExternalName = DNS-алиас до балансировщика БД

Порядок старта:
  migration Job -> API Deployment -> sync Deployment -> Service -> Ingress

TeamCity:
  build -> docker build -> push -> BUILD_VERSION
  deploy -> decrypt secrets -> delete old job -> apply migration -> apply all -> rollout restart
```

# CI/CD пайплайн и инфраструктура

---

## Главная схема

```
git push (ветка master / envs/stage1 / envs/prod5 / ...)
    │
    ▼ workflow.rules: ветка → среда (dev2, stage1, prod5, ...)
    │
    ├── STAGE 1: unittest  → make test           (только если изменился код)
    ├── STAGE 2: build     → docker build + push → Yandex Registry
    ├── STAGE 3: deploy    → kubectl apply / rollout restart → Kubernetes
    └── STAGE 4: e2etest   → make testall        (только если изменился код)
```

Стадии выполняются **последовательно**. Внутри стадии джобы могут работать **параллельно**.
`needs:` позволяет нарушить этот порядок — например `deploy` с `needs: []` стартует сразу, не дожидаясь предыдущей стадии.

---

## Паттерны организации — подробно

### GitOps

**Что это:** подход при котором вся конфигурация инфраструктуры хранится в Git и является единственным источником правды (single source of truth). Изменения в инфраструктуре делаются через коммиты, не через ручные команды на серверах.

**Как реализован здесь:**
- Папка `infrastructure/` содержит Kubernetes-манифесты всех сред
- Чтобы изменить конфигурацию пода в prod — нужно изменить `infrastructure/prod5/app.yml` и запушить
- Нельзя просто зайти на сервер и что-то поменять руками — при следующем деплое `kubectl apply` перезапишет изменения тем что в Git

**Преимущества:** история изменений инфраструктуры в `git log`, ревью через MR, откат через `git revert`.

---

### Branch-to-Environment mapping (ветка = среда)

**Что это:** каждая Git-ветка соответствует конкретной среде деплоя. Пуш в ветку = деплой в среду.

**Как работает механически** — через `workflow.rules` в начале `.gitlab-ci.yml`:

```yaml
workflow:
  rules:
    - if: '$CI_COMMIT_REF_NAME == "master"'
      variables:
        ENV_NAME: dev2
        IMAGE_TAG_LATEST: "latest-dev2"
        DEPLOY_ENV_NAME: "dev2"
        BUILD_MODE: "development"
    - if: '$CI_COMMIT_REF_NAME == "envs/prod5"'
      variables:
        ENV_NAME: prod5
        IMAGE_TAG_LATEST: "latest-prod5"
        DEPLOY_ENV_NAME: "prod5"
        BUILD_MODE: "production"
    - when: always   # ← любая другая ветка: пайплайн запустится,
                     #   но переменные будут пустыми
```

`when: always` в конце важен — без него пайплайн вообще не запустился бы на feature-ветках. С ним он запускается, но большинство джобов не сработают (потому что `if: "$DEPLOY_ENV_NAME"` будет false для пустой переменной).

**Отображение:**
```
master       → dev2   (development, latest-dev2)
envs/dev2    → dev2   (development, latest-dev2)  ← дублирует master
envs/stage1  → stage1 (production, latest-stage1)
envs/stage3  → stage3 (production, latest-stage3)
envs/prod1   → prod1  (production, latest-prod1)
envs/prod4   → prod4  (production, latest-prod4)
envs/prod5   → prod5  (production, latest-prod5)
```

---

### Environment-per-folder (папка = среда)

**Что это:** для каждой среды — отдельная папка с манифестами. Файлы называются одинаково, но содержат разные значения (образы, URL, реплики, секреты).

```
infrastructure/dev2/backend-api-v2.yml   ← replicas: 2, image: :latest-dev2
infrastructure/prod5/backend-api-v2.yml  ← replicas: 5, image: :latest-prod5
```

**Сравнение с Helm** — Helm решает ту же задачу через шаблонизацию:
```
Helm:         один шаблон + values-dev.yaml + values-prod.yaml
Этот репо:    отдельный YAML на каждую среду (значения захардкожены)
```
Helm гибче при большом числе сред, но требует изучения синтаксиса шаблонов. Плоские YAML-файлы проще читать и дебажить.

---

### SRP (Single Responsibility Principle) на CI-джобы

**Что такое SRP:** принцип из SOLID — каждый модуль/класс/функция должны иметь **одну причину для изменения** и **одну ответственность**. Придуман Робертом Мартином. Применяется не только к коду, но и к CI-джобам.

**Проблема без SRP:** один джоб делает всё — тестирует, собирает, деплоит. Если сломается деплой, придётся перезапускать тесты и сборку заново. Если изменился только YAML-манифест — всё равно запускается полная сборка образа.

**Как реализован здесь:** 3 джоба на каждый сервис, каждый с **одной ответственностью** и **одним триггером**:

```
build-<app>    ответственность: создать Docker-образ
               триггер: изменился исходный код
               причина изменить джоб: изменился способ сборки образа

deploy-<app>   ответственность: синхронизировать манифест с кластером
               триггер: изменился YAML-файл
               причина изменить джоб: изменился способ применения манифеста

upgrade-<app>  ответственность: перезапустить поды с новым образом
               триггер: изменился исходный код (после build)
               причина изменить джоб: изменился способ перекатки подов
```

**Практическое следствие:** если сломался `upgrade` (например недоступен кластер), `build` уже выполнен и образ в registry — при повторном запуске не нужно пересобирать образ.

---

### Template Method (шаблонный метод)

**Что такое паттерн Template Method:** классический паттерн из ООП (GoF). Абстрактный класс определяет скелет алгоритма, оставляя заполнение деталей подклассам. Подклассы переопределяют только то, что отличается.

**В ООП:**
```python
class BuildTemplate:           # абстрактный класс
    def run(self):
        self.login()           # одинаково для всех
        self.build(self.app_name())  # вызывает абстрактный метод
        self.push()            # одинаково для всех

    def app_name(self): ...    # ← переопределяется в подклассе

class BuildBackendApi(BuildTemplate):
    def app_name(self): return "backend-api-v2"

class BuildWorker(BuildTemplate):
    def app_name(self): return "worker-v1"
```

**В GitLab CI через `extends:`:**
```yaml
# "Абстрактный класс" — шаблон
.build-template:
  script:
    - docker login ...           # одинаково для всех
    - docker-compose build $APP_NAME  # использует "абстрактный метод"
    - docker-compose push $APP_NAME

# "Подклассы" — конкретные джобы, переопределяют только APP_NAME
build-backend-api-v2:
  extends: .build-template
  variables:
    APP_NAME: "backend-api-v2"

build-worker-v1:
  extends: .build-template
  variables:
    APP_NAME: "worker-v1"
```

**Без этого паттерна** пришлось бы скопировать ~20 строк скрипта 8 раз. Изменение одной строки (например версия Docker) потребовало бы правки в 8 местах.

---

### Event-Driven Deployment

**Что это:** деплой происходит не по расписанию и не вручную, а как **реакция на событие** — изменение конкретного файла в Git.

```yaml
deploy-backend-api-v2:
  rules:
    - if: "$DEPLOY_ENV_NAME"       # условие 1: мы на нужной ветке
      changes:                     # условие 2: изменился этот файл
        - "infrastructure/**/backend-api-v2.yml"
        - ".gitlab-ci.yml"
```

Оба условия должны выполниться одновременно. `**` — glob-паттерн, работает для любой среды (`dev2`, `prod5`, ...).

---

## Три джоба на каждый сервис

```
build-<app>
  триггер:   изменился исходный код (golang/**/* и т.д.)
  depends:   needs: ["unittest-<lang>"]  ← только после тестов
  делает:    docker-compose build → docker-compose push → cr.yandex/:latest-dev2
  результат: новый образ в registry, в кластере ничего не изменилось

deploy-<app>
  триггер:   изменился YAML (infrastructure/dev2/app.yml или .gitlab-ci.yml)
  depends:   needs: []  ← запускается сразу, никого не ждёт
  делает:    kubectl apply -f infrastructure/dev2/app.yml
  результат: k8s применил новую конфигурацию (реплики, env, порты...)
             образ НЕ пересобирается, поды НЕ перезапускаются

upgrade-<app>
  триггер:   изменился исходный код (тот же что и build)
  depends:   needs: ["build-<app>"]  ← только после сборки
  делает:    kubectl rollout restart deployment/app
  результат: k8s перезапустил поды → скачал свежий образ из registry
             YAML-манифест НЕ читается из репозитория
```

---

## Как k8s подхватывает новый образ

Само приложение **не знает** что вышла новая версия — это оркестрирует Kubernetes.

```
1. CI пушит новый образ с тем же тегом:
   docker push cr.yandex/.../backend-api-v2:latest-dev2
   (тег не изменился, но содержимое образа новое)

2. CI говорит k8s перезапустить деплоймент:
   kubectl rollout restart deployment/backend-api-v2

3. k8s запускает Rolling Update:
   → создаёт новый под параллельно со старым
   → смотрит в манифест: imagePullPolicy: Always
   → идёт в registry, скачивает свежий :latest-dev2
   → дожидается Ready нового пода
   → убивает старый под
   → повторяет для каждого пода (по одному)

4. В итоге все поды работают на новом коде
   Старый код работал всё время обновления → Zero Downtime
```

**`imagePullPolicy: Always`** — без него k8s использовал бы кешированный образ с тем же тегом и ничего бы не обновилось.

**Почему `rollout restart` а не "само"** — тег `:latest-dev2` не изменился. k8s следит за изменением тега или digest образа. Если ни то ни другое не изменилось — k8s считает что обновлять нечего.

**Rolling Update vs Recreate** — стратегия по умолчанию в Deployment. Альтернатива `Recreate` — убить всё сразу, поднять заново (downtime). Rolling лучше для production.

**Почему тег `latest-dev2` а не `latest`** — один тег `latest` на все среды означает что prod может скачать dev-образ. Тег `latest-<среда>` изолирует среды друг от друга.

---

## Связь .gitlab-ci.yml ↔ infrastructure/

Ключевая строка в `.deploy-template`:

```bash
kubectl apply -f infrastructure/$DEPLOY_ENV_NAME/$APP_NAME.yml
```

Обе переменные подставляются автоматически из разных мест:
- `DEPLOY_ENV_NAME` — из `workflow.rules` по имени ветки (`master` → `dev2`)
- `APP_NAME` — из `variables:` конкретного джоба (`"backend-api-v2"`)

```
ветка master  +  APP_NAME=backend-api-v2  →  infrastructure/dev2/backend-api-v2.yml
ветка envs/prod5  +  APP_NAME=worker-v1   →  infrastructure/prod5/worker-v1.yml
```

**Триггер деплоя** задаётся через `changes:` в том же джобе:
```yaml
rules:
  - changes:
      - "infrastructure/**/backend-api-v2.yml"
```
Glob `**` означает "любая вложенность" — сработает и для `dev2/`, и для `prod5/`.

---

## Что запускается при каждом сценарии

```
Сценарий А: пуш Go-кода (golang/**/* изменился)
──────────────────────────────────────────────────
STAGE unittest:  unittest-golang
STAGE build:     build-backend-api-v2  (needs: unittest ✓)
                 build-worker-v1       (needs: unittest ✓)  ← параллельно
STAGE deploy:    upgrade-backend-api-v2 (needs: build ✓)
                 upgrade-worker-v1      (needs: build ✓)
STAGE e2etest:   e2etest-golang

deploy-backend-api-v2 — НЕ запускается (YAML не менялся)


Сценарий Б: пуш YAML-манифеста infrastructure/dev2/backend-api-v2.yml
──────────────────────────────────────────────────
STAGE deploy:  deploy-backend-api-v2  (needs: [] → сразу)

Всё остальное — НЕ запускается (код не менялся)


Сценарий В: пуш кода + YAML одновременно
──────────────────────────────────────────────────
STAGE unittest:  unittest-golang
STAGE build:     build-backend-api-v2
STAGE deploy:    upgrade-backend-api-v2  (ждёт build)
                 deploy-backend-api-v2   (needs: [] → параллельно с upgrade)
STAGE e2etest:   e2etest-golang
```

---

## Шаблоны в .gitlab-ci.yml — детально

Имена начинаются с `.` — GitLab не запускает их как самостоятельные джобы.

| Шаблон | Что делает | Ключевой механизм |
|--------|-----------|-------------------|
| `.unittest-template` | `cd $LANGUAGE && make test` | параметр `$LANGUAGE` |
| `.build-template` | docker login → build → push | DinD, `$APP_NAME`, `$BUILD_MODE` |
| `.deploy-template` | kubectl apply манифест | `$DEPLOY_ENV_NAME` + `$APP_NAME` |
| `.upgrade-template` | rollout restart Deployment | `kubectl config use-context` |
| `.upgrade-db-template` | rollout restart StatefulSet | StatefulSet вместо Deployment (для БД) |
| `.job-template` | apply Job → polling → логи | while-loop, exit код пода |
| `.e2etest-template` | make testall + артефакты | JUnit report, coverage regex |
| `.freezedeployment` | блок деплоя в outage window | `CI_DEPLOY_FREEZE` переменная |

**`.upgrade-template` vs `.upgrade-db-template`** — разница в типе k8s-ресурса:
- `Deployment` — для stateless приложений (backend, worker, frontend). Поды взаимозаменяемы.
- `StatefulSet` — для stateful (PostgreSQL). Каждый под имеет стабильный идентификатор и отдельный том. `rollout restart statefulset` обновляет их по одному в строгом порядке.

**`before_script`** в пайплайне — выполняется перед каждым джобом:
```yaml
before_script:
  - echo "Branch: '$CI_COMMIT_BRANCH', tag: '$IMAGE_TAG_LATEST', cluster: '$DEPLOY_ENV_NAME'"
```
Полезно при дебаггинге — в логах каждого джоба видно в каком контексте он запустился.

---

## docker-compose как абстракция сборки

В `.build-template` используется `docker-compose`, а не напрямую `docker build`:

```bash
docker-compose --env-file ./envs/build.env build --build-arg BUILD_MODE=$BUILD_MODE $APP_NAME
docker-compose --env-file ./envs/build.env push $APP_NAME
```

**Зачем docker-compose а не просто docker:**
- `docker-compose.build.yml` хранит имя образа, registry URL, теги — не нужно дублировать их в CI
- `--env-file build.env` позволяет передавать build-аргументы через файл, а не хардкодить в CI
- `$APP_NAME` — docker-compose знает какой сервис собирать по имени из `docker-compose.yml`

---

## Структура infrastructure/

```
infrastructure/
├── dev2/                        ← среда разработки
│   ├── backend-api-v2.yml       ← Deployment: Go API
│   ├── worker-v1.yml            ← Deployment: воркеры
│   ├── frontend-v3.yml          ← Deployment: фронтенд
│   ├── telegram-bot-v1.yml
│   ├── py-task-executor-v1.yml
│   ├── redis.yml                ← PV + PVC + Deployment + Service + ConfigMap
│   ├── pgbouncer.yml            ← пул соединений к PostgreSQL
│   ├── ingress.yml              ← nginx Ingress + cert-manager (TLS)
│   ├── configs-and-secrets.yml  ← Secrets, ServiceAccount, imagePullSecret
│   ├── dbmate.yml               ← k8s Job: миграции БД
│   ├── temporal.yml             ← Temporal workflow orchestration
│   ├── minio.yml                ← S3-совместимое хранилище
│   ├── grafana.yml              ← дашборды (namespace: observability)
│   ├── itv-supabase-db.yml      ← PostgreSQL (StatefulSet)
│   ├── itv-supabase-auth.yml    ← GoTrue (авторизация)
│   ├── itv-supabase-kong.yml    ← Kong API gateway
│   ├── itv-supabase-meta.yml    ← Supabase Meta API
│   ├── itv-supabase-storage.yml
│   ├── itv-supabase-studio.yml  ← веб-UI Supabase
│   └── manual/                  ← ручная настройка кластера (один раз)
│       ├── gitlab-agent.yaml         ← helm install k8s-agent
│       ├── gitlab-runner-chart-values.yaml
│       ├── kube-prom-stack-values.yaml
│       └── loki-values.yaml
├── stage1/    ← те же имена файлов, другие значения
├── prod1/
└── prod5/
```

**Паттерн multi-document YAML** — один файл может содержать несколько k8s-ресурсов разделённых `---`. Например `redis.yml` содержит: PersistentVolume + PersistentVolumeClaim + Deployment + Service + ConfigMap. Логически связанные ресурсы держатся вместе.

**ConfigMap как shared config** — один сервис публикует своё подключение, другие читают:
```yaml
# redis.yml публикует:
kind: ConfigMap
data:
  REDIS_URL: "redis:6379"

# backend-api-v2.yml читает:
env:
  - name: REDIS_ADDR
    valueFrom:
      configMapKeyRef:
        name: redis
        key: REDIS_URL
```

---

## Как CI попадает в Kubernetes (GitLab Agent)

Прямой доступ к kube API небезопасен (нужно открывать порт наружу, хранить kubeconfig). Используется GitLab Kubernetes Agent:

```
GitLab (облако)
    ▲
    │  wss:// WebSocket — кластер сам устанавливает исходящее соединение
    │  Кластер не нужно открывать наружу
    │
k8s-agent Pod (внутри кластера dev2)
    │
    ▼
Kubernetes API Server (localhost внутри кластера)
```

**Как это работает для CI:**
1. Runner запускает джоб
2. GitLab пробрасывает kubeconfig через агента (runner видит контекст `platform/core:k8s-agent-dev2`)
3. Все `kubectl` команды идут через этот туннель

```bash
kubectl config use-context platform/core:k8s-agent-dev2
# формат: <gitlab-группа>/<репозиторий>:k8s-agent-<ENV_NAME>
```

**Конфигурация агента** в репо: `.gitlab/agents/k8s-agent-dev2/config.yaml`
```yaml
ci_access:
  groups:
    - id: platform   # GitLab-группа которой разрешён доступ к этому кластеру
```

Для каждой среды — свой агент: `k8s-agent-dev2`, `k8s-agent-stage1`, `k8s-agent-prod5`.
Установка агента: `infrastructure/dev2/manual/gitlab-agent.yaml` (один раз, через Helm).

---

## Роутинг джобов на нужный runner

Все джобы в `.gitlab-ci.yml` имеют:
```yaml
tags:
  - ${ENV_NAME}
```

`ENV_NAME` подставляется из `workflow.rules` по ветке. Каждый runner зарегистрирован с тегом = имя среды.

```
push master → ENV_NAME=dev2 → все джобы уходят на Runner-dev2 → kubectl k8s-agent-dev2
push envs/prod5 → ENV_NAME=prod5 → Runner-prod5 → kubectl k8s-agent-prod5
```

Runner'ы установлены прямо в k8s-кластерах (Helm `gitlab/gitlab-runner`, конфиги в `manual/`).

**Важно:** runner с тегом `dev2` имеет доступ только к агенту `k8s-agent-dev2`. Джоб dev2 физически не может задеплоить что-то в prod5 — он просто не подключится к тому контексту.

---

## Docker в CI (DinD — Docker-in-Docker)

Docker в CI **не запускает приложение** — только собирает образ и пушит.

```
Runner Pod (в k8s)
├── Job-контейнер: docker:20.10.18
│     ← здесь выполняется скрипт CI
│     ← docker-compose build ... (команды идут к демону ниже)
│     ← docker-compose push ...
│
└── Service-контейнер: docker:20.10.18-dind
      ← это Docker-демон (dockerd)
      ← слушает на tcp://docker:2376 (TLS)
      ← получает команды от job-контейнера
      ← здесь живёт собранный образ, нигде не запускается
```

**Почему TCP а не сокет** — job и dind это разные контейнеры в одном поде. Docker-сокет `/var/run/docker.sock` не шарится между контейнерами. Поэтому демон слушает на TCP.

**TLS** — `DOCKER_TLS_CERTDIR`, `DOCKER_TLS_VERIFY`, `DOCKER_CERT_PATH` — сертификаты для безопасного соединения между job и dind. Автоматически генерируются при старте dind.

После пуша CI-контейнер умирает:
```
CI job → умер    Yandex Registry → хранит образ    k8s Pod → работает
```

---

## Аутентификация в registry

`key.json` — файл сервисного аккаунта Yandex Cloud с правами `container-registry.images.pusher/puller`.
Хранится в **GitLab CI/CD Variables** как тип File (не в репо).

```bash
# В каждом build/deploy-джобе:
cat key.json | docker login --username json_key --password-stdin cr.yandex
```

**Два уровня аутентификации:**
1. **CI → Registry** — через `key.json` при push/pull в скриптах
2. **k8s → Registry** — через `imagePullSecret` в кластере (задаётся в `configs-and-secrets.yml`, привязывается к ServiceAccount `default`)

Второй уровень нужен потому что `kubectl rollout restart` заставляет k8s самостоятельно идти в registry за образом — CI уже не участвует в этот момент.

---

## Helm — где используется

Helm **не используется** для деплоя собственных сервисов. Только для кластерной инфраструктуры в `manual/` (один раз вручную):

```bash
helm upgrade --install k8s-agent-dev2 gitlab/gitlab-agent ...
helm upgrade kube-prom-stack prometheus-community/kube-prometheus-stack ...
helm upgrade loki grafana/loki-stack ...
```

**Helm-подобная структура без Helm:** папки сред воспроизводят идею `values.yaml`, но без шаблонизации. Каждый YAML захардкожен под свою среду. Это проще для чтения и дебага, но сложнее поддерживать при большом числе сред — при изменении структуры манифеста нужно менять файл в каждой папке вручную.

---

## Тонкости синтаксиса

### `needs` — явные зависимости между джобами

```yaml
build-backend-api-v2:
  needs: ["unittest-golang"]    # ждать юнит-тесты

upgrade-backend-api-v2:
  needs: ["build-backend-api-v2"]  # ждать сборку

deploy-backend-api-v2:
  needs: []   # НЕ ждать никого
```

Без `needs: []` GitLab по умолчанию ждёт **все джобы предыдущей стадии**. Для `deploy` это критично: если `build` не запустился (код не менялся), GitLab посчитает что предыдущая стадия не завершена и не запустит `deploy`.

### `rules` vs `only` — два синтаксиса

```yaml
# Старый (only) — только path-фильтр, нельзя комбинировать с if
build-worker-v1:
  only:
    changes: ["golang/**/*"]

# Новый (rules) — if + changes вместе, первое совпавшее правило побеждает
deploy-backend-api-v2:
  rules:
    - if: "$DEPLOY_ENV_NAME"          # сначала проверяем среду
      changes:
        - "infrastructure/**/backend-api-v2.yml"
        - ".gitlab-ci.yml"
      when: always
```

В файле смешаны оба стиля — legacy от старых джобов.

### `when` — условие запуска

```yaml
when: on_success  # (дефолт) — только если предыдущие шаги прошли
when: always      # — всегда, даже если что-то упало
when: manual      # — только нажатием кнопки в GitLab UI
when: never       # — никогда (используется для отключения правила)
```

`when: always` у инфра-джобов (redis, pgbouncer...) — чтобы деплойились независимо от других джобов.

### `interruptible`

```yaml
.e2etest-template:
  interruptible: true   # этот джоб можно прервать при auto_cancel
```

Только e2etest помечен как прерываемый. Build и deploy не прерываются при новом пуше — пересобирать образ это дорого, а деплой нельзя прерывать на полуслове.

### `environment`

```yaml
.upgrade-template:
  environment:
    name: $DEPLOY_ENV_NAME/$APP_NAME
```

Создаёт GitLab Environment — раздел в UI где видны последние деплои каждого сервиса в каждую среду, история деплоев, кто и когда деплоил.

---

## Специальные случаи

**`run-migrations`** — единственный `.job-template`:
```
удалить старый Job (kubectl delete job --ignore-not-found)
→ apply нового Job из dbmate.yml
→ polling в цикле (sleep 10) пока pod.status == Succeeded/Failed
→ вывести kubectl logs пода
→ exit 0 (Succeeded) или exit 1 (Failed)
```
CI блокируется до завершения миграций и видит их результат прямо в логах пайплайна.

**`deploy-grafana`** — переопределяет `KUBECTL_NAMESPACE: "observability"`. Единственный сервис не в namespace `default`. Namespace'ы изолируют ресурсы внутри одного кластера — Grafana живёт в `observability` вместе с Prometheus и Loki.

**`build-android-v1`** — `when: manual`, нет `deploy`. Собирает APK, прикладывает как артефакт (хранится 2 дня), скачивается из GitLab UI.

**`apply-blueprint`** — `when: manual`. Ручное применение Supabase blueprint (схема данных/RLS-правила). Не автоматизировано — требует осознанного запуска.

**`lint-golang`** — стоит в стадии `e2etest` вместо `unittest`. Комментарий TODO: перенести. Сейчас там из-за того что линтер ломал пайплайн при активной разработке — в e2etest он не блокирует сборку.

**`pack-migrations`** — собирает Docker-образ dbmate с SQL-миграциями. Отдельный джоб от `run-migrations` — сначала упакуй, потом запусти.

**Закомментированные джобы** — `deploy-itv-supabase-realtime/rest`: компоненты временно отключены, манифесты в `infrastructure/` сохранены.

---

## Оптимизации CI

**`GIT_DEPTH: 1`** — shallow clone, только последний коммит без истории. Ускоряет старт.

**`GIT_STRATEGY: fetch`** — не `clone` (скачать заново), а `fetch` поверх существующей копии на runner'е. Значительно быстрее на больших репо.

**`auto_cancel.on_new_commit: interruptible`** — новый пуш отменяет старый пайплайн для джобов с `interruptible: true`. Экономит ресурсы runner'ов.

**`needs`** для параллелизма — все `build-*` джобы одной стадии запускаются параллельно (каждый ждёт только свой unittest, не чужой).

---

## Подводные камни

**`imagePullPolicy: Always` обязателен** — иначе k8s кеширует образ по тегу. При `IfNotPresent` (дефолт) с мутирующим тегом `latest-dev2` новый код никогда не приедет без ручного удаления кеша на ноде.

**`imagePullSecrets` в ServiceAccount** — если не настроен, поды не смогут скачать образ из приватного registry. Ошибка: `ErrImagePull` / `ImagePullBackOff`. Секрет задаётся в `configs-and-secrets.yml`, привязывается к ServiceAccount `default` — применяется ко всем подам namespace автоматически.

**`needs: []` обязателен у `deploy-*`** — без него deploy ждёт все джобы стадии build. Если build не запустился (код не менялся), GitLab решит что стадия "не завершена" и не запустит deploy вовсе.

**`only` нельзя комбинировать с `if`** — если нужно `if: "$DEPLOY_ENV_NAME"` + `changes:`, джоб обязательно должен использовать `rules:`, а не `only:`.

**Тег среды на runner'е** — runner с тегом `dev2` не сможет использовать контекст `k8s-agent-prod5`. Это одновременно защита (нельзя случайно задеплоить в прод с dev-ветки) и ограничение (нельзя запустить deploy в prod5 с runner'а dev2).

---

## Дебаггинг

| Проблема | Что проверить |
|----------|--------------|
| Джоб не запустился | `workflow.rules` (ветка подходит?), `changes` (нужный файл менялся?), `if: "$DEPLOY_ENV_NAME"` (переменная не пустая?), `needs` (зависимость выполнена?) |
| kubectl: no context | Pod k8s-agent запущен? → `kubectl get pods -n gitlab-agent-k8s-agent-dev2` |
| Образ не обновился | `upgrade-*` запустился? `imagePullPolicy: Always` стоит? `build` не упал? |
| Pod не может скачать образ | `imagePullSecrets` настроены в `configs-and-secrets.yml`? Применены через ServiceAccount? |
| Деплой заблокирован | `CI_DEPLOY_FREEZE` выставлена в GitLab → запустить вручную через UI |
| Миграции упали | Логи прямо в CI в шаге `run-migrations` (job-template их печатает) |
| deploy запустился без build | Это нормально — разные триггеры. deploy на YAML-изменение, build на код-изменение |

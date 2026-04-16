# Helm и CI/CD: GitLab + Kubernetes

---

## 1. Helm

> Аналогия: Helm — это **менеджер пакетов** (как apt или npm) для Kubernetes.
> Chart — пакет (как deb-файл), Release — конкретная установка Chart в кластер.
> Вместо 10 отдельных YAML-файлов — один `helm install`.

### Ключевые понятия

| Понятие | Что это | Аналогия |
|---|---|---|
| Chart | Набор YAML-шаблонов + values.yaml | npm-пакет / deb-пакет |
| Release | Конкретная установка Chart с определёнными values | Установленное приложение |
| Repository | Каталог Charts | npm registry |
| Values | Параметры установки (image.tag, replicas...) | Конфигурационный файл |

### Основные команды

```bash
# Репозитории
helm repo add myrepo https://charts.example.com
helm repo update
helm search repo myrepo

# Установка / обновление
helm install myapp myrepo/myapp -f values.yaml -n production
helm upgrade --install myapp myrepo/myapp -f values.yaml -n production  # idempotent (CI/CD)
helm rollback myapp 1 -n production

# Управление
helm list -n production
helm get values myapp -n production
helm get values myapp -n production --all   # включая дефолтные
helm uninstall myapp -n production
```

### install vs upgrade --install

| | `helm install` | `helm upgrade --install` |
|---|---|---|
| Release уже есть | Ошибка | Обновит |
| Release нет | Установит | Установит |
| Idempotent | Нет | Да ✅ |
| Применение | Первый раз вручную | **Всегда в CI/CD** |

### Helm hooks

> Аналогия: hooks — это **задачи до/после переезда**.
> pre-install: "упаковать вещи перед переездом" (миграции БД перед деплоем).
> post-install: "распаковать вещи после переезда" (smoke-тест после деплоя).

| Hook | Когда запускается |
|---|---|
| `pre-install` | До создания ресурсов |
| `post-install` | После создания всех ресурсов |
| `pre-upgrade` | До обновления |
| `pre-delete` | До удаления |

---

## 2. GitLab Runner в Kubernetes

```bash
# Добавить helm-репозиторий GitLab
helm repo add gitlab https://charts.gitlab.io
helm repo update gitlab

# Создать namespace
kubectl create ns gitlab-runner

# Установить Runner (токен берём из GitLab → Settings → CI/CD → Runners)
helm upgrade --install -n gitlab-runner gitlab-runner gitlab/gitlab-runner -f values.yaml
```

---

## 3. Пример .gitlab-ci.yml

```yaml
stages:
  - build
  - test
  - deploy

variables:
  IMAGE: $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

build:
  stage: build
  script:
    - docker build -t $IMAGE .
    - docker push $IMAGE

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
  only:
    - main
```

---

## 4. Kaniko — сборка образов без Docker daemon

> Аналогия: обычная сборка через Docker требует "ключи от сервера" (docker socket).
> Kaniko собирает образы изнутри Pod без привилегий — безопаснее docker-in-docker (DinD).

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

### Docker daemon vs Kaniko vs DinD

| | Docker daemon | DinD | Kaniko |
|---|---|---|---|
| Нужен docker socket | Да (небезопасно) | Нет | Нет |
| Privileged mode | Да | Да | Нет ✅ |
| Безопасность | Низкая | Средняя | Высокая |
| Применение | Локально | Legacy CI | Production CI/CD в K8s |

---

## 5. GitOps

> Аналогия: GitOps — это когда Git — это **единственная правда**.
> Хочешь изменить что-то в production — делаешь PR в Git.
> Оператор (ArgoCD/Flux) сам видит изменения и применяет их.

### GitLab CI vs ArgoCD

| | GitLab CI (push-модель) | ArgoCD (pull-модель) |
|---|---|---|
| Кто деплоит | Pipeline явно вызывает helm/kubectl | Оператор в кластере следит за Git |
| Источник правды | Pipeline конфиг | Git репозиторий |
| Multi-cluster | Сложнее | Проще |
| Audit | Через pipeline logs | Через ArgoCD UI / Git history |
| Рекомендация | Простые проекты | Серьёзный production, много кластеров |

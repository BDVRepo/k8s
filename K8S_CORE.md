# Kubernetes — Архитектура и Основные объекты

> Аналогия: Kubernetes — это **логистическая компания** (DHL, FedEx).
> - Кластер — вся компания целиком.
> - Control Plane — головной офис (принимает заявки, планирует, контролирует).
> - Worker Node — склад/площадка, где реально стоят и едут грузовики (контейнеры).
> - Pod — грузовик с одним или несколькими контейнерами на борту.

---

## 1. Архитектура кластера

### Control Plane (Master-ноды) — "Головной офис"

| Компонент | Роль | Аналогия |
|---|---|---|
| `kube-apiserver` | Единая точка входа. Принимает запросы kubectl/REST, валидирует, записывает в etcd | Ресепшн, куда приходят все заявки |
| `etcd` | Распределённое key-value хранилище. Хранит **всё** состояние кластера | Главный реестр компании (база данных) |
| `kube-scheduler` | Выбирает, на какую Node поставить Pod (смотрит на ресурсы, affinity, taints) | Логист-планировщик, который решает на какой склад отправить груз |
| `kube-controller-manager` | Запускает контроллеры в бесконечном цикле reconcile (desired vs actual state) | Менеджер, который следит чтобы всё было "как договорились" |

### Worker Node — "Склады/площадки"

| Компонент | Роль | Аналогия |
|---|---|---|
| `kubelet` | Агент на каждой ноде. Получает PodSpec, запускает контейнеры через CRI | Прораб на площадке |
| `kube-proxy` | Управляет сетевыми правилами (iptables/ipvs) для Services | Диспетчер маршрутизации трафика |
| Container Runtime | Запускает контейнеры (`containerd`, `CRI-O`) | Сам механизм (кран, погрузчик) |

### Как работает `kubectl apply`

```
kubectl apply -f pod.yaml
        ↓
kube-apiserver (валидация + авторизация RBAC)
        ↓
etcd (сохранение desired state)
        ↓
kube-scheduler (выбирает Node для Pod)
        ↓
kubelet на выбранной Node (скачивает образ, запускает контейнер)
        ↓
Pod Running
```

### Reconcile loop — принцип работы контроллеров

> Аналогия: термостат. Желаемая температура = 22°. Если реальная температура отличается — включается обогрев/кондиционер.
> Так и K8s: желаемое состояние (в etcd) vs реальное → контроллер корректирует.

```
desired state (etcd): replicas=3
actual state (ноды):  replicas=2   → controller-manager создаёт +1 Pod
```

---

## 2. Кто кого создаёт и контролирует

| Сущность | Для чего нужна | Кто создаёт | Кто контролирует |
|---|---|---|---|
| `Deployment` | Stateless-приложения, rolling update/rollback | Пользователь | `deployment-controller` |
| `ReplicaSet` | Держать N одинаковых Pod | `Deployment` | `replicaset-controller` |
| `StatefulSet` | Stateful-приложения (БД): стабильные имена, PVC | Пользователь | `statefulset-controller` |
| `DaemonSet` | По одному Pod на каждую ноду | Пользователь | `daemonset-controller` |
| `Job` | Одноразовая задача до успешного завершения | Пользователь или `CronJob` | `job-controller` |
| `CronJob` | Job по расписанию | Пользователь | `cronjob-controller` |
| `Pod` | Минимальная единица запуска контейнеров | Обычно контроллеры | `kubelet` на Node |
| `Service` | Стабильный доступ к Pod | Пользователь | `endpoints-controller`, `kube-proxy` |
| `Node` | Сервер, где запускаются Pod | Регистрируется `kubelet` | `node-controller` |

### Цепочка зависимостей

```
Deployment → ReplicaSet → Pod
CronJob    → Job        → Pod
DaemonSet              → Pod (по одному на каждую Node)
StatefulSet            → Pod (с персистентным именем и PVC)

Pod → Node (выбирает kube-scheduler)
Pod → Container (запускает kubelet через container runtime)
```

---

## 3. Основные объекты

### 3.1 Namespace

> Аналогия: Namespace — это **отдел в компании** (dev, staging, production).
> Сотрудники одного отдела знают друг друга по коротким именам,
> но чтобы обратиться в другой отдел — нужно полное имя.

```bash
kubectl get namespaces
kubectl create ns my-namespace
kubectl -n my-namespace get pods
```

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: home-dev
```

---

### 3.2 Pod

> Аналогия: Pod — это **грузовик** с одним или несколькими контейнерами.
> Все контейнеры в Pod: один IP, одни тома, одна нода.
> Pod — минимальная единица в K8s, как атом в химии.

**Жизненный цикл Pod:**

| Статус | Значение |
|---|---|
| `Pending` | Принят, ждёт планирования или загрузки образа |
| `Running` | Хотя бы один контейнер работает |
| `Succeeded` | Все контейнеры завершились с кодом 0 |
| `Failed` | Хотя бы один контейнер завершился с ненулевым кодом |
| `Unknown` | Нет связи с нодой |

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
      resources:
        requests:
          cpu: "100m"
          memory: "64Mi"
        limits:
          cpu: "200m"
          memory: "128Mi"
```

```bash
kubectl apply -f pod.yaml
kubectl -n home-dev get pods
kubectl -n home-dev describe pod simple-web
kubectl -n home-dev logs simple-web
kubectl -n home-dev port-forward pod/simple-web 8080:8080
kubectl delete -f pod.yaml
```

**Почему Pod не запускают напрямую в production?**
Pod не восстанавливается сам после падения ноды. Для этого нужны контроллеры (Deployment, ReplicaSet).

---

### 3.3 ReplicaSet

> Аналогия: ReplicaSet — это **контракт на поставку**: "всегда должно быть 3 грузовика в рейсе".
> Если один сломался — немедленно запускают замену.

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

**Ограничение:** ReplicaSet не умеет обновлять уже запущенные Pods при смене образа. Нужен Deployment.

---

### 3.4 Deployment

> Аналогия: Deployment — это **менеджер автопарка**.
> Он управляет ReplicaSet'ами и умеет обновлять грузовики один за другим (rolling update)
> без остановки всей логистики.

```
Deployment → ReplicaSet (v2, 3 pods running)
           → ReplicaSet (v1, 0 pods, хранится для rollback)
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goapp-deployment
  namespace: web-app-stage
  labels:
    app: goapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: goapp
  template:
    metadata:
      labels:
        app: goapp
    spec:
      containers:
        - name: web
          image: nginx
          ports:
            - containerPort: 8070
```

```bash
kubectl -n home-dev apply -f deployment.yaml
kubectl -n home-dev rollout status deployment/go-http-server
kubectl -n home-dev rollout history deployment/go-http-server
kubectl -n home-dev rollout undo deployment/go-http-server
kubectl -n home-dev rollout undo deployment/go-http-server --to-revision=2
kubectl -n home-dev rollout restart deployment/go-http-server
kubectl -n home-dev scale deployment/go-http-server --replicas=5
kubectl diff -f deployment.yaml   # что изменится перед apply
```

---

### 3.5 DaemonSet

> Аналогия: DaemonSet — это **охранник на каждом складе**.
> Сколько бы складов ни открылось — на каждом должен быть свой охранник.

**Типичное применение:**
- Сбор логов (Fluentd, Filebeat)
- Мониторинг (node-exporter)
- Сетевые плагины (Calico, Flannel)

При добавлении новой Node Pod автоматически создаётся на ней.

---

### 3.6 Job

> Аналогия: Job — это **разовый подрядчик**. Нанял, выполнил задачу, ушёл.

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
          command: ["python", "-c", "from math import pi; print(f'{pi:.20f}')"]
      restartPolicy: Never   # Never или OnFailure
```

```bash
kubectl apply -f job.yaml
kubectl -n home-dev get jobs
kubectl -n home-dev logs <pod-name>
kubectl explain job.spec.backoffLimit
```

---

### 3.7 CronJob

> Аналогия: CronJob — это **расписание уборщика**. Каждую ночь в 03:00 приходит и делает бэкап.

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

```
*/2 * * * *   → каждые 2 минуты
0 * * * *     → каждый час
0 3 * * *     → каждый день в 03:00
0 3 * * 1     → каждый понедельник в 03:00
```

```bash
kubectl -n home-dev get cronjob
kubectl -n home-dev create job --from=cronjob/hello hello-manual  # запустить вручную
```

---

### 3.8 StatefulSet

> Аналогия: StatefulSet — это **нумерованные кабинеты в архиве** (кабинет №1, №2, №3).
> У каждого своё место, свои файлы, своё имя. Их нельзя просто переставить местами.

**Отличия от Deployment:**

| | Deployment | StatefulSet |
|---|---|---|
| Имена Pod | Случайные суффиксы (pod-xyz) | Стабильные (pod-0, pod-1) |
| Порядок запуска | Параллельно | Последовательно (0, 1, 2...) |
| Хранилище | Общий PVC | Отдельный PVC на каждую реплику |
| DNS-имя Pod | Нет | Стабильное (`pod-0.service.ns`) |
| Применение | Stateless (API, веб) | БД, Kafka, Zookeeper |

---

## 4. Сравнение workload-объектов

| Объект | Когда использовать | Аналогия |
|---|---|---|
| `Pod` | Отладка, тесты. Никогда в production напрямую | Разовая поездка на такси |
| `ReplicaSet` | Почти никогда напрямую — используйте Deployment | Автопарк без менеджера |
| `Deployment` | Stateless приложения (API, веб-сервисы) | Автопарк с менеджером |
| `DaemonSet` | Один Pod на каждой ноде (мониторинг, логи) | Охранник на каждом складе |
| `StatefulSet` | Stateful приложения (БД, Kafka) | Нумерованные архивные кабинеты |
| `Job` | Одноразовые задачи (миграции, расчёты) | Разовый подрядчик |
| `CronJob` | Периодические задачи (бэкапы, отчёты) | Расписание уборщика |

---

## 5. Labels, Annotations, Selectors

> Аналогия: Labels — это **стикеры на коробке** (app=web, env=prod).
> Selector — это **запрос**: "принеси мне все коробки с стикером app=web".

**Отличие Labels от Annotations:**

| | Labels | Annotations |
|---|---|---|
| Назначение | Идентификация и выборка | Метаданные для инструментов/людей |
| Используется в selector | Да | Нет |
| Длина значения | Короткое | Может быть длинной строкой |
| Пример | `app: myapp`, `env: prod` | `description: "main api service"` |

```bash
kubectl get pods -l app=myapp               # фильтрация по label
kubectl get pods -l app=myapp,env=prod      # несколько labels
```

---

## 6. Scheduling: Taints, Tolerations, Affinity

> Аналогия: Taint — это **табличка "Только для VIP"** на ноде.
> Toleration — это **VIP-пропуск** в Pod: "я могу сюда".
> Affinity — это **предпочтение**: "хочу быть рядом с GPU-нодой".

### Taints и Tolerations

```bash
# Добавить taint на ноду (GPU-только)
kubectl taint nodes gpu-node dedicated=gpu:NoSchedule
```

| Эффект | Поведение |
|---|---|
| `NoSchedule` | Не планировать новые Pods без toleration |
| `PreferNoSchedule` | Избегать, но не запрещать |
| `NoExecute` | Выселить уже запущенные Pods без toleration |

### nodeSelector (простой способ)

```yaml
nodeSelector:
  disktype: ssd
```

### Node Affinity (гибкий способ)

```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:  # жёсткое требование
      nodeSelectorTerms:
        - matchExpressions:
            - key: disktype
              operator: In
              values: ["ssd"]
```

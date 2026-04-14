# Kubernetes Cheatsheet — Быстрый справочник

> Компактный справочник на основе манифестов из этой папки.
> Для подготовки к собеседованию по инфраструктурным вопросам.

---

## Сравнительные таблицы

### Workload-объекты: ReplicaSet vs Deployment

| | ReplicaSet | Deployment |
|--|-----------|-----------|
| Управляет | Pods напрямую | ReplicaSet (→ Pods) |
| Rolling Update | Нет | Да |
| Rollback | Нет | Да (`rollout undo`) |
| История версий | Нет | Да (хранит старые RS) |
| Использовать | Почти никогда напрямую | Всегда для stateless |

### Job vs CronJob

| | Job | CronJob |
|--|-----|---------|
| Запуск | Один раз | По расписанию (cron) |
| Управляет | Pods | Job (→ Pods) |
| `backoffLimit` | Сколько раз повторить | Задаётся в jobTemplate |
| `restartPolicy` | `Never` или `OnFailure` | То же |
| Запустить вручную | `kubectl apply -f` | `kubectl create job --from=cronjob/name` |

### Service типы: ClusterIP vs NodePort vs LoadBalancer vs Ingress

| | ClusterIP | NodePort | LoadBalancer | Ingress |
|--|-----------|----------|--------------|---------|
| Доступен снаружи | Нет | Да (порт ноды) | Да (внешний IP) | Да (HTTP/HTTPS) |
| Диапазон портов | - | 30000–32767 | - | 80 / 443 |
| Протокол | TCP/UDP | TCP/UDP | TCP/UDP | HTTP/HTTPS |
| Маршрутизация | Только по IP | IP:порт | IP | hostname + path |
| Нужен контроллер | Нет | Нет | Облако / MetalLB | IngressController |
| Типичное применение | Между сервисами | Dev/тест | Production (облако) | Production (HTTP API) |

### ConfigMap vs Secret

| | ConfigMap | Secret |
|--|-----------|--------|
| Данные | Несекретные | Чувствительные |
| Хранение в etcd | Plaintext | Base64 (можно включить encryption at rest) |
| Типы | - | Opaque, tls, dockerconfigjson, ... |
| Подключение | envFrom, env, volume | envFrom, env, volume |
| Обновление в Pod | Через restart | Через restart (или при volume-монтировании — авто) |

### Probes

| | livenessProbe | readinessProbe | startupProbe |
|--|---------------|----------------|--------------|
| При провале | Перезапуск контейнера | Убрать из балансировки | Перезапуск контейнера |
| Цель | Обнаружить зависший контейнер | Не слать трафик на неготовый Pod | Защитить медленно стартующий контейнер |
| Проверки | httpGet, exec, tcpSocket | httpGet, exec, tcpSocket | httpGet, exec, tcpSocket |

---

## YAML-паттерны

### Namespace

```yaml
# ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: web-app-dev
```

### Pod

```yaml
# pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: simple-web
  namespace: web-app-dev
spec:
  containers:
    - name: web
      image: nginx
      imagePullPolicy: Always
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

### ReplicaSet

```yaml
# rs.yaml (→ rs-test/rs.yaml)
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: goapp-replicaset
  namespace: rs-test
spec:
  replicas: 2
  selector:
    matchLabels:
      app: goapp           # ← должен совпадать с labels в template
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

### Deployment

```yaml
# deployment.yaml
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

### Deployment с RollingUpdate

```yaml
# rolling.yaml
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1   # макс. недоступных Pods во время обновления
      maxSurge: 1         # макс. дополнительных Pods сверх replicas
  selector:
    matchLabels:
      app: goapp
  template:
    ...
```

### Job

```yaml
# job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi-job
spec:
  backoffLimit: 4
  template:
    spec:
      containers:
        - name: pi
          image: python:3.9
          command: ["python", "-c", "from math import pi; print(f'{pi:.20f}')"]
      restartPolicy: Never
```

### CronJob

```yaml
# cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: pi-cronjob
spec:
  schedule: "*/2 * * * *"   # каждые 2 минуты
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: pi
              image: python:3.9
              command: ["python", "-c", "from math import pi; print(f'{pi:.20f}')"]
          restartPolicy: Never
```

**Cron-синтаксис:** `минуты часы день_месяца месяц день_недели`
```
*/2 * * * *   → каждые 2 минуты
0 * * * *     → каждый час (в 00 минут)
0 3 * * *     → каждый день в 03:00
0 3 * * 1     → каждый понедельник в 03:00
```

### Service — ClusterIP

```yaml
# demo-clusterIP/03_service.yaml
apiVersion: v1
kind: Service
metadata:
  name: goweb-svc
  namespace: demo-clusterip
spec:
  selector:
    app: goapp          # выбирает Pods с этим label
  ports:
    - protocol: TCP
      port: 80          # порт Service
      targetPort: 8070  # порт Pod
```

### Service — NodePort

```yaml
# nodeport/np.yaml
apiVersion: v1
kind: Service
metadata:
  name: goweb-svc-nodeport
  namespace: demo-nodeport
spec:
  type: NodePort
  selector:
    app: goapp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8070
      nodePort: 30050   # порт на ноде (30000-32767)
```

### Service — LoadBalancer

```yaml
# loadbalancer/lb.yaml
apiVersion: v1
kind: Service
metadata:
  name: goweb-svc-lb
  namespace: demo-lb
spec:
  type: LoadBalancer
  selector:
    app: goapp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

### Ingress

```yaml
# ingress/04_ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  namespace: demo-ingress
spec:
  ingressClassName: nginx
  rules:
    - host: hello-world.example
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

### ConfigMap + envFrom

```yaml
# configmap/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: env-config
  namespace: demo-cm
data:
  LOG_LEVEL: "debug"
  DATABASE_URL: "postgres://db:5432/mydb"
```

```yaml
# configmap/deployment-env.yaml — подключение через envFrom
spec:
  containers:
    - name: app
      image: busybox
      envFrom:
        - configMapRef:
            name: env-config   # все ключи становятся переменными окружения
```

### ConfigMap — монтирование как файл

```yaml
# configmap/configmap.yaml — файл конфига
data:
  config.yaml: |
    log_level: "debug"
    database_address: 156.23.15.2
```

```yaml
# deployment — монтирование
spec:
  volumes:
    - name: config-vol
      configMap:
        name: first-cm
  containers:
    - volumeMounts:
        - name: config-vol
          mountPath: /etc/config   # файл появится как /etc/config/config.yaml
```

### Utility Pod для отладки сети (curl)

```yaml
# curl.yaml
apiVersion: v1
kind: Pod
metadata:
  name: curl
spec:
  containers:
    - name: curl
      image: radial/busyboxplus:curl
      command: ["sleep", "3600"]
```

```bash
kubectl apply -f curl.yaml
kubectl exec -it curl -- curl http://goweb-svc.demo-clusterip.svc.cluster.local/health
```

---

## Быстрые kubectl-команды

```bash
# Применить все манифесты в текущей папке
kubectl apply -f .

# Наблюдать за Pods в реальном времени
watch kubectl -n <ns> get pods

# Получить YAML развёрнутого ресурса
kubectl -n <ns> get deployment myapp -o yaml

# Сгенерировать YAML без применения
kubectl create deployment myapp --image=nginx --dry-run=client -o yaml

# Посмотреть, что изменится
kubectl diff -f deployment.yaml

# Откатить Deployment
kubectl -n <ns> rollout undo deployment/myapp

# История обновлений
kubectl -n <ns> rollout history deployment/myapp

# Масштабировать
kubectl -n <ns> scale deployment/myapp --replicas=5

# Запустить Job вручную из CronJob
kubectl -n <ns> create job manual-run --from=cronjob/pi-cronjob

# Проверить DNS внутри кластера
kubectl run dns-test -it --image=busybox --restart=Never --rm -- nslookup goweb-svc.demo-clusterip

# Отладка через временный Pod
kubectl run debug -it --image=ubuntu --restart=Never --rm -- bash
```

---

## 100 вопросов к собеседованию с ответами

**1. Чем Pod отличается от контейнера?**
Pod — минимальная единица в K8s. Содержит один или несколько контейнеров с общей сетью и storage. Контейнер — процесс в изолированном окружении.

**2. Почему не запускают Pod напрямую в production?**
Pod нет авто-восстановления после падения ноды. Используют Deployment/ReplicaSet, которые гарантируют нужное количество реплик.

**3. Что произойдёт, если удалить Pod из ReplicaSet?**
ReplicaSet немедленно создаст новый Pod, чтобы поддержать нужное количество реплик (reconcile loop).

**4. Чем Deployment лучше ReplicaSet?**
Deployment добавляет: rolling update без downtime, rollback на предыдущую версию, хранение истории RS.

**5. Как работает RollingUpdate?**
Постепенно создаёт новые Pods (не превышая `replicas + maxSurge`) и убивает старые (не опускаясь ниже `replicas - maxUnavailable`). Нет downtime.

**6. Как откатить Deployment на предыдущую версию?**
`kubectl rollout undo deployment/myapp` или на конкретную ревизию: `kubectl rollout undo deployment/myapp --to-revision=2`

**7. Когда использовать DaemonSet?**
Когда нужен один Pod на каждой ноде: сбор логов (Fluentd), мониторинг (node-exporter), сетевые агенты (Calico).

**8. Чем Job отличается от CronJob?**
Job — одноразовый запуск задачи. CronJob — периодический запуск Job по расписанию (cron-синтаксис).

**9. Что такое `backoffLimit` в Job?**
Сколько раз K8s попытается перезапустить Pod при ошибке перед тем, как Job помечается как Failed.

**10. Когда Service имеет тип ClusterIP?**
Когда нужен доступ только внутри кластера — между микросервисами. Получает виртуальный IP, доступный только изнутри кластера.

**11. Как Service находит нужные Pods?**
Через `selector` по labels. Service направляет трафик на Pods, у которых есть совпадающие labels. Endpoint Controller следит за актуальным списком IP-адресов.

**12. Зачем нужен Ingress если есть LoadBalancer?**
LoadBalancer = один внешний IP на один Service. Ingress = один IP, маршрутизация по hostname и path на множество Services. Дешевле и гибче.

**13. Что такое IngressController?**
Pod, который читает Ingress-объекты и настраивает реальный балансировщик (nginx, traefik). Без него Ingress-объекты ничего не делают.

**14. Чем ConfigMap отличается от Secret?**
ConfigMap — несекретная конфигурация в открытом виде. Secret — чувствительные данные в base64 (это не шифрование, но поддерживает encryption at rest в etcd).

**15. Три способа подключить ConfigMap к Pod?**
1. `envFrom.configMapRef` — все ключи как переменные окружения
2. `env.valueFrom.configMapKeyRef` — один ключ как переменная
3. `volumes.configMap` + `volumeMounts` — ключи как файлы

**16. Чем liveness отличается от readiness probe?**
liveness — если провалена, контейнер перезапускается. readiness — если провалена, Pod убирается из Endpoints Service (не получает трафик), но не перезапускается.

**17. Зачем нужна startupProbe?**
Защищает медленно стартующий контейнер от срабатывания livenessProbe раньше времени. Пока startupProbe не пройдёт — liveness не проверяется.

**18. Что такое requests и limits?**
requests — сколько CPU/памяти scheduler резервирует на ноде при размещении Pod. limits — максимум, который Pod может использовать (CPU throttling, OOMKill для памяти).

**19. Pod находится в статусе CrashLoopBackOff. Ваши действия?**
1. `kubectl logs <pod> --previous` — логи упавшего контейнера
2. `kubectl describe pod <pod>` — посмотреть events
3. Найти причину: ошибка в приложении, неверный CMD, нет нужной переменной/файла, не достигает БД

**20. Как Pod в одном Namespace обратится к Service в другом?**
Через полное DNS-имя: `<service-name>.<namespace>.svc.cluster.local`. Пример: `goweb-svc.demo-clusterip.svc.cluster.local:80`.

---

### Архитектура и компоненты

**21. Что хранится в etcd?**
Всё состояние кластера: объекты (Pods, Deployments, Services, Secrets...), конфигурация, статусы. etcd — единственный stateful компонент Control Plane. Без etcd кластер не работает.

**22. Что произойдёт если упадёт kube-apiserver?**
Нельзя будет делать изменения через kubectl, но уже запущенные Pods продолжат работать. kubelet на нодах продолжает поддерживать состояние локально.

**23. Что делает kube-scheduler?**
Выбирает Node для размещения Pod, который ещё не назначен на ноду. Учитывает: доступные ресурсы, taints/tolerations, affinity/anti-affinity, topologySpreadConstraints.

**24. Что такое controller-manager?**
Запускает встроенные контроллеры в одном процессе: Node Controller, Deployment Controller, ReplicaSet Controller, Job Controller и др. Каждый контроллер в цикле reconcile сверяет desired vs actual state.

**25. Чем kubeadm отличается от k3s?**
kubeadm — официальный инструмент для установки "полноценного" K8s (отдельные процессы компонентов, podman/docker обязателен). k3s — лёгкий дистрибутив Rancher, всё в одном бинарнике, встроенный containerd, подходит для edge/IoT/dev.

**26. Что такое kubelet и где он запускается?**
Агент Kubernetes на каждой Worker Node. Получает PodSpec от apiserver, запускает контейнеры через Container Runtime Interface (CRI). Следит за здоровьем контейнеров и отчитывается apiserver.

**27. Что такое kube-proxy и зачем он нужен?**
Работает на каждой ноде. Реализует сетевые правила (iptables или ipvs) для Services — перенаправляет трафик на нужный Pod. Именно kube-proxy обеспечивает работу ClusterIP/NodePort.

**28. Что такое Container Runtime Interface (CRI)?**
Стандартный API между kubelet и container runtime. Позволяет использовать разные runtime: containerd, CRI-O. Docker больше не поддерживается напрямую в K8s (убран dockershim в 1.24).

**29. Что такое CNI (Container Network Interface)?**
Плагины для настройки сети между Pods. Примеры: Calico, Flannel, Cilium, Weave. CNI обеспечивает: каждый Pod получает уникальный IP, Pods могут общаться между нодами.

**30. Как K8s гарантирует уникальность имён ресурсов?**
Имена уникальны в пределах namespace и kind. Pod "nginx" может существовать в namespace "dev" и "prod" одновременно, но не два Pod с именем "nginx" в одном namespace.

---

### Labels, Annotations, Selectors

**31. Чем Labels отличаются от Annotations?**
Labels — ключ/значение для идентификации и выборки объектов (используются в selector). Annotations — метаданные для инструментов и людей (не используются в selector), могут быть длинными строками.

**32. Для чего используются Labels на практике?**
Выборка Pods в Service/Deployment selector; фильтрация `kubectl get pods -l app=myapp`; топология для PodAffinity; nodeSelector. Labels — главный механизм связи объектов в K8s.

**33. Что такое selector в Service и как он работает?**
`selector` в Service описывает labels, которые должны быть у целевых Pods. Endpoint Controller постоянно обновляет список IP-адресов Pods, соответствующих selector. При несовпадении selector — трафик некуда слать.

---

### Volumes и хранилище

**34. Что такое PersistentVolume (PV)?**
Абстракция хранилища в K8s, подготовленная администратором или динамически через StorageClass. Существует независимо от Pod — данные сохраняются после удаления Pod.

**35. Что такое PersistentVolumeClaim (PVC)?**
Запрос от Pod на хранилище (размер, режим доступа). K8s находит подходящий PV и связывает их (binding). Pod монтирует PVC, не зная о конкретном PV.

**36. Какие режимы доступа у PV?**
- `ReadWriteOnce (RWO)` — один узел читает и пишет
- `ReadOnlyMany (ROX)` — много узлов только читают
- `ReadWriteMany (RWX)` — много узлов читают и пишут (поддерживают NFS, CephFS)

**37. Что такое StorageClass?**
Описывает "класс" хранилища и provisioner (aws-ebs, gce-pd, nfs, local-path). Позволяет динамически создавать PV при появлении PVC. Можно задать `storageClassName: standard` в PVC.

**38. Чем emptyDir отличается от hostPath?**
`emptyDir` — временный том, создаётся при старте Pod, удаляется вместе с ним. Используется для обмена данными между контейнерами в Pod. `hostPath` — монтирует путь с Node, данные остаются после удаления Pod (не переносимо между нодами).

---

### StatefulSet

**39. Чем StatefulSet отличается от Deployment?**
StatefulSet гарантирует: стабильные имена Pods (pod-0, pod-1...), стабильные DNS-имена, упорядоченный запуск/остановку, отдельный PVC для каждой реплики. Нужен для БД (Postgres, Kafka, Cassandra).

**40. Почему нельзя использовать Deployment для базы данных?**
Deployment не гарантирует порядок запуска и стабильные идентификаторы. Реплики БД должны знать, кто master, кто replica — это требует предсказуемых имён и отдельного хранилища на каждом Pod.

---

### RBAC (Role-Based Access Control)

**41. Что такое RBAC в Kubernetes?**
Механизм контроля доступа. Определяет кто (ServiceAccount/User/Group) может делать что (verbs: get/list/create/delete) с какими ресурсами (pods/deployments/secrets).

**42. Чем Role отличается от ClusterRole?**
`Role` — права в пределах одного Namespace. `ClusterRole` — права на уровне всего кластера (или для non-namespaced ресурсов: nodes, PV).

**43. Чем RoleBinding отличается от ClusterRoleBinding?**
`RoleBinding` привязывает Role/ClusterRole к субъекту в конкретном Namespace. `ClusterRoleBinding` — привязывает ClusterRole к субъекту на уровне всего кластера.

**44. Что такое ServiceAccount?**
Идентификатор для процессов внутри Pod. Каждый Pod автоматически получает токен default ServiceAccount своего namespace. Используется для обращений к K8s API из кода.

**45. Как Pod использует ServiceAccount?**
kubelet монтирует токен ServiceAccount в `/var/run/secrets/kubernetes.io/serviceaccount/token`. Приложение использует этот токен для аутентификации в kube-apiserver.

---

### Scheduling: Taints, Tolerations, Affinity

**46. Что такое Taints и Tolerations?**
`Taint` — "метка отпугивания" на Node: не ставить сюда Pods без специального разрешения. `Toleration` — разрешение в Pod: "я могу работать на ноде с таким taint". Используется для выделения нод под специфические задачи (GPU, spot-instances).

**47. Три эффекта Taint?**
- `NoSchedule` — не планировать новые Pods без toleration
- `PreferNoSchedule` — избегать, но не запрещать
- `NoExecute` — выселить уже запущенные Pods без toleration

**48. Что такое Node Affinity?**
Правила притяжения/отталкивания Pod к конкретным нодам по labels. `requiredDuringScheduling` — жёсткое требование (Pod не запустится если нет подходящей ноды). `preferredDuringScheduling` — мягкое предпочтение.

**49. Что такое Pod Anti-Affinity и зачем?**
Правило: не ставить этот Pod рядом с другими Pods, соответствующими selector. Используется для распределения реплик по разным нодам/зонам для высокой доступности.

**50. Что такое nodeSelector?**
Простейший способ ограничить, на каких нодах запускается Pod: указать labels ноды. Менее гибкий чем Node Affinity, но проще в конфигурации.
```yaml
nodeSelector:
  disktype: ssd
```

---

### Autoscaling и ресурсы

**51. Что такое HPA (Horizontal Pod Autoscaler)?**
Автоматически масштабирует количество реплик Deployment/ReplicaSet/StatefulSet на основе метрик (CPU, memory, custom). Требует metrics-server.

**52. Что такое VPA (Vertical Pod Autoscaler)?**
Автоматически подбирает requests/limits для контейнеров на основе исторического потребления. Может перезапускать Pods для применения новых значений.

**53. Что такое Resource Quota?**
Объект Namespace-уровня: ограничивает суммарные ресурсы в namespace (количество Pods, CPU, memory, PVC). Если лимит исчерпан — новые объекты не создаются.

**54. Что такое LimitRange?**
Задаёт дефолтные и максимальные requests/limits для Pods и контейнеров в Namespace. Если Pod не указал limits — LimitRange применяет дефолтные значения автоматически.

**55. Что такое QoS классы и как они влияют на вытеснение?**
При нехватке памяти на ноде K8s вытесняет Pods в порядке: `BestEffort` (нет requests/limits) → `Burstable` (requests < limits) → `Guaranteed` (requests == limits вытесняются последними).

---

### Сеть: углублённо

**56. Как работает kube-dns / CoreDNS?**
CoreDNS — системный Pod в namespace `kube-system`. Обрабатывает DNS-запросы внутри кластера. Создаёт A-записи для Services и Pods. Настраивается через ConfigMap `coredns`.

**57. Что такое NetworkPolicy?**
Объект K8s для ограничения сетевого трафика между Pods (L3/L4 firewall). По умолчанию все Pods могут общаться со всеми. NetworkPolicy применяется только если CNI поддерживает (Calico, Cilium, Weave).

**58. Что такое headless Service?**
Service с `clusterIP: None`. Не создаёт виртуальный IP — DNS возвращает напрямую IP-адреса Pods. Используется в StatefulSet для прямого обращения к конкретным репликам (pod-0.service.ns).

**59. Как Pod получает IP-адрес?**
При создании Pod kubelet вызывает CNI-плагин, который выделяет IP из Pod CIDR и настраивает сетевой интерфейс. IP Pod уникален в кластере, но не персистентен — меняется при пересоздании.

**60. Что такое Endpoints объект?**
Автоматически создаваемый K8s объект, связанный с Service. Содержит актуальный список IP:порт всех Pods, соответствующих selector Service. `kubectl get endpoints -n <ns>`.

---

### ConfigMap и Secret: углублённо

**61. Обновляется ли ConfigMap в Pod автоматически?**
При монтировании как volume — да, обновляется автоматически (с задержкой ~1 минута). При использовании как env-переменные — нет, нужен перезапуск Pod.

**62. Безопасен ли Secret?**
Base64 — это не шифрование, а кодирование. По умолчанию Secrets хранятся в etcd в открытом виде. Для безопасности нужно включить encryption at rest в etcd и настроить RBAC, ограничив доступ к Secrets.

**63. Что такое imagePullSecret?**
Secret типа `kubernetes.io/dockerconfigjson` с credentials для Docker Registry. Указывается в Pod spec как `imagePullSecrets`, чтобы kubelet мог скачать приватный образ.
```yaml
spec:
  imagePullSecrets:
    - name: registry-secret
```

---

### Debug и troubleshooting: углублённо

**64. Как посмотреть события в кластере?**
```bash
kubectl get events -n <ns> --sort-by='.lastTimestamp'
kubectl get events --field-selector reason=OOMKilling
```

**65. Pod завис в статусе Terminating. Что делать?**
Pod не может завершиться (finalizer не снят или процесс не реагирует на SIGTERM). Принудительное удаление:
```bash
kubectl delete pod <pod> --force --grace-period=0
```

**66. Что такое `kubectl port-forward` и когда использовать?**
Проброс порта с локальной машины на Pod/Service в кластере. Используется для отладки без создания Service.
```bash
kubectl -n <ns> port-forward pod/mypod 8080:8080
kubectl -n <ns> port-forward svc/myservice 8080:80
```

**67. Как посмотреть ресурсы, потребляемые Pods?**
```bash
kubectl top pods -n <ns>
kubectl top nodes
```
Требует установленного metrics-server.

**68. Что такое `kubectl exec` и чем отличается от `kubectl debug`?**
`exec` — выполняет команду в уже запущенном контейнере (контейнер должен содержать нужный бинарник). `debug` — добавляет новый ephemeral container с нужным образом (busybox, ubuntu) без перезапуска Pod.

**69. Pod не может подключиться к другому сервису. Алгоритм проверки?**
1. `kubectl exec` в Pod, `curl <service-dns>:<port>` — проверить DNS и подключение
2. `kubectl get endpoints <svc>` — есть ли реальные IP за Service
3. `kubectl describe svc` — проверить selector
4. `kubectl get pods -l <selector>` — есть ли Pods с нужными labels
5. Проверить NetworkPolicy если настроен

**70. Как проверить, какой Node используется для Pod?**
```bash
kubectl get pod <pod> -o wide           # колонка NODE
kubectl describe pod <pod> | grep Node
```

---

### Helm

**71. Что такое Helm и зачем он нужен?**
Пакетный менеджер для Kubernetes. Chart — набор шаблонов YAML + values. Helm генерирует манифесты из шаблонов, применяет их и отслеживает версии (release). Упрощает установку и обновление сложных приложений.

**72. Что такое Helm Release?**
Конкретная установка Chart в кластер с определёнными values. Один Chart можно установить несколько раз под разными именами (разные Release).

**73. Чем `helm install` отличается от `helm upgrade --install`?**
`install` — только первый раз, упадёт если release уже есть. `upgrade --install` — idempotent: установит если нет, обновит если есть. В CI/CD всегда используют `upgrade --install`.

**74. Как посмотреть текущие values установленного release?**
```bash
helm get values <release> -n <ns>
helm get values <release> -n <ns> --all   # включая дефолтные
```

**75. Что такое Helm hooks?**
Хуки позволяют запускать Jobs в определённые моменты жизненного цикла release: `pre-install`, `post-install`, `pre-upgrade`, `pre-delete`. Пример: миграции БД перед деплоем.

---

### Docker: углублённо

**76. Что такое Docker layer cache и как его использовать?**
Каждая инструкция в Dockerfile создаёт слой. Слои кешируются: если инструкция и контекст не изменились — берётся кеш. Для эффективности: сначала COPY файлов зависимостей (package.json, go.mod), потом RUN install, потом COPY исходников.

**77. Что такое .dockerignore?**
Файл, аналогичный .gitignore. Исключает файлы и папки из Docker build context. Уменьшает размер контекста и время сборки. Обязательно добавлять: `node_modules/`, `.git/`, `*.log`.

**78. Как посмотреть слои образа?**
```bash
docker history <image>
docker inspect <image>
```

**79. Что такое docker-compose и чем отличается от K8s?**
docker-compose — инструмент для локального запуска multi-container приложений на одном хосте. K8s — production-grade оркестратор для распределённых систем на множестве нод. docker-compose не масштабируется горизонтально.

**80. Что такое ENTRYPOINT в виде shell формы vs exec формы?**
Shell форма: `ENTRYPOINT ./app` — запускает через `/bin/sh -c`, PID 1 = sh, процесс не получает сигналы (SIGTERM). Exec форма: `ENTRYPOINT ["./app"]` — приложение = PID 1, получает сигналы напрямую. Всегда используйте exec форму.

---

### Жизненный цикл Pod и контейнера

**81. Что такое Init Containers?**
Контейнеры, которые запускаются и завершаются до старта основных контейнеров Pod. Выполняются последовательно. Используются для: ожидания БД, миграций, подготовки конфигов.
```yaml
initContainers:
  - name: wait-db
    image: busybox
    command: ['sh', '-c', 'until nc -z db-service 5432; do sleep 2; done']
```

**82. Что такое sidecar контейнер?**
Дополнительный контейнер в одном Pod с основным приложением. Разделяет сеть и volumes. Примеры: log-shipper (Fluentd), service mesh proxy (Envoy/Istio), secrets injector (Vault Agent).

**83. Что происходит при удалении Pod? Что такое terminationGracePeriodSeconds?**
K8s отправляет SIGTERM процессу. Если через `terminationGracePeriodSeconds` (default: 30 сек) процесс не завершился — отправляет SIGKILL. Приложение должно обрабатывать SIGTERM для graceful shutdown.

**84. Что такое lifecycle hooks в контейнере?**
`postStart` — выполняется сразу после запуска контейнера (асинхронно). `preStop` — выполняется перед остановкой (до SIGTERM). Пример preStop для graceful drain:
```yaml
lifecycle:
  preStop:
    exec:
      command: ["/bin/sh", "-c", "sleep 5"]
```

**85. Что такое PodDisruptionBudget (PDB)?**
Гарантирует минимальное количество доступных Pods во время voluntary disruptions (drain ноды, обновление кластера). Защищает от одновременного вывода слишком многих реплик.
```yaml
spec:
  minAvailable: 2      # или maxUnavailable: 1
  selector:
    matchLabels:
      app: myapp
```

---

### Безопасность

**86. Что такое SecurityContext?**
Настройки безопасности для Pod или контейнера: от какого пользователя запускать, read-only filesystem, capabilities, seccomp профиль.
```yaml
securityContext:
  runAsUser: 1000
  runAsNonRoot: true
  readOnlyRootFilesystem: true
```

**87. Что такое Pod Security Standards (PSS)?**
Встроенные политики безопасности K8s (заменили PodSecurityPolicy): `privileged` (без ограничений), `baseline` (базовая защита), `restricted` (строгие ограничения). Применяются на уровне Namespace через labels.

**88. Зачем запускать контейнер не от root?**
Если злоумышленник выйдет за пределы контейнера — у него не будет прав root на ноде. Минимизация ущерба. Best practice: `runAsNonRoot: true`, создавать пользователя в Dockerfile.

---

### Мониторинг и observability

**89. Что такое metrics-server?**
Лёгкий агрегатор метрик ресурсов (CPU/memory) с нод и Pods. Нужен для HPA и команды `kubectl top`. Не для долгосрочного хранения — для этого Prometheus.

**90. Что такое liveness probe с точки зрения операций?**
Способ сказать K8s "если приложение не отвечает на этот endpoint — оно зависло, перезапусти его". Должен проверять только базовую работоспособность (не зависимости), иначе каскадные перезапуски.

---

### CI/CD и GitOps

**91. Что такое GitOps?**
Подход, где Git — единственный источник истины для инфраструктуры. Изменения вносятся через PR/MR. Оператор (ArgoCD, Flux) автоматически синхронизирует кластер с состоянием в Git.

**92. Чем ArgoCD отличается от GitLab CI для деплоя в K8s?**
GitLab CI — push-модель: pipeline явно вызывает `helm upgrade`. ArgoCD — pull-модель: оператор в кластере следит за Git и сам применяет изменения. ArgoCD лучше для multi-cluster и audit.

**93. Что такое Kaniko и зачем он нужен?**
Инструмент для сборки Docker-образов внутри K8s Pod без Docker daemon (не нужен privileged mode). Безопаснее docker-in-docker (DinD). Популярен в GitLab CI/CD внутри K8s.

---

### Разное и практика

**94. Как быстро создать YAML манифест не пиша с нуля?**
```bash
# Deployment
kubectl create deployment myapp --image=nginx --dry-run=client -o yaml > deployment.yaml
# Service
kubectl expose deployment myapp --port=80 --dry-run=client -o yaml > service.yaml
# ConfigMap
kubectl create configmap myconfig --from-literal=key=value --dry-run=client -o yaml
```

**95. Что такое `kubectl explain`?**
Встроенная документация по полям объектов K8s прямо в терминале:
```bash
kubectl explain pod.spec.containers
kubectl explain deployment.spec.strategy.rollingUpdate
```

**96. Что такое Namespace `kube-system`?**
Зарезервированный namespace для системных компонентов K8s: CoreDNS, kube-proxy, metrics-server, ingress-controller. Не стоит деплоить туда пользовательские приложения.

**97. Как переключаться между кластерами и namespace?**
```bash
kubectl config get-contexts                    # список контекстов
kubectl config use-context <context-name>      # переключить кластер
kubectl config set-context --current --namespace=my-ns  # дефолтный ns
# Или использовать утилиты: kubectx (кластеры), kubens (namespace)
```

**98. Что такое `finalizer` в K8s?**
Строка в `metadata.finalizers` объекта. K8s не удаляет объект пока список finalizers не пуст. Контроллеры добавляют finalizer, выполняют cleanup, потом снимают finalizer. Причина "зависших" объектов при удалении.

**99. Чем `kubectl apply` отличается от `kubectl create` и `kubectl replace`?**
`create` — создаёт объект, ошибка если уже есть. `replace` — полностью заменяет (объект должен существовать). `apply` — idempotent: создаёт если нет, обновляет если есть, хранит аннотацию с предыдущим состоянием (three-way merge).

**100. Что нужно сделать перед удалением ноды из кластера?**
1. `kubectl cordon <node>` — запретить планировать новые Pods на ноду
2. `kubectl drain <node> --ignore-daemonsets --delete-emptydir-data` — выселить все Pods
3. После обслуживания: `kubectl uncordon <node>` — вернуть ноду в ротацию

# Вопросы и ответы к собеседованию — Docker + Kubernetes

> Источник: INTERVIEW_THEORY.md + CHEATSHEET.md
> Формат: Q&A, сгруппировано по темам.

---

## Docker

**Q: Чем отличается CMD от ENTRYPOINT?**
ENTRYPOINT — неизменяемая точка входа (бинарник), CMD — аргументы по умолчанию к ней. CMD легко переопределить в `docker run` (полностью заменяется). ENTRYPOINT — только через `--entrypoint`. Shell-форма (`ENTRYPOINT ./app`) плохая — PID 1 = sh, SIGTERM не доходит. Exec-форма (`ENTRYPOINT ["./app"]`) — правильная.

**Q: Что такое multi-stage build и зачем?**
Позволяет использовать большой образ для сборки (golang, node) и маленький для финального образа. Итоговый образ не содержит компилятор/исходники — меньше размер, меньше уязвимостей.

**Q: Чем volume отличается от bind mount?**
Volume управляется Docker, хранится в `/var/lib/docker/volumes`, переносимый. Bind mount — монтирование конкретного пути хоста, зависит от структуры FS. Volume лучше для production, bind mount удобен при разработке.

**Q: Что такое Docker layer cache и как его использовать?**
Каждая инструкция в Dockerfile создаёт слой. Если инструкция и контекст не изменились — берётся кеш. Для эффективности: сначала `COPY package*.json .`, потом `RUN npm ci`, потом `COPY . .` — зависимости кешируются отдельно от исходников.

**Q: Что такое .dockerignore?**
Файл, аналогичный .gitignore. Исключает файлы из Docker build context. Уменьшает размер контекста и время сборки. Обязательно: `node_modules/`, `.git/`, `*.log`, тесты.

**Q: Как посмотреть слои образа?**
`docker history <image>` и `docker inspect <image>`.

**Q: Что такое docker-compose и чем отличается от K8s?**
docker-compose — инструмент для локального запуска multi-container приложений на одном хосте. K8s — production-grade оркестратор для распределённых систем на множестве нод. docker-compose не масштабируется горизонтально.

---

## Kubernetes — Workload объекты

**Q: Чем Pod отличается от контейнера?**
Pod — минимальная единица в K8s. Содержит один или несколько контейнеров с общей сетью и storage. Контейнер — процесс в изолированном окружении.

**Q: Почему не запускают Pod напрямую в production?**
Pod нет авто-восстановления после падения ноды. Используют Deployment/ReplicaSet, которые гарантируют нужное количество реплик через reconcile loop.

**Q: Что произойдёт, если удалить Pod из ReplicaSet?**
ReplicaSet немедленно создаст новый Pod, чтобы поддержать нужное количество реплик (reconcile loop: desired state = 3, actual = 2 → создать Pod).

**Q: Чем Deployment лучше ReplicaSet?**
Deployment добавляет: rolling update без downtime, rollback на предыдущую версию (`rollout undo`), хранение истории RS.

**Q: Как откатить Deployment на предыдущую версию?**
`kubectl rollout undo deployment/myapp` или на конкретную ревизию: `kubectl rollout undo deployment/myapp --to-revision=2`

**Q: Как работает RollingUpdate?**
Постепенно создаёт новые Pods (не превышая `replicas + maxSurge`) и убивает старые (не опускаясь ниже `replicas - maxUnavailable`). Нет downtime.

**Q: Когда использовать DaemonSet?**
Когда нужен один Pod на каждой ноде: сбор логов (Fluentd), мониторинг (node-exporter), сетевые агенты (Calico).

**Q: Чем StatefulSet отличается от Deployment?**
StatefulSet гарантирует: стабильные имена Pods (pod-0, pod-1...), стабильные DNS-имена, упорядоченный запуск/остановку, отдельный PVC для каждой реплики. Нужен для БД (Postgres, Kafka, Cassandra).

**Q: Почему нельзя использовать Deployment для базы данных?**
Deployment не гарантирует порядок запуска и стабильные идентификаторы. Реплики БД должны знать, кто master, кто replica — это требует предсказуемых имён и отдельного хранилища на каждом Pod.

**Q: Чем Job отличается от CronJob?**
Job — одноразовый запуск задачи до успешного завершения. CronJob — периодический запуск Job по расписанию (cron-синтаксис).

**Q: Что такое `backoffLimit` в Job?**
Сколько раз K8s попытается перезапустить Pod при ошибке перед тем, как Job помечается как Failed.

---

## Kubernetes — Архитектура и компоненты

**Q: Из чего состоит Control Plane?**
API Server (единая точка входа), etcd (хранилище состояния), Scheduler (выбирает Node для Pod), Controller Manager (reconcile loop).

**Q: Что хранится в etcd?**
Всё состояние кластера: объекты (Pods, Deployments, Services, Secrets...), конфигурация, статусы. etcd — единственный stateful компонент Control Plane. Без etcd кластер не работает.

**Q: Что произойдёт если упадёт kube-apiserver?**
Нельзя будет делать изменения через kubectl, но уже запущенные Pods продолжат работать. kubelet на нодах продолжает поддерживать состояние локально.

**Q: Что делает kube-scheduler?**
Выбирает Node для размещения Pod, который ещё не назначен на ноду. Учитывает: доступные ресурсы, taints/tolerations, affinity/anti-affinity, topologySpreadConstraints.

**Q: Что такое controller-manager?**
Запускает встроенные контроллеры в одном процессе: Node Controller, Deployment Controller, ReplicaSet Controller, Job Controller и др. Каждый контроллер в цикле reconcile сверяет desired vs actual state.

**Q: Что такое kubelet и где он запускается?**
Агент Kubernetes на каждой Worker Node. Получает PodSpec от apiserver, запускает контейнеры через Container Runtime Interface (CRI). Следит за здоровьем контейнеров и отчитывается apiserver.

**Q: Что такое kube-proxy и зачем он нужен?**
Работает на каждой ноде. Реализует сетевые правила (iptables или ipvs) для Services — перенаправляет трафик на нужный Pod. Именно kube-proxy обеспечивает работу ClusterIP/NodePort.

**Q: Что такое Container Runtime Interface (CRI)?**
Стандартный API между kubelet и container runtime. Позволяет использовать разные runtime: containerd, CRI-O. Docker больше не поддерживается напрямую в K8s (убран dockershim в 1.24).

**Q: Что такое CNI (Container Network Interface)?**
Плагины для настройки сети между Pods. Примеры: Calico, Flannel, Cilium, Weave. CNI обеспечивает: каждый Pod получает уникальный IP, Pods могут общаться между нодами.

**Q: Чем kubeadm отличается от k3s?**
kubeadm — официальный инструмент для установки "полноценного" K8s. k3s — лёгкий дистрибутив Rancher, всё в одном бинарнике, встроенный containerd, подходит для edge/IoT/dev.

---

## Kubernetes — Сеть

**Q: Когда Service имеет тип ClusterIP?**
Когда нужен доступ только внутри кластера — между микросервисами. Получает виртуальный IP, доступный только изнутри кластера.

**Q: Как Service находит нужные Pods?**
Через `selector` по labels. Endpoint Controller следит за актуальным списком IP-адресов Pods с совпадающими labels. При несовпадении selector — трафик некуда слать.

**Q: Зачем нужен Ingress если есть LoadBalancer?**
LoadBalancer = один внешний IP на один Service. Ingress = один IP, маршрутизация по hostname и path на множество Services. Дешевле и гибче.

**Q: Что такое IngressController?**
Pod, который читает Ingress-объекты и настраивает реальный балансировщик (nginx, traefik). Без него Ingress-объекты ничего не делают.

**Q: Как Pod'ы находят друг друга?**
Через Service DNS: `<service>.<namespace>.svc.cluster.local`. CoreDNS создаёт записи для каждого Service автоматически.

**Q: Как Pod в одном Namespace обратится к Service в другом?**
Через полное DNS-имя: `<service-name>.<namespace>.svc.cluster.local`. Пример: `goweb-svc.demo.svc.cluster.local:80`.

**Q: Как работает kube-dns / CoreDNS?**
CoreDNS — системный Pod в namespace `kube-system`. Обрабатывает DNS-запросы внутри кластера. Создаёт A-записи для Services и Pods. Настраивается через ConfigMap `coredns`.

**Q: Что такое NetworkPolicy?**
Объект K8s для ограничения сетевого трафика между Pods (L3/L4 firewall). По умолчанию все Pods могут общаться со всеми. Работает только если CNI поддерживает (Calico, Cilium, Weave).

**Q: Что такое headless Service?**
Service с `clusterIP: None`. Не создаёт виртуальный IP — DNS возвращает напрямую IP-адреса Pods. Используется в StatefulSet для прямого обращения к конкретным репликам (pod-0.service.ns).

**Q: Что такое Endpoints объект?**
Автоматически создаваемый K8s объект, связанный с Service. Содержит актуальный список IP:порт всех Pods, соответствующих selector Service. `kubectl get endpoints -n <ns>`.

---

## Kubernetes — Labels, Volumes

**Q: Чем Labels отличаются от Annotations?**
Labels — ключ/значение для идентификации и выборки объектов (используются в selector). Annotations — метаданные для инструментов и людей (не используются в selector), могут быть длинными строками.

**Q: Для чего используются Labels на практике?**
Выборка Pods в Service/Deployment selector; фильтрация `kubectl get pods -l app=myapp`; топология для PodAffinity; nodeSelector. Labels — главный механизм связи объектов в K8s.

**Q: Что такое PersistentVolume (PV)?**
Абстракция хранилища в K8s, подготовленная администратором или динамически через StorageClass. Существует независимо от Pod — данные сохраняются после удаления Pod.

**Q: Что такое PersistentVolumeClaim (PVC)?**
Запрос от Pod на хранилище (размер, режим доступа). K8s находит подходящий PV и связывает их (binding). Pod монтирует PVC, не зная о конкретном PV.

**Q: Какие режимы доступа у PV?**
- `ReadWriteOnce (RWO)` — один узел читает и пишет
- `ReadOnlyMany (ROX)` — много узлов только читают
- `ReadWriteMany (RWX)` — много узлов читают и пишут (NFS, CephFS)

**Q: Чем emptyDir отличается от hostPath?**
`emptyDir` — временный том, создаётся при старте Pod, удаляется вместе с ним (обмен данными между контейнерами). `hostPath` — монтирует путь с Node, данные остаются после удаления Pod (не переносимо между нодами).

---

## Kubernetes — Lifecycle: Probes, Resources, ConfigMap, Secret

**Q: Разница между liveness и readiness probe?**
liveness — если провалена, контейнер перезапускается. readiness — если провалена, Pod убирается из балансировки Service (не получает трафик), но не перезапускается.

**Q: Зачем нужна startupProbe?**
Защищает медленно стартующий контейнер от срабатывания livenessProbe раньше времени. Пока startupProbe не пройдёт — liveness не проверяется.

**Q: Чем requests отличается от limits?**
requests — сколько CPU/памяти scheduler резервирует на ноде при размещении Pod. limits — максимум, который Pod может использовать (CPU throttling, OOMKill для памяти).

**Q: Что такое QoS классы и как они влияют на вытеснение?**
При нехватке памяти на ноде K8s вытесняет Pods в порядке: `BestEffort` (нет requests/limits) → `Burstable` (requests < limits) → `Guaranteed` (requests == limits — вытесняются последними).

**Q: Что такое HPA?**
Horizontal Pod Autoscaler — автоматически масштабирует количество реплик Deployment/ReplicaSet/StatefulSet на основе метрик (CPU, memory, custom). Требует metrics-server.

**Q: Что такое Resource Quota?**
Объект Namespace-уровня: ограничивает суммарные ресурсы в namespace (количество Pods, CPU, memory, PVC). Если лимит исчерпан — новые объекты не создаются.

**Q: Чем ConfigMap отличается от Secret?**
ConfigMap — несекретная конфигурация (plain text). Secret — чувствительные данные (base64, не шифрование!). Секреты можно настроить на encryption at rest в etcd.

**Q: Три способа подключить ConfigMap к Pod?**
1. `envFrom.configMapRef` — все ключи как переменные окружения
2. `env.valueFrom.configMapKeyRef` — один ключ как переменная
3. `volumes.configMap` + `volumeMounts` — ключи как файлы

**Q: Обновляется ли ConfigMap в Pod автоматически?**
При монтировании как volume — да, обновляется автоматически (с задержкой ~1 минута). При использовании как env-переменные — нет, нужен перезапуск Pod.

**Q: Безопасен ли Secret?**
Base64 — это не шифрование, а кодирование. По умолчанию Secrets хранятся в etcd в открытом виде. Для безопасности нужно включить encryption at rest в etcd и настроить RBAC.

**Q: Что такое imagePullSecret?**
Secret типа `kubernetes.io/dockerconfigjson` с credentials для Docker Registry. Указывается в Pod spec как `imagePullSecrets`, чтобы kubelet мог скачать приватный образ.

**Q: Как работает RollingUpdate?**
Создаёт новый ReplicaSet с новым образом, постепенно поднимает новые Pods и убивает старые. `maxUnavailable` — сколько Pods может быть недоступно, `maxSurge` — сколько дополнительных Pods можно создать.

**Q: Что такое Init Containers?**
Контейнеры, которые запускаются и завершаются до старта основных контейнеров Pod. Выполняются последовательно. Используются для: ожидания БД, миграций, подготовки конфигов.

**Q: Что такое sidecar контейнер?**
Дополнительный контейнер в одном Pod с основным приложением. Разделяет сеть и volumes. Примеры: log-shipper (Fluentd), service mesh proxy (Envoy/Istio), secrets injector (Vault Agent).

**Q: Что происходит при удалении Pod? Что такое terminationGracePeriodSeconds?**
K8s отправляет SIGTERM процессу. Если через `terminationGracePeriodSeconds` (default: 30 сек) процесс не завершился — отправляет SIGKILL. Приложение должно обрабатывать SIGTERM для graceful shutdown.

**Q: Что такое PodDisruptionBudget (PDB)?**
Гарантирует минимальное количество доступных Pods во время voluntary disruptions (drain ноды, обновление кластера). Защищает от одновременного вывода слишком многих реплик.

---

## Kubernetes — RBAC и безопасность

**Q: Что такое RBAC в Kubernetes?**
Механизм контроля доступа. Определяет кто (ServiceAccount/User/Group) может делать что (verbs: get/list/create/delete) с какими ресурсами (pods/deployments/secrets).

**Q: Чем Role отличается от ClusterRole?**
`Role` — права в пределах одного Namespace. `ClusterRole` — права на уровне всего кластера (или для non-namespaced ресурсов: nodes, PV).

**Q: Что такое ServiceAccount?**
Идентификатор для процессов внутри Pod. Каждый Pod автоматически получает токен default ServiceAccount своего namespace. Используется для обращений к K8s API из кода.

**Q: Что такое SecurityContext?**
Настройки безопасности для Pod или контейнера: от какого пользователя запускать, read-only filesystem, capabilities, seccomp профиль.

**Q: Зачем запускать контейнер не от root?**
Если злоумышленник выйдет за пределы контейнера — у него не будет прав root на ноде. Минимизация ущерба. Best practice: `runAsNonRoot: true`, создавать пользователя в Dockerfile.

**Q: Что такое Taints и Tolerations?**
`Taint` — "метка отпугивания" на Node: не ставить сюда Pods без специального разрешения. `Toleration` — разрешение в Pod: "я могу работать на ноде с таким taint". Используется для выделения нод под GPU, spot-instances.

**Q: Что такое Node Affinity?**
Правила притяжения/отталкивания Pod к конкретным нодам по labels. `requiredDuringScheduling` — жёсткое требование. `preferredDuringScheduling` — мягкое предпочтение.

---

## Kubernetes — Debug и Troubleshooting

**Q: Pod находится в статусе CrashLoopBackOff. Ваши действия?**
1. `kubectl logs <pod> --previous` — логи упавшего контейнера
2. `kubectl describe pod <pod>` — посмотреть events
3. Найти причину: ошибка в приложении, неверный CMD, нет нужной переменной/файла, не достигает БД

**Q: Как попасть внутрь контейнера для отладки?**
`kubectl exec -it <pod> -- bash` — если bash есть. `kubectl debug -it <pod> --image=ubuntu` — если bash нет (добавляет ephemeral container). `nsenter` — если нужен доступ на уровне ноды.

**Q: Что такое DaemonSet?**
Гарантирует запуск одного Pod на каждой Node. Используется для агентов мониторинга, сбора логов, сетевых плагинов.

**Q: Как посмотреть события в кластере?**
`kubectl get events -n <ns> --sort-by='.lastTimestamp'` и `kubectl get events --field-selector reason=OOMKilling`

**Q: Pod завис в статусе Terminating. Что делать?**
Pod не может завершиться (finalizer не снят или процесс не реагирует на SIGTERM). Принудительное удаление: `kubectl delete pod <pod> --force --grace-period=0`

**Q: Что такое `kubectl port-forward` и когда использовать?**
Проброс порта с локальной машины на Pod/Service в кластере. Используется для отладки без создания Service.

**Q: Pod не может подключиться к другому сервису. Алгоритм проверки?**
1. `kubectl exec` в Pod, `curl <service-dns>:<port>` — проверить DNS и подключение
2. `kubectl get endpoints <svc>` — есть ли реальные IP за Service
3. `kubectl describe svc` — проверить selector
4. `kubectl get pods -l <selector>` — есть ли Pods с нужными labels
5. Проверить NetworkPolicy если настроен

**Q: Как проверить на какой Node стоит Pod?**
`kubectl get pod <pod> -o wide` (колонка NODE) или `kubectl describe pod <pod> | grep Node`

**Q: Что нужно сделать перед удалением ноды из кластера?**
1. `kubectl cordon <node>` — запретить планировать новые Pods
2. `kubectl drain <node> --ignore-daemonsets --delete-emptydir-data` — выселить все Pods
3. После обслуживания: `kubectl uncordon <node>` — вернуть в ротацию

**Q: Что такое `finalizer` в K8s?**
Строка в `metadata.finalizers` объекта. K8s не удаляет объект пока список finalizers не пуст. Контроллеры добавляют finalizer, выполняют cleanup, потом снимают finalizer. Причина "зависших" объектов при удалении.

**Q: Как посмотреть ресурсы, потребляемые Pods?**
`kubectl top pods -n <ns>` и `kubectl top nodes`. Требует установленного metrics-server.

---

## Kubernetes — Инструменты

**Q: Чем `kubectl apply` отличается от `kubectl create` и `kubectl replace`?**
`create` — создаёт объект, ошибка если уже есть. `replace` — полностью заменяет (объект должен существовать). `apply` — idempotent: создаёт если нет, обновляет если есть, хранит аннотацию с предыдущим состоянием (three-way merge).

**Q: Как быстро создать YAML манифест не пиша с нуля?**
```bash
kubectl create deployment myapp --image=nginx --dry-run=client -o yaml > deployment.yaml
kubectl expose deployment myapp --port=80 --dry-run=client -o yaml > service.yaml
kubectl create configmap myconfig --from-literal=key=value --dry-run=client -o yaml
```

**Q: Что такое `kubectl explain`?**
Встроенная документация по полям объектов K8s прямо в терминале: `kubectl explain pod.spec.containers`

**Q: Что такое Namespace `kube-system`?**
Зарезервированный namespace для системных компонентов K8s: CoreDNS, kube-proxy, metrics-server, ingress-controller. Не стоит деплоить туда пользовательские приложения.

**Q: Как переключаться между кластерами и namespace?**
```bash
kubectl config get-contexts
kubectl config use-context <context-name>
kubectl config set-context --current --namespace=my-ns
# Или: kubectx (кластеры), kubens (namespace)
```

---

## Helm

**Q: Что такое Helm и зачем он нужен?**
Пакетный менеджер для Kubernetes. Chart — набор шаблонов YAML + values. Helm генерирует манифесты из шаблонов, применяет их и отслеживает версии (release). Упрощает установку и обновление сложных приложений.

**Q: Что такое Helm Release?**
Конкретная установка Chart в кластер с определёнными values. Один Chart можно установить несколько раз под разными именами.

**Q: Чем `helm install` отличается от `helm upgrade --install`?**
`install` — только первый раз, упадёт если release уже есть. `upgrade --install` — idempotent: установит если нет, обновит если есть. В CI/CD всегда используют `upgrade --install`.

**Q: Что такое Helm hooks?**
Хуки позволяют запускать Jobs в определённые моменты жизненного цикла release: `pre-install`, `post-install`, `pre-upgrade`, `pre-delete`. Пример: миграции БД перед деплоем.

---

## CI/CD и GitOps

**Q: Что такое GitOps?**
Подход, где Git — единственный источник истины для инфраструктуры. Изменения вносятся через PR/MR. Оператор (ArgoCD, Flux) автоматически синхронизирует кластер с состоянием в Git.

**Q: Чем ArgoCD отличается от GitLab CI для деплоя в K8s?**
GitLab CI — push-модель: pipeline явно вызывает `helm upgrade`. ArgoCD — pull-модель: оператор в кластере следит за Git и сам применяет изменения. ArgoCD лучше для multi-cluster и audit.

**Q: Что такое Kaniko и зачем он нужен?**
Инструмент для сборки Docker-образов внутри K8s Pod без Docker daemon (не нужен privileged mode). Безопаснее docker-in-docker (DinD). Популярен в GitLab CI/CD внутри K8s.

---

## Мониторинг

**Q: Что такое metrics-server?**
Лёгкий агрегатор метрик ресурсов (CPU/memory) с нод и Pods. Нужен для HPA и команды `kubectl top`. Не для долгосрочного хранения — для этого Prometheus.

**Q: Что такое liveness probe с точки зрения операций?**
Способ сказать K8s "если приложение не отвечает на этот endpoint — оно зависло, перезапусти его". Должен проверять только базовую работоспособность (не зависимости), иначе каскадные перезапуски.

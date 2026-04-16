# Kubernetes — Debug, Troubleshooting и Инструменты

---

## 1. Алгоритм диагностики (от общего к частному)

```bash
# 1. Смотрим статус Pods
kubectl get pods -n <namespace>
kubectl get pods -n <namespace> -o wide   # + на какой Node

# 2. Детали и события
kubectl describe pod <pod> -n <ns>

# 3. Логи
kubectl logs <pod> -n <ns>
kubectl logs <pod> -n <ns> --previous    # логи упавшего контейнера (после crash)
kubectl logs <pod> -n <ns> -f            # в реальном времени
kubectl logs <pod> -n <ns> -c <container>  # конкретный контейнер

# 4. События кластера
kubectl get events -n <ns> --sort-by='.lastTimestamp'
kubectl get events --field-selector reason=OOMKilling
```

---

## 2. Статусы Pod и их причины

| Статус | Причина | Что делать |
|---|---|---|
| `ImagePullBackOff` / `ErrImagePull` | Образ не найден или нет доступа к registry | Проверить имя образа, тег, `imagePullSecrets` |
| `ContainerCreating` (зависло) | PVC не смонтирован, Secret/ConfigMap не найден | `kubectl describe pod`, секция Events |
| `Pending` | Нет ресурсов на нодах, не прошёл scheduling | `describe pod` → Events: `Insufficient CPU/memory`, taints |
| `READY 0/1` | readinessProbe не проходит | Лог приложения, проверить health endpoint |
| `CrashLoopBackOff` | Контейнер падает сразу после старта | `kubectl logs --previous`, ошибка в CMD/ENTRYPOINT |
| `Terminating` (завис) | Finalizer не снят, SIGTERM игнорируется | `kubectl delete pod <pod> --force --grace-period=0` |
| Pods нет (0 ready) | Неправильный selector в Deployment/RS | `describe deployment` → проверить `selector.matchLabels` и labels в template |
| Networking issue | Неправильный selector в Service или неверный port | `describe svc` → проверить `selector` и `targetPort` |

---

## 3. Вход в контейнер

```bash
# Стандартный вход (если есть bash/sh)
kubectl exec -it <pod> -n <ns> -- bash
kubectl exec -it <pod> -n <ns> -- sh

# Если bash нет — ephemeral container (без перезапуска Pod)
kubectl debug -it <pod-name> --image=ubuntu -- bash

# Создать копию Pod с debug-образом
kubectl debug <pod-name> -it \
  --copy-to=debug-pod \
  --image=ubuntu \
  --container=debug-container -- sh
```

### exec vs debug

| | kubectl exec | kubectl debug |
|---|---|---|
| Что делает | Запускает команду в существующем контейнере | Добавляет новый ephemeral контейнер |
| Нужен bash в образе? | Да | Нет (образ задаёте сами) |
| Перезапуск Pod | Нет | Нет |
| Когда использовать | Образ с bash (ubuntu, debian) | scratch, distroless, alpine без утилит |

---

## 4. Отладка сети

```bash
# Временный Pod с curl
kubectl run curl -it --image=radial/busyboxplus:curl --restart=Never --rm -- sh
# Внутри:
curl http://goweb-svc.demo.svc.cluster.local:80/health

# Проверить DNS
kubectl run dns-test -it --image=busybox --restart=Never --rm -- nslookup goweb-svc.demo

# Проверить TCP-порт
echo "yes" > /dev/tcp/10.0.11.3/5432 && echo "open" || echo "close"

# Алгоритм: Pod не достучится до Service
# 1. Проверить DNS
# 2. Проверить endpoints (есть ли реальные Pods за Service)
kubectl get endpoints <svc> -n <ns>
# 3. Проверить selector Service
kubectl describe svc <svc> -n <ns>
# 4. Проверить labels Pods
kubectl get pods -l app=myapp -n <ns>
# 5. Проверить NetworkPolicy
kubectl get networkpolicy -n <ns>
```

---

## 5. nsenter — отладка на уровне Node

> Когда `kubectl exec` недоступен (distroless, scratch, нет shell).

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

---

## 6. Ресурсы и Node

```bash
# Потребление CPU/memory
kubectl top pods -n <ns>
kubectl top nodes

# На какой Node стоит Pod
kubectl get pod <pod> -o wide           # колонка NODE
kubectl describe pod <pod> | grep Node

# Вывод ноды из ротации (обслуживание)
kubectl cordon <node>                   # запретить новые Pods
kubectl drain <node> --ignore-daemonsets --delete-emptydir-data  # выселить Pods
kubectl uncordon <node>                 # вернуть в ротацию
```

---

## 7. Алиасы kubectl

```bash
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgpo='kubectl get pods -o wide'
alias kgs='kubectl get secret'
alias kdp='kubectl describe pod'
alias kl='kubectl logs'

# Добавить в ~/.bashrc + source ~/.bashrc
```

---

## 8. Автодополнение

```bash
sudo apt install bash-completion
echo 'source <(kubectl completion bash)' >> ~/.bashrc
echo 'alias k=kubectl' >> ~/.bashrc
echo 'complete -o default -F __start_kubectl k' >> ~/.bashrc
source ~/.bashrc
```

---

## 9. Полезные команды kubectl

```bash
# Получить YAML существующего ресурса
kubectl get pod nginx -o yaml
kubectl get pod nginx -o yaml > nginx.yaml

# Сухой прогон (не применять, только вывести YAML)
kubectl run nginx --image=nginx --dry-run=client -o yaml
kubectl create deployment myapp --image=nginx --dry-run=client -o yaml > deployment.yaml
kubectl expose deployment myapp --port=80 --dry-run=client -o yaml > service.yaml
kubectl create configmap myconfig --from-literal=key=value --dry-run=client -o yaml

# Посмотреть что изменится до apply
kubectl diff -f deployment.yaml

# Смотреть ресурсы в реальном времени
watch kubectl -n home-dev get pods

# Быстрое масштабирование
kubectl -n <ns> scale deployment/myapp --replicas=5

# Документация по полям прямо в терминале
kubectl explain pod.spec.containers
kubectl explain deployment.spec.strategy.rollingUpdate

# Переключение между кластерами и namespace
kubectl config get-contexts
kubectl config use-context <context-name>
kubectl config set-context --current --namespace=my-ns
# или: kubectx (кластеры), kubens (namespace)

# Дерево ресурсов
kubectl krew install tree
kubectl tree deployment my-deploy

# Потребление ресурсов
kubectl krew install resource-capacity
kubectl resource-capacity --pods
```

---

## 10. Быстрые сценарии

### Pod завис в Terminating
```bash
kubectl delete pod <pod> --force --grace-period=0
```

### Port-forward для отладки без Service
```bash
kubectl -n <ns> port-forward pod/mypod 8080:8080
kubectl -n <ns> port-forward svc/myservice 8080:80
```

### Запустить Job вручную из CronJob
```bash
kubectl -n <ns> create job manual-run --from=cronjob/pi-cronjob
```

### Rollback Deployment
```bash
kubectl -n <ns> rollout undo deployment/myapp
kubectl -n <ns> rollout undo deployment/myapp --to-revision=2
kubectl -n <ns> rollout history deployment/myapp
```

# Kubernetes — Сеть и доступ

> Аналогия:
> - **Pod** — сотрудник с рабочим телефоном (IP меняется при переезде).
> - **Service** — корпоративный номер отдела (не меняется, звонки переадресуются на свободного сотрудника).
> - **Ingress** — ресепшн компании (принимает внешние звонки и направляет в нужный отдел по имени).

---

## 1. Service

Service — стабильный сетевой эндпоинт для набора Pods (через selector по labels).

> Зачем нужен? Pod перезапускается → меняется IP.
> Service — это постоянный адрес, который всегда знает актуальные IP Pods.

### Типы Services

| Тип | Доступен снаружи | Протокол | Когда использовать | Аналогия |
|---|---|---|---|---|
| `ClusterIP` | Нет | TCP/UDP | Внутренние микросервисы | Внутренний телефон в офисе |
| `NodePort` | Да (порт ноды 30000–32767) | TCP/UDP | Dev/тесты | Дверь сбоку здания |
| `LoadBalancer` | Да (внешний IP) | TCP/UDP | Production в облаке | Парадный вход с охраной |
| `ExternalName` | — (редирект DNS) | — | Интеграция с внешними сервисами | Переадресация звонка на внешний номер |
| `Ingress` | Да (HTTP/HTTPS) | HTTP/HTTPS | Production HTTP API | Ресепшн с маршрутизацией |

---

### ClusterIP (по умолчанию)

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

### NodePort

```yaml
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
    - port: 80
      targetPort: 8070
      nodePort: 30050   # порт на ноде (30000-32767)
```

Доступ: `http://<node-ip>:30050`

### LoadBalancer

```yaml
apiVersion: v1
kind: Service
metadata:
  name: goweb-svc-lb
spec:
  type: LoadBalancer
  selector:
    app: goapp
  ports:
    - port: 80
      targetPort: 8080
```

---

## 2. Service Discovery — как Pods находят друг друга

> Аналогия: CoreDNS — это **телефонная книга** кластера.
> Каждый Service получает запись: `имя.namespace.svc.cluster.local`.

```
<service-name>.<namespace>.svc.<cluster-domain>
# Пример:
goweb-svc.demo.svc.cluster.local

# Из того же namespace достаточно:
goweb-svc
```

```bash
# Проверка DNS из Pod
kubectl run -it --image=busybox --restart=Never dns-test -- nslookup goweb-svc.demo
```

### Как Service находит Pods

```
Service (selector: app=goapp)
    ↓
Endpoint Controller → следит за Pods с label app=goapp
    ↓
Endpoints object   → актуальный список IP:порт Pods
    ↓
kube-proxy         → настраивает iptables/ipvs правила на каждой ноде
```

```bash
kubectl get endpoints -n demo  # посмотреть реальные IP за Service
```

---

## 3. Ingress

> Аналогия: Ingress — это **ресепшн** компании.
> Один адрес снаружи, а внутри маршрутизирует по имени/пути в нужный отдел.
> LoadBalancer — отдельный вход для каждого отдела (дорого и неудобно).

**LoadBalancer vs Ingress:**

| | LoadBalancer | Ingress |
|---|---|---|
| IP на Service | Один IP на один Service | Один IP на все Services |
| Маршрутизация | По IP | По hostname + path |
| Протокол | TCP/UDP (любой) | HTTP/HTTPS только |
| Стоимость в облаке | Платить за каждый LB | Один LB для всего |
| Нужен контроллер | Облако / MetalLB | IngressController (nginx, traefik) |

**Установка nginx Ingress Controller:**
```bash
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
    - host: api.example.com           # маршрут по hostname
      http:
        paths:
          - path: /users
            pathType: Prefix
            backend:
              service:
                name: users-svc
                port:
                  number: 80
          - path: /orders
            pathType: Prefix
            backend:
              service:
                name: orders-svc
                port:
                  number: 80
```

---

## 4. NetworkPolicy

> Аналогия: NetworkPolicy — это **файрвол между отделами**.
> По умолчанию все Pod могут общаться со всеми.
> NetworkPolicy: "отдел финансов принимает запросы только от отдела бухгалтерии".

**Важно:** работает только если CNI поддерживает (Calico, Cilium, Weave). Flannel — не поддерживает.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-only-frontend
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: backend           # применяется к backend Pods
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend  # разрешить только от frontend
      ports:
        - port: 8080
```

---

## 5. CoreDNS и headless Service

### CoreDNS
- Системный Pod в `kube-system`.
- Создаёт A-записи для каждого Service.
- Настраивается через ConfigMap `coredns`.

```bash
# Добавить внутреннее DNS-имя
kubectl -n kube-system edit cm coredns
# Добавить в NodeHosts: 10.0.2.15 my-internal-service.example.com

# Перезапустить CoreDNS
kubectl -n kube-system delete pod -l k8s-app=kube-dns
```

### Headless Service

> Аналогия: обычный Service — звонок в отдел (ответит любой).
> Headless Service — звонок конкретному сотруднику по имени (pod-0, pod-1...).

```yaml
spec:
  clusterIP: None   # headless — нет виртуального IP
  selector:
    app: mydb
```

Используется в StatefulSet: `pod-0.mydb-svc.ns.svc.cluster.local` — стабильный адрес конкретной реплики.

---

## 6. Как Pod получает IP

```
Pod создан
   ↓
kubelet вызывает CNI-плагин
   ↓
CNI выделяет IP из Pod CIDR (пул адресов для Pods)
   ↓
Настраивает сетевой интерфейс
   ↓
Pod получает уникальный IP в кластере
```

**Важно:** IP Pod непостоянный — меняется при пересоздании. Поэтому всегда обращайтесь через Service, не через IP Pod.

---

## 7. Utility Pod для отладки сети

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
kubectl exec -it curl -- curl http://goweb-svc.demo.svc.cluster.local/health

# Проверить DNS
kubectl run dns-test -it --image=busybox --restart=Never --rm -- nslookup goweb-svc.demo

# Проверить TCP-порт
echo "yes" > /dev/tcp/10.0.11.3/5432 && echo "open" || echo "close"
```

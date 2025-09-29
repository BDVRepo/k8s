# Gitlab Runner

```sh
export PS1='\u@\h:\W\$ '
```

1. Добавляем репозиторий с чартом
```sh
helm repo add gitlab https://charts.gitlab.io
helm repo update gitlab
```

2. Создаем Неймспейс
```sh
kubectl create ns gitlab-runner
```

3. Скачиваем сертификат и деплоим Секрет с сертификатом
```sh
openssl s_client -showcerts -connect gitlab.example.com:443 -servername gitlab.example.com < /dev/null 2>/dev/null | openssl x509 -outform PEM > gitlab.example.com.crt

kubectl -n gitlab-runner create secret generic gitlab-tls --from-file gitlab.example.com.crt
# kubectl -n gitlab-runner create secret generic gitlab-tls --from-file gitlab.example.com.crt --dry-run=client -o yaml > secret.yaml
```

4. Устанавливаем релиз с Гитлаб Раннером
```sh
helm upgrade --install -n gitlab-runner gitlab-runner gitlab/gitlab-runner -f values.yaml
```

## DNS issue


Если у Гитлаба произвольное доменное имя, не прописанное нигде в DNS, то Раннер не сможет запуститься.

Это будет видно по сообщению "no such host" в логах
```sh
kubectl -n gitlab-runner logs gitlab-runner-57444c949f-jd5m4
```

```
ERROR: Verifying runner... failed                   runner=t1_xmVwE2 status=couldn't execute POST against https://gitlab.example.com/api/v4/runners/verify: Post "https://gitlab.example.com/api/v4/runners/verify": dial tcp: lookup gitlab.example.com on 10.43.0.10:53: no such host
PANIC: Failed to verify the runner.  
```
1. Отредактируйте ConfigMap Core DNS: 
```sh
kubectl -n kube-system edit cm coredns
```

В разделе NodeHosts добавьте запись:
```
  NodeHosts: |
    10.0.2.15 ubuntu2004                                                                                     
    178.20.45.22 gitlab.example.com
```

3. Перезагрузите Core DNS:
```sh
kubectl -n kube-system delete po coredns-77ccd57875-hsw9t
```

4. С Ноды k3s можно проверить, что доменное имя резолвится успешно: 
```sh
dig @10.43.0.10 gitlab.example.com
```

Где 10.43.0.10 - это IP адрес Сервиса Core DNS

```sh
kubectl -n kube-system get svc kube-dns
```
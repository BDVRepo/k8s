# Gitlab Runner

1. Добавляем репозиторий с чартом Гитлаба: 
```sh
helm repo add gitlab https://charts.gitlab.io
helm repo update gitlab
```

2. Создаем Неймспейс  для Раннера:
```sh
kubectl create ns gitlab-runner
```

3. В gitlab.com получите Токен для будущего Раннера и заполните values.yaml

4. Устанавливаем Раннер: 
```sh
helm upgrade --install -n gitlab-runner gitlab-runner gitlab/gitlab-runner -f values.yaml
```
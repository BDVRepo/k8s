## Commands

```bash
cd administrator-kubernetes-junior/02_Docker/023-step2

# Собираем обычный образ
docker build . -t go-big:0.0.1

# Собираем образ, используя прием Multistage. Через ключ -f передается путь к Dockerfile
docker build . -t go-tiny:0.0.1 -f Dockerfile_multistage

docker images 

# Также Вы можете проверить, что приложение работает
docker run -p 8080:8080 go-tiny:0.0.1
curl localhost:8080/k8s
```

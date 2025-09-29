## HTTP-server

Accepts requests and log URL path.

#### Сборка с Docker 
```bash
docker build . -t bdv21/payment-server-test1:0.0.1
docker push bdv21/payment-server-test1:0.0.1
```

#### Сборка с Kaniko

```sh
export REGISTRY_URL="https://index.docker.io/v1/"
export REGISTRY_USERNAME="[REDACTED]"
export REGISTRY_PASSWORD="[REDACTED]"

mkdir  /tmp/.docker
sudo chmod 777 /tmp/.docker

echo "{\"auths\":{\"$REGISTRY_URL\":{\"auth\":\"$(echo -n "${REGISTRY_USERNAME}:${REGISTRY_PASSWORD}" | base64 | tr -d '\n')\"}}}" > /tmp/.docker/config.json
```

```sh
docker run -v .:/workspace/ \
           -v /tmp/.docker/config.json:/kaniko/.docker/config.json \
           gcr.io/kaniko-project/executor:v1.23.2 \
           --destination bdv21/payment-server-test2:0.0.1
```

#### unit-test

```sh
go test ./...
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile.cov
```

#### Деплой с Helm

```sh
helm repo add m8x https://MadEngineX.github.io/helm-charts/
helm repo update

kubectl create ns payment-server-dev

helm -n payment-server-dev upgrade --install payment-server m8x/common-chart -f deploy/dev/values.yaml
```

Проверка:
```sh
echo '192.168.10.1 payment-server-dev.akj.ru' | sudo tee -a /etc/hosts
curl payment-server-dev.akj.ru/helm
```
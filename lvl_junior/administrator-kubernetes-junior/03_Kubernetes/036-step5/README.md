# UpdateStrategy

```bash
kubectl -n home-dev get deploy go-http-server -o yaml

kubectl explain deployment
kubectl explain deployment.spec
kubectl explain deployment.spec.strategy
kubectl explain deployment.spec.strategy.rollingUpdate

watch kubectl -n home-dev get po
kubectl apply -f deployment.yaml
```
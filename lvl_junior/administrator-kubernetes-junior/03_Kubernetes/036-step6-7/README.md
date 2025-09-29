# Deployment management

```bash
kubectl -n home-dev get deploy
kubectl -n home-dev get po 
kubectl -n home-dev get po  go-http-server-7cd9cb8c9-tttj5 -o yaml | grep -i image
ubectl -n home-dev rollout history deployment/go-http-server
kubectl -n home-dev get rs
watch kubectl -n home-dev get po

kubectl -n home-dev rollout undo deployment/go-http-server
kubectl -n home-dev get pods go-http-server-7cd9cb8c9-wlqz4 -o yaml | grep -i image
kubectl -n home-dev rollout restart deployment/go-http-server
kubectl -n home-dev delete po go-http-server-65896f556c-k8b75
kubectl -n home-dev scale deployment/go-http-server --replicas=1
kubectl diff -f deployment.yaml 
```
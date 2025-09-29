## Pod

```bash
kubectl apply -f kubernetes/pod.yaml 
kubectl get pods
kubectl --namespace home-dev get pods
kubectl -n home-dev get pods
kubectl -n home-dev port-forward pod/http-server 8080:8080
kubectl delete -f kubernetes/pod.yaml 
```
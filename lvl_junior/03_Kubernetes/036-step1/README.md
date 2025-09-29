## Upgpdate with ReplicaSet

```bash
cd http-server/
docker build . -t bdv21/http-server:0.0.2
docker run -it bdv21/http-server:0.0.2
docker login
docker push bdv21/http-server:0.0.2

kubectl -n home-dev get pod -w
kubectl apply -f replicaSet.yaml 
kubectl -n home-dev get pod http-server-5rvdk -o yaml | grep -i image
kubectl -n home-dev delete pods http-server-5rvdk http-server-xgxsr
kubectl -n home-dev get pods http-server-bdwxl -o yaml | grep -i image
kubectl -n home-dev logs http-server-bdwxl
```
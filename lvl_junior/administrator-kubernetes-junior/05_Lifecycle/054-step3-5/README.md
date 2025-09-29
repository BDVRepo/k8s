## Probes

```bash
docker build . -t bdv21/healthchecks:1.0.0
docker push  bdv21/healthchecks:1.0.0

kubectl apply -f 01_deployment.yaml 
watch kubectl -n probes-test get po
kubectl apply -f 02_deployment.yaml 

kubectl -n probes-test scale deployment/go-http-server --replicas=0
kubectl run curl --image=radial/busyboxplus:curl -i --tty --rm
curl http://healthchecks-svc.probes-test:8080/hi
kubectl -n probes-test scale deployment/go-http-server --replicas=1
curl http://healthchecks-svc.probes-test:8080/hi
```
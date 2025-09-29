# Job

```bash
kubectl apply -f job.yaml 
kubectl -n home-dev get po
kubectl -n home-dev get job 
kubectl -n home-dev logs echo-job-alpine-4p77s
kubectl explain job.spec.backoffLimit
```


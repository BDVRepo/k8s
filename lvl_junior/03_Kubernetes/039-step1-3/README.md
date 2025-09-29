# CronJob

```bash
watch kubectl -n home-dev get po

kubectl apply -f cronJob.yaml 
kubectl -n home-dev get cronjob
kubectl -n home-dev get jobs

kubectl -n home-dev describe job hello-28373589
kubectl -n home-dev describe cronjob hello
kubectl -n home-dev create job --from=cronjob/hello hello-manual
```
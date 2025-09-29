# Secret. Mount as file

```bash
kubectl -n read-secret exec -it $(kubectl -n read-secret get pod -l "app=go-file-reader" -o name) -- df -h
```
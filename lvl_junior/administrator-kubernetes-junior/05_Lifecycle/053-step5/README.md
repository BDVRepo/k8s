# Memory Limit

```bash
docker login -u kinualx@gmail.com

docker build . -t bdv21/memory-bomb:1.0.0
docker push bdv21/memory-bomb:1.0.0

kubectl -n memory-test top pods
```
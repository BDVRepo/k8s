## Skooner

https://skooner.io/

#### Installation

```sh
kubectl apply -f skooner.yaml
```

```sh
## kubectl -n skooner create sa skooner-sa

kubectl apply -f crb.yaml 
kubectl -n skooner create token skooner-sa
```
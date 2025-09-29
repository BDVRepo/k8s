# Declaration

```bash
kubectl apply -f kubernetes/namespace.yaml 
kubectl get ns
kubectl create ns wallet-dev2
kubectl get ns
kubectl run nginx --image=nginx --restart=Never
kubectl -n default get po
```
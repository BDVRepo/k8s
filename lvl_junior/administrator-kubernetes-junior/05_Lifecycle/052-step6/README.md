# LoadBalancer

```bash
kubectl apply -f deployment.yaml 
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.12/config/manifests/metallb-native.yaml
kubectl apply -f pool.yaml

sudo vim /etc/systemd/system/k3s.service

# ExecStart=/usr/local/bin/k3s \
#    server \
#        '--write-kubeconfig-mode' \
#        '644' \
#        '-disable=servicelb'

systemctl daemon-reload && systemctl restart k3s
kubectl get nodes -w

kubectl apply -f svc.yaml

curl http://192.168.10.14:8080/cakes
```
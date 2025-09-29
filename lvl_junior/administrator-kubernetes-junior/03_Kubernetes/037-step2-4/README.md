# DaemonSet

```bash
## sudo k3s server --write-kubeconfig-mode 644 --token K102bd1edc6a58dbe79a5199bbf97c75f003d70365d694e46bba8979fa68b2a3c16::server:fab2c838bd493639e71c25136f498a66 

sudo cat /var/lib/rancher/k3s/server/node-token
kubectl get nodes -o wide

docker run  --privileged  -it -p 6445:6445 rancher/k3s:v1.24.10-k3s1   agent --server https://10.0.2.15:6443 --token K102bd1edc6a58dbe79a5199bbf97c75f003d70365d694e46bba8979fa68b2a3c16::server:fab2c838bd493639e71c25136f498a66 --lb-server-port 6445 --node-name worker1
```


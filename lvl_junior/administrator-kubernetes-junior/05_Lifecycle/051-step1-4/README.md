# Service Discovery

Get cluster name: 
```bash
kubectl run -it --image=ubuntu --restart=Never shell -- \
sh -c 'apt-get update > /dev/null && apt-get install -y dnsutils > /dev/null && \
nslookup kubernetes.default | grep Name | sed "s/Name:\skubernetes.default//"'
```

Then delete Pod in `default` namespace.
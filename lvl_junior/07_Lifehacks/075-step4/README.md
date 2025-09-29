## OKD console

https://github.com/openshift/console

https://gitlab.com/av1o/charts/-/tree/master/charts/openshift-console

#### Installation

```sh
kubectl apply -f 01_ns.yaml

helm repo add av1o https://av1o.gitlab.io/charts
helm -n console upgrade --install okd-console av1o/openshift-console -f values.yaml

echo '192.168.10.1 console.akj.ru' | sudo tee -a /etc/hosts

kubectl apply -f 02_crb.yaml
```

#### Integration with Keycloak

https://github.com/openshift/console/issues/9754#issuecomment-2071734589
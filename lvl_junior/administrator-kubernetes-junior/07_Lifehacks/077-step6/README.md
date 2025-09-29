## Debug

#### No Pods exists

```sh
cd 077-step6
kubectl apply -f .
echo '192.168.10.1 debug-6.akj.ru' | sudo tee -a /etc/hosts

kubectl krew install tree
```

Browse http://debug-6.akj.ru/debugging

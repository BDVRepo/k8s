## Debug

#### Networking issue

```sh
cd 077-step7
kubectl apply -f .
echo '192.168.10.1 debug-7.akj.ru' | sudo tee -a /etc/hosts
```

Browse http://debug-7.akj.ru/debugging

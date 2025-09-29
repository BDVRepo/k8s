## Debug

#### READY 0/1

```sh
cd 077-step4
kubectl apply -f .
echo '192.168.10.1 debug-4.akj.ru' | sudo tee -a /etc/hosts
```

Browse http://debug-4.akj.ru/debugging

## Debug

#### ContainerCreating 

```sh
cd 077-step2
kubectl apply -f .
echo '192.168.10.1 debug-2.akj.ru' | sudo tee -a /etc/hosts
```

Browse http://debug-2.akj.ru/debugging

## Debug

#### ImagePuullBackoff (ErrImagePull)

```sh
cd 077-step1
kubectl apply -f .
echo '192.168.10.1 debug-1.akj.ru' | sudo tee -a /etc/hosts
```

Browse http://debug-1.akj.ru/debugging

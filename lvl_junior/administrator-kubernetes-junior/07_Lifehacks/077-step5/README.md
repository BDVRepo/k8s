## Debug

#### READY 0/1 RESTARTS 5 (29s ago) (CrashLoopBackOff)

```sh
cd 077-step5
kubectl apply -f .
echo '192.168.10.1 debug-5.akj.ru' | sudo tee -a /etc/hosts
```

Browse http://debug-5.akj.ru/debugging

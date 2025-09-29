## Debug

#### kubectl debug

```sh
cd 078-step1
kubectl apply -f .
```

```sh
telnet 10.0.11.3 5432

echo "yes" > /dev/tcp/10.0.11.3/5432 && echo "open" || echo "close"
```

```sh
kubectl debug -it --image=ubuntu  go-http-server-6568d89d97-gkwq7
```

```sh
kubectl debug go-http-server-6568d89d97-dmxvn -it --copy-to=my-debugger --image=ubuntu --container=mycontainer -- sh
```

Browse http://debug-8.akj.ru/debugging

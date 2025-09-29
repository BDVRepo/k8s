## Create your own kubectl plugin

#### Bash

```sh
chmod +x kubectl-access.sh
sudo cp kubectl-access.sh /usr/local/bin/kubectl-access
```

```sh
kubectl access
```

#### Golang

```sh
cd kubectl-access
go build -o kubectl-access main.go
chmod +x kubectl-access
sudo mv kubectl-access /usr/local/bin/kubectl-access
```

```sh
kubectl access
```
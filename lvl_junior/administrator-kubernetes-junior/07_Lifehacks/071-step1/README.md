## alias

```bash
echo 'alias k="kubectl"' >> ~/.bashrc
source ~/.bashrc
```

```bash
alias k='kubectl'
alias kgp='kubectl get pods'
alias kgpo='kubectl get pods -o wide'
alias kgs='kubectl get secret'
alias kdp='kubectl describe pod'
alias kl='kubectl logs'
```

```bash
kgs -n weather-dev weather-bot-env -o yaml
kubectl --context k3s  get secret -n weather-dev weather-bot-env -o yaml
```
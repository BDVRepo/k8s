## Commands completition

https://kubernetes.io/docs/reference/kubectl/generated/kubectl_completion/

#### Ubuntu

```sh
sudo apt update && sudo apt install bash-completion
type _init_completion
# source '/usr/share/bash-completion/bash_completion' >>~/.bashrc

echo 'source <(kubectl completion bash)' >>~/.bashrc

echo 'alias k=kubectl' >>~/.bashrc
echo 'complete -o default -F __start_kubectl k' >>~/.bashrc

source ~/.bashrc
```


#### Mac OS

```sh
brew install bash-completion@2
## Добавить в ~/.bash_profile 
## brew_etc="$(brew --prefix)/etc" && [[ -r "${brew_etc}/profile.d/bash_completion.sh" ]] && . "${brew_etc}/profile.d/bash_completion.sh"
## Перезагрузить терминал
echo 'source <(kubectl completion bash)' >>~/.bash_profile

kubectl completion bash >/usr/local/etc/bash_completion.d/kubectl

echo 'alias k=kubectl' >>~/.bash_profile
echo 'complete -o default -F __start_kubectl k' >>~/.bash_profile
```
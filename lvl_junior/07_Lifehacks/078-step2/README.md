## Debug


#### nsenter

```sh
cd 078-step2
kubectl apply -f .

kubectl get po -o wide

## ssh on node with pod

sudo crictl ps | grep go-http-server-6568d89d97-rmkfc ## pod name

sudo crictl inspect --output go-template --template '{{.info.pid}}'  0d3927393a30d ## container ID

# nsenter -t <pid of container> -n <command and args>
top
sudo nsenter -t 124792 -m top

sudo nsenter -t 124792 -n netstat -tulpn | grep LISTEN

sudo nsenter -t 124792 -m dig
sudo nsenter -t 124792 -n dig google.com
sudo nsenter -t 124792 -m cat /etc/resolv.conf
sudo nsenter -t 124792 -n dig @10.43.0.10 google.com

sudo nsenter -t 124792 -n ip link show
```

https://github.com/pabateman/kubectl-nsenter

#### init 1 node cluster

VM Ubuntu 22.04
2 CPU 
4 GB RAM
100 GB HDD

93.183.92.134

```sh
su -

## Disable swap
swapoff -a
swapoff -a && sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

## Core modules for containerd
echo overlay >> /etc/modules-load.d/containerd.conf 
echo br_netfilter >> /etc/modules-load.d/containerd.conf 

modprobe overlay
modprobe br_netfilter

## Check if modules work
lsmod | egrep "br_netfilter|overlay"

echo 'net.bridge.bridge-nf-call-iptables = 1' >> /etc/sysctl.d/99-kubernetes-cri.conf
echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.d/99-kubernetes-cri.conf
echo 'net.bridge.bridge-nf-call-ip6tables = 1' >> /etc/sysctl.d/99-kubernetes-cri.conf

sysctl --system

apt-get update
apt-get install -y apt-transport-https ca-certificates curl software-properties-common

# Add the Docker repository
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# Add the Kubernetes repository
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.30/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

apt-get update
apt-get install -y docker-ce docker-ce-cli containerd

apt-get update
apt-get install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl

apt-get install -y cri-tools socat conntrack ipvsadm ebtables git curl wget runc


iptables -P FORWARD ACCEPT
ufw allow 6443/tcp


containerd config default | sudo tee /etc/containerd/config.toml 

sed -i 's/SystemdCgroup \= false/SystemdCgroup \= true/g' /etc/containerd/config.toml

systemctl enable --now containerd
systemctl enable kubelet.service

kubeadm config images pull

hostnamectl set-hostname masternode

echo 93.183.92.134 masternode >> /etc/hosts

## Init cluster
kubeadm init --pod-network-cidr=10.244.0.0/16

echo "export KUBECONFIG=/etc/kubernetes/admin.conf" >> /root/.bashrc
source .bashrc
export KUBECONFIG=/etc/kubernetes/admin.conf

kubectl taint nodes --all node-role.kubernetes.io/control-plane-

kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.0/manifests/calico.yaml

kubectl get pod -n kube-system

## probably need
vim /lib/systemd/system/kubelet.service
## add Environment="KUBELET_EXTRA_ARGS=--hostname-override=masternode" in [Service] section
```


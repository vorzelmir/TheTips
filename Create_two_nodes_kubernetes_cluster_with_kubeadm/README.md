### Create two nodes Kubernetes cluster with kubeadm

control-plane node OS is 

`Fedora Linux 37 (Workstation Edition)`

worker-node OS is

`Fedora Linux 41 (Server Edition)`

nodes's hostnames must be different the control-plane is  __fedora__ and the worker node is __fedora-node__ in this case

The kubernetes is 1.31 version, crio version 1.24.1

#### Control-plane and worker node configure

enable crio, kubelet

`fedora<fedora-node># systemctl enable --now crio`

`fedora<fedora-node># systemctl enable --now kubelet`

turn off swap

`fedora<fedora-node># dnf remove zram-generator-defaults`

and 

`fedora<fedora-node># swapoff -a`

set SELinux permissive mode

`fedora<fodora-node># setenforce 0`

#### Control-plane node configure

configure firewall

`fedora# firewall-cmd --add-port={6443,2379,2380,10250,10259,10257}/tcp`

`fedora# firewall-cmd --add-port={6443,2379,2380,10250,10259,10257}/tcp --permanent`

init control-plane node with crio default container runtime interface

`fedora# kubeadm init --cri-socket unix:///var/run/crio/crio.sock`

create in $HOME .kube directory and copy admin.conf to it

`fedora$ mkdir $HOME/.kube`

`fedora# cp /etc/kubernetes/admin.conf /home/USER/.kube/conf`

make USER owner of it

`fedora$ sudo chown $USER:$USER $HOME/.kube/config`

get created node

```
fedora$ kubectl get node
NAME          STATUS   ROLES           AGE   VERSION
fedora        Ready    control-plane   2m   v1.31.7

```

run Calico on the cluster

`fedora$ curl https://raw.githubusercontent.com/projectcalico/calico/v3.29.3/manifests/calico.yaml -O
`

`fedora$ kubectl apply -f calico.yaml`

#### Worker node configure

firewall settings

`fedora-node# firewall-cmd --add-port={10250,10256,30000-32767}/tcp`

`fedora-node# firewall-cmd --add-port={10250,10256,30000-32767}/tcp --permanent`

create a directory and copy from the control-plane node config file to it

`fedora-node$ mkdir $HOME/.kube`

`fedora# scp /etc/kubernetes/admin.conf fedora-node.ip:/home/USER/.kube/config`

make USER owner of it

`fedora-node$ sudo chown $USER:$USER $HOME/.kube/config`

create new node

`fedora-node# kubeadm join --discovery-file /home/USER/.kube/config --cri-socket unix:///var/run/crio/crio.sock`

make global variable 

`fedora-node$ export KUBECONFIG=/home/USER/.kube/config`

get nodes on any node

```fedora<fedora-node>$ kubectl get nodes
NAME          STATUS   ROLES           AGE     VERSION
fedora        Ready    control-plane   26m     v1.31.7
fedora-node   Ready    <none>          2m48s   v1.31.7
```















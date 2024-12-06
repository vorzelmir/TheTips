### Run kubernetes cluster on-premises environment on RHEL-9.4

operation system on the virtual machine

`localhost$ hostanamectl`

```
 Static hostname: rhel
       Icon name: computer-vm
         Chassis: vm 🖴
      Machine ID: bf792a7bfce2455dab8abb477297845b
         Boot ID: 31971a1929c84d47b4c474153e8c89fc
  Virtualization: kvm
Operating System: Red Hat Enterprise Linux 9.4 (Plow)     
     CPE OS Name: cpe:/o:redhat:enterprise_linux:9::baseos
          Kernel: Linux 5.14.0-427.13.1.el9_4.x86_64
    Architecture: x86-64
 Hardware Vendor: QEMU
  Hardware Model: Standard PC _Q35 + ICH9, 2009_
Firmware Version: 1.16.2-1.fc37

```

add tcp ports to run api-server on kubernetes cluster


`localhost# firewall-cmd --add-port={6443/tcp,2379-2380/tcp,10250/tcp,10259/tcp,10257/tcp}` 

`localhost# firewall-cmd --add-port={6443/tcp,2379-2380/tcp,10250/tcp,10259/tcp,10257/tcp} --permanent` 

next steps from Kubernetes docks, firstly add kubernetes repo 

`localhost# vim /etc/yum.repos.d/kubernetes.repo`

and add next content

```
[kubernetes]
name=Kubernetes
baseurl=https://pkgs.k8s.io/core:/stable:/v1.31/rpm/
enabled=1
gpgcheck=1
gpgkey=https://pkgs.k8s.io/core:/stable:/v1.31/rpm/repodata/repomd.xml.key
```

this step is not obvious but needed(in my case), commented the next line in the /etc/containerd/config.toml

`#disabled_plugins = ["cri"]`

install tools

`localhost# yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes`

enable and start kubelet.service

`localhost# systemctl enable --now kubelet`

check status and enable if needed containerd daemon

`localhost# systemctl status containerd`

if active disable system swap and make changes in /etc/fstab to comment swap line

`localhost# swapoff -a`

run kubernetes cluster

`localhost# kubeadm init`

folowing the instructions copy the config file to the $USER home .kube directory

`localhost# cp -i /etc/kubernetes/admin.conf $HOME/.kube/config`

finally, check the result

`localhost$ kubectl get all`

```
NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   24m

```


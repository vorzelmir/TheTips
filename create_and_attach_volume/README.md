### Create volume and add it to virtual machine

with qemu-img

`localhost# qemu-img create -f qcow2 /var/lib/libvirt/images/qcow_disk.qcow2 3G `

or with **virsh** tools

create dump xml file from existing on localhost machine another virtual disk and edit the tag <name> 

`localhost# virsh vol-dumpxml --pool default vol-alloc | sed '/name/s/vol-alloc/new-vol/' > new-vol.xml`

create new virtual disk from file

`localhost$ virsh vol-create default new-vol.xml`

where default is a pool name and second parameter is the path to xml file



check creation of disks

`localhost# ls -l /var/lib/libvirt/images`



get the info of virtual disk

`localhost# virsh vol-info --pool default new-vol`

```
Name:           new-vol
Type:           file
Capacity:       2.00 GiB
Allocation:     2.00 GiB

```

get the available virtual machines

`localhost# virsh list --all`

```
 Id   Name          State
------------------------------
 -    eve           shut off
 -    fedora38-2    shut off
 -    ol9.3         shut off
 -    rhel9.3       shut off
```

run Virtual Machine based on Oralcle Linux Server

`localhost# virsh start ol9.3`

connect to VM

`localhost$ virt-viewer --connect qemu:///system`

or connect to VM throw ssh

`localhost$ ssh username@ol9.3-ip-address`



linux server on VM is

`VM$ hostnamectl | grep Operating`

```
Operating System: Oracle Linux Server 9.3
```

check existing virtual disks

`ol9.3$ lsblk`

```
NAME                MAJ:MIN RM   SIZE RO TYPE MOUNTPOINTS
sr0                  11:0    1  1024M  0 rom  
vda                 251:0    0    20G  0 disk 
├─vda1              251:1    0   512M  0 part /boot
└─vda2              251:2    0    11G  0 part 
  ├─ol-root         252:0    0    10G  0 lvm  /
  └─ol-swap         252:1    0     1G  0 lvm  [SWAP]
vdb                 251:16   0     2G  0 disk 
└─vdb1              251:17   0     2G  0 part 
  └─vg_home-lv_home 252:2    0   1.5G  0 lvm  /home

```
attach disk to the VM

`localhost# virsh attach-disk --domain ol9.3 --source /var/lib/libvirt/images/new-vol --target vdc --persistent`

check the appeareance of new virtual disk

`ol9.3$ lsblk`

```
NAME                MAJ:MIN RM   SIZE RO TYPE MOUNTPOINTS
sr0                  11:0    1  1024M  0 rom  
vda                 251:0    0    20G  0 disk 
├─vda1              251:1    0   512M  0 part /boot
└─vda2              251:2    0    11G  0 part 
  ├─ol-root         252:0    0    10G  0 lvm  /
  └─ol-swap         252:1    0     1G  0 lvm  [SWAP]
vdb                 251:16   0     2G  0 disk 
└─vdb1              251:17   0     2G  0 part 
  └─vg_home-lv_home 252:2    0   1.5G  0 lvm  /home
vdc                 251:32   0     2G  0 disk
```

if need to detach disk back

`localhost# virsh detach-disk --domain ol9.3 --target vdc --persistent`

check the result

```
NAME                MAJ:MIN RM  SIZE RO TYPE MOUNTPOINTS
sr0                  11:0    1 1024M  0 rom  
vda                 251:0    0   20G  0 disk 
├─vda1              251:1    0  512M  0 part /boot
└─vda2              251:2    0   11G  0 part 
  ├─ol-root         252:0    0   10G  0 lvm  /
  └─ol-swap         252:1    0    1G  0 lvm  [SWAP]
vdb                 251:16   0    2G  0 disk 
└─vdb1              251:17   0    2G  0 part 
  └─vg_home-lv_home 252:2    0  1.5G  0 lvm  /home
```

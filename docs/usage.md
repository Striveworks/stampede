## Ubuntu
The standard install will install stampede and it's dependencies. A systemd service
will be created for stampede and run upon install.
```
git clone git@github.com:Striveworks/stampede.git
cd stampede
```
MicroK8s
```
make install-microk8s
```
Kubeadm
```
make install-kubeadm
```


## Vagrant
Vagrant can be used test stampede or run on a set of VM's.

[Install Vagrant](https://www.vagrantup.com/docs/installation)

```
git clone git@github.com:Striveworks/stampede.git
cd stampede
```
MicroK8s
```
export STAMPEDE_ENV=microk8s && vagrant up
```
Kubeadm
```
export STAMPEDE_ENV=kubeadm && vagrant up
```

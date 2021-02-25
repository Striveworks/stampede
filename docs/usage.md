## Ubuntu
The standard install will install go, microk8s and stampede. A systemd service
will be created for stampede and run upon install.
```
git clone git@github.com:Striveworks/stampede.git
cd stampede
make install
```


## Vagrant
Vagrant can be used test stampede or run on a set of VM's.

[Install Vagrant](https://www.vagrantup.com/docs/installation)

```
git clone git@github.com:Striveworks/stampede.git
cd stampede
vagrant up
```

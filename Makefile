test:
	vagrant up
	vagrant ssh k8s1 -- -t 'sudo -s; systemctl status stampede; microk8s.kubectl get nodes'

clean:
	vagrant destroy --force

install:
	bash ./scripts/install

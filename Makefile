go-test:
	cd pkg && go test

test-microk8s:
	export STAMPEDE_ENV=microk8s
	vagrant up
	vagrant ssh k8s1 -- -t 'sudo systemctl status stampede; sudo microk8s.kubectl get nodes'

test-kubeadm:
	export STAMPEDE_ENV=kubeadm
	vagrant up
	vagrant ssh k8s1 -- -t 'sudo systemctl status stampede'

clean:
	vagrant destroy --force

docs:
	bash ./scripts/docs.sh

install-microk8s:
	bash ./scripts/install-microk8s.sh

install-kubeadm:
	bash ./scripts/install-kubeadm.sh

go-test:
	cd pkg && go test
	
test:
	vagrant up
	vagrant ssh k8s1 -- -t 'sudo systemctl status stampede; sudo microk8s.kubectl get nodes'

clean:
	vagrant destroy --force

docs:
	bash ./scripts/docs.sh

install:
	bash ./scripts/install.sh

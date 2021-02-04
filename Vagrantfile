Vagrant.configure("2") do |config|

  config.vm.provision "shell", path: "scripts/install.sh"

  config.vm.synced_folder ".", "/opt/stampede"

  config.vm.define "k8s1" do |k8s1|
    k8s1.vm.hostname = "k8s1"
    k8s1.vm.box = "ubuntu/bionic64"
    k8s1.vm.network "private_network", type: "dhcp"
  end

  config.vm.define "k8s2" do |k8s2|
    k8s2.vm.hostname = "k8s2"
    k8s2.vm.box = "ubuntu/bionic64"
    k8s2.vm.network "private_network", type: "dhcp"
  end

  config.vm.define "k8s3" do |k8s3|
    k8s3.vm.hostname = "k8s3"
    k8s3.vm.box = "ubuntu/bionic64"
    k8s3.vm.network "private_network", type: "dhcp"

  end

end

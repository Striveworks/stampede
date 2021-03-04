Vagrant.configure("2") do |config|

 if (ENV['STAMPEDE_ENV'] == "kubeadm")
    config.vm.provision "shell" do |s|
       s.path = "scripts/install-kubeadm.sh"
       s.args = "enp0s8"
    end

    config.vm.synced_folder ".", "/opt/stampede"

    config.vm.provider "virtualbox" do |v|
        v.memory = 2048
        v.cpus = 2
    end

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

  if (ENV['STAMPEDE_ENV'] == "microk8s")
     config.vm.provision "shell" do |s|
        s.path = "scripts/install-microk8s.sh"
        s.args = "enp0s8"
     end

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
end

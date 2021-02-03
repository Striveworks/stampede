#/bin/env bash
set -ex


function setup_user {
  # Setup Docker for non-root user
   groupadd -f docker

  # Create RKE User
   useradd $RKE_USER \
  --groups docker \
  --create-home \
  --password $RKE_PASS \
  || true

  passwd $RKE_USER --stdin <<< "$RKE_PASS"
  # Add RKE User to sudoers
  grep -qxF "$RKE_USER    ALL=(ALL:ALL) ALL"  /etc/sudoers || echo "$RKE_USER    ALL=(ALL:ALL) ALL"  >> /etc/sudoers

}

function create_ssh_keys {
  apt-get install sshpass -y
  rm -rf /home/$RKE_USER/.ssh/
  mkdir -p /home/$RKE_USER/.ssh/
  mkdir -p $HOME/.ssh
  ssh-keygen -b 2048 -t rsa -f /home/$RKE_USER/.ssh/id_rsa -q -N ""
  cat /home/$RKE_USER/.ssh/id_rsa.pub >> /home/$RKE_USER/.ssh/authorized_keys
  chown $RKE_USER:$RKE_USER /home/$RKE_USER/.ssh/id_rsa
  chown $RKE_USER:$RKE_USER /home/$RKE_USER/.ssh/id_rsa.pub
  chown $RKE_USER:$RKE_USER /home/$RKE_USER/.ssh/authorized_keys

  for host in $1 $2 $3
    do
    if [ "$(hostname -f)" != "$host" ];
    then
      sshpass -f /root/${RKE_USER}-password.txt ssh-copy-id -i /home/$RKE_USER/.ssh/id_rsa.pub $RKE_USER@$host
    fi
    done
}

function disable_swap {
  # Disable Swap
  swapoff -a
  sed -i '/ swap / s/^/#/' /etc/fstab

}


function kernel_setup {
  apt-get update
  apt-get dist-upgrade -y

  # Check all kernel modules are present
  for module in br_netfilter ip6_udp_tunnel ip_set ip_set_hash_ip ip_set_hash_net iptable_filter iptable_nat iptable_mangle iptable_raw nf_conntrack_netlink nf_conntrack nf_conntrack_ipv4   nf_defrag_ipv4 nf_nat nf_nat_ipv4 nf_nat_masquerade_ipv4 nfnetlink udp_tunnel veth vxlan x_tables xt_addrtype xt_conntrack xt_comment xt_mark xt_multiport xt_nat xt_recent xt_set  xt_statistic xt_tcpudp;
       do
         modprobe $module
         if ! lsmod | grep -q $module; then
           echo "module $module is not present";
           exit 1
         fi;
       done


  echo "br_netfilter" > /etc/modules-load.d/bridge.conf
  # Ensure net.bridge.bridge-nf-call-iptables is enabled in the kernel
  sysctl net.bridge.bridge-nf-call-iptables=1
  sysctl net.bridge.bridge-nf-call-ip6tables=1
  sysctl net.bridge.bridge-nf-call-arptables=1

}


function install_docker {
  # Remove any old versions of Docker
  apt-get remove docker docker-engine docker.io containerd runc

  # Instal deps
  apt-get update
  apt-get install -y \
      apt-transport-https \
      ca-certificates \
      curl \
      gnupg-agent \
      software-properties-common


  # Add docker CE repo
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  add-apt-repository \
     "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
     $(lsb_release -cs) \
     stable"

  apt-get update
  #Instal docker
  apt-get -y install docker-ce=18.06.1~ce~3-0~ubuntu containerd.io --allow-downgrades

  #Enable docker on startup
  systemctl restart docker.service
  systemctl enable docker.service
  systemctl enable containerd.service

}


function allow_tcp_ssh_forwarding {
  grep -qxF "AllowTcpForwarding yes" /etc/ssh/sshd_config || echo "AllowTcpForwarding yes" >> /etc/ssh/sshd_config
}

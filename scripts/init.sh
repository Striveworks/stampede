#/bin/env bash
set -ex

source install.sh
source pre-install.sh

host1="chariotmaster01.america.striveworks.us"
host2="chariotmaster02.america.striveworks.us"
host3="chariotmaster03.america.striveworks.us"

RKE_USER="rke"

if [ ! -f /root/${RKE_USER}-password.txt]; then
    RKE=$(cat "/root/${RKE_USER}-password.txt")
else
  RKE_PASS=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
  echo $RKE_PASS > "/root/${RKE_USER}-password.txt"
fi


# Setup RKE user
setup_user

# Disable swap
disable_swap

# Kernel param config
kernel_setup

# Install docker
install_docker

# Allow tcp forwarding
allow_tcp_ssh_forwarding


if [ "$CONTROLLER" = "1" ];
  then

  # Create SSH keys for RKE/cluster
  create_ssh_keys $host1 $host2 $host3
  # Get RKE binary
  get_rke_binary
  # Build RKE config
  build_cluster_config $host1 $host2 $host3
  # Init RKE
  cd /opt
  rke_up
fi

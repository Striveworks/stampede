#/bin/env bash
set -ex


if [ $# -eq 1 ]
  then
    INTERFACE_OVERRIDE=$1
fi

mkdir -p /opt/stampede

#Install Go
snap install go --classic

#Build binary
cd /opt/stampede && go build -o /usr/local/bin/stampede .

swapoff -a
sed -i '/ swap / s/^/#/' /etc/fstab

cat <<EOF | tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system


apt-get update -y
apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg -y

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update -y
apt-get install docker-ce docker-ce-cli containerd.io -y

systemctl start docker
systemctl enable docker

apt-get update && sudo apt-get install -y apt-transport-https curl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF | tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
apt-get update
apt-get install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl

#Send multicast traffic through default interface, or specified override
if [ -z "${INTERFACE_OVERRIDE}" ]
  then
    ip route add 224.0.0.0/4 dev $(route | grep '^default' | grep -o '[^ ]*$')
    ADVERTISE_ADDRESS=$(ip addr show dev $(route | grep '^default' | grep -o '[^ ]*$') | grep -m1 inet | egrep -o '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' | head -1)
  else
    ip route add 224.0.0.0/4 dev "${INTERFACE_OVERRIDE}"
    ADVERTISE_ADDRESS=$(ip addr show dev "${INTERFACE_OVERRIDE}" | grep -m1 inet | egrep -o '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' | head -1)
fi


#Setup stampede systemd service
cat << EOF > /lib/systemd/system/stampede.service
[Unit]
Description=stampede

[Service]
Type=simple
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/stampede init --cluster-type kubeadm --advertise-address $ADVERTISE_ADDRESS

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable stampede
systemctl start stampede

#/bin/env bash
set -ex

mkdir -p /opt/stampede

#Install Go
snap install go --classic

#Build binary
cd /opt/stampede && go build -o /usr/local/bin/stampede .

#Install microk8s
snap install microk8s --classic --channel=1.18/stable

#Send multicast traffic through default interface
ip route add 224.0.0.0/4 dev $(route | grep '^default' | grep -o '[^ ]*$')

#Setup stampede systemd service
cat << EOF > /lib/systemd/system/stampede.service
[Unit]
Description=stampede is a microk8s bootstrapping utility to elect a leader and add nodes

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/usr/local/bin/stampede

[Install]
WantedBy=multi-user.target
EOF

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

#Install microk8s
snap install microk8s --classic --channel=1.18/stable

#Send multicast traffic through default interface, or specified override
if [ -z "${INTERFACE_OVERRIDE}" ]
  then
    ip route add 224.0.0.0/4 dev $(route | grep '^default' | grep -o '[^ ]*$')
  else
    ip route add 224.0.0.0/4 dev "${INTERFACE_OVERRIDE}"
fi

#Setup stampede systemd service
cat << EOF > /lib/systemd/system/stampede.service
[Unit]
Description=stampede

[Service]
Type=simple
Restart=on-failure
RestartSec=5s
ExecStart=/usr/local/bin/stampede init

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable stampede
systemctl start stampede

#!/usr/bin/sh
echo "Installing managerd..."
# app directories
MNGRD_LOGS="/var/dynamite/managerd/logs"
MNGRD_APP="/opt/dynamite/managerd/bin"
MNGRD_CONF="/etc/dynamite/managerd"
MNGRD_PID="/var/run/dynamite/managerd"

# create dirs 
echo "Creating application directories..."
mkdir -p $MNGRD_APP
mkdir -p $MNGRD_CONF
mkdir -p $MNGRD_LOGS
mkdir -p $MNGRD_PID

# place files 
echo "Installing application files..."
cp ../cmd/managerd $MNGRD_APP/.
chmod +x $MNGRD_APP/managerd
cp ../internal/conf/config.yml $MNGRD_CONF/.
cp ../etc/managerd.service /etc/systemd/system/.

# set up managerd service 
echo "Enabling managerd service..."
systemctl daemon-reload 
systemctl enable managerd.service


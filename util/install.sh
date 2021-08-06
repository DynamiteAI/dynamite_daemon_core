#!/usr/bin/bash
echo "Installing dynamited..."
# app directories
DYND_LOGS="/var/dynamite/dynamited/logs"
DYND_APP="/opt/dynamite/dynamited/bin"
DYND_CONF="/etc/dynamite/dynamited"
DYND_PID="/var/run/dynamite/dynamited"

# create dirs 
echo "Creating application directories..."
mkdir -p $DYND_APP
mkdir -p $DYND_CONF
mkdir -p $DYND_LOGS
mkdir -p $DYND_PID

# place files 
echo "Installing application files..."
cp ../conf/config.yml $DYND_CONF/.
cp ../systemd/dynamited.service /etc/systemd/system/.

# set up dynamited service 
echo "Enabling dynamited service..."
systemctl daemon-reload 
systemctl enable dynamited.service

echo "Dynamited service installed. Run `systemctl start dynamited` to start."

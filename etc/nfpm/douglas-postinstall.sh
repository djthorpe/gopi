#!/bin/bash

# Enable service
deb-systemd-helper enable douglas.service

# Add user and set groups
id -u gopi &>/dev/null || useradd --system gopi
usermod -a -G spi,gpio,video gopi

# Add directories and permissions
install -o gopi -g gopi -m 775 -d /opt/gopi/var/images

# Start service
deb-systemd-invoke start douglas.service

#!/bin/bash

# Enable service
deb-systemd-helper enable httpserver.service

# Add user and set groups
id -u gopi &>/dev/null || useradd --system gopi

# Add directories and permissions
install -o gopi -g gopi -m 775 -d /opt/gopi/var/htdocs

# Start service
deb-systemd-invoke start httpserver.service

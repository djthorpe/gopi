#!/bin/bash

# Enable service
deb-systemd-helper enable rotel.service

# Add user and set groups
id -u gopi &>/dev/null || useradd --system gopi
usermod -a -G dialout gopi

# Start service
deb-systemd-invoke start rotel.service

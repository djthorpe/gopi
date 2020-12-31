#!/bin/bash

# Enable service
deb-systemd-helper enable googlecast.service

# Add user and set groups
id -u gopi &>/dev/null || useradd --system gopi

# Start service
deb-systemd-invoke start googlecast.service

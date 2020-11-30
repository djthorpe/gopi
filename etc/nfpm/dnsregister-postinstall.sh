#!/bin/bash

# Enable service
systemctl enable dnsregister.service

# Add user
id -u gopi &>/dev/null || useradd --system -G i2c,video gopi

# Start service
systemctl start dnsregister.service

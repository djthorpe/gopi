#!/bin/bash

# Daemon reload
systemctl --system daemon-reload

# Stop service
deb-systemd-invoke stop dnsregister.service

# Purge
deb-systemd-helper purge dnsregister.service
deb-systemd-helper unmask dnsregister.service



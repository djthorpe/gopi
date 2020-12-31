#!/bin/bash

# Daemon reload
systemctl --system daemon-reload

# Stop service
deb-systemd-invoke stop googlecast.service

# Purge
deb-systemd-helper purge googlecast.service
deb-systemd-helper unmask googlecast.service


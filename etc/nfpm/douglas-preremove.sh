#!/bin/bash

# Daemon reload
systemctl --system daemon-reload

# Stop service
deb-systemd-invoke stop douglas.service

# Purge
deb-systemd-helper purge douglas.service
deb-systemd-helper unmask douglas.service



#!/bin/bash

# Daemon reload
systemctl --system daemon-reload

# Stop service
deb-systemd-invoke stop argonone.service

# Purge
deb-systemd-helper purge argonone.service
deb-systemd-helper unmask argonone.service


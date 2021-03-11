#!/bin/bash

# Daemon reload
systemctl --system daemon-reload

# Stop service
deb-systemd-invoke stop rotel.service

# Purge
deb-systemd-helper purge rotel.service
deb-systemd-helper unmask rotel.service



#!/bin/bash

# Daemon reload
systemctl --system daemon-reload

# Stop service
deb-systemd-invoke stop httpserver.service

# Purge
deb-systemd-helper purge httpserver.service
deb-systemd-helper unmask httpserver.service





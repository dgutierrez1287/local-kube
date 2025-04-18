#!/usr/bin/env bash

echo "[Resolve]" > /etc/systemd/resolved.conf
echo "DNS=8.8.8.8 8.8.4.4" >> /etc/systemd/resolved.conf

sudo systemctl restart systemd-resolved

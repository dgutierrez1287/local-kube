#!/usr/bin/env bash

DISK="/dev/sda"

FREE_SPACE=$(sudo parted $DISK print free | awk '/Free Space/ {print $3}')

if [[ -z "$FREE_SPACE" ]]; then
  echo "Drive was not expanded, no free space to expand to"
  exit 0
else 
  sudo sgdisk -e $DISK
  sudo parted -s -a opt $DISK "resizepart 3 100%"
  pvresize $DISK
  sudo lvextend -l +100%FREE /dev/ubuntu-vg/ubuntu-lv
  sudo resize2fs /dev/ubuntu-vg/ubuntu-lv
fi

exit 0

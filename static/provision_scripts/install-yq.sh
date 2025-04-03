#!/usr/bin/env bash

echo "Installing yq ppa repo"
sudo add-apt-repository -y ppa:rmescandon/yq

sudo apt update

echo "installing yq"
sudo apt install -y yq


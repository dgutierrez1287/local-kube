#!/usr/bin/env bash

# Args

## Setup ##
echo "creating directory for dynamic ansible variables"
sudo mkdir /etc/ansible/vars/dynamic
sudo rm -f /etc/ansible/vars/dynamic/vars-dynamic.yml
sudo touch /etc/ansible/vars/dynamic/vars-dynamic.yml
sudo chmod 777 /etc/ansible/vars/dynamic/vars-dynamic.yml

## Generate Dynamic Vars off of settings ##
echo "generating dynamic variables"
sudo python3 /vagrant/scripts/single_node/generate_dynamic_vars.py

## Run Ansible ##
echo "Running ansible"
/usr/local/bin/ansible-playbook /etc/ansible/playbook/playbook.yml

## Kube config stuff ##
echo "Setting kube config for vagrant user"
mkdir /home/vagrant/.kube
python3 /vagrant/scripts/correct_kubeconfig.py "single-node"
chmod 777 /home/vagrant/.kube/config

echo "correcting perms to kube config to shared folder"
chmod 777 /vagrant/kubeconfig/config.yaml

#!/usr/bin/env bash

# Args

## Setup ##
echo "creating directory for dynamic ansible variables"
sudo mkdir /etc/ansible/vars/dynamic

sudo rm -f /etc/ansible/vars/dynamic/vars-dynamic-control.yml
sudo rm -f /etc/ansible/vars/dynamic/vars-dynamic-worker.yml

sudo touch /etc/ansible/vars/dynamic/vars-dynamic-control.yml
sudo touch /etc/ansible/vars/dynamic/vars-dynamic-worker.yml

sudo chmod 777 /etc/ansible/vars/dynamic/*

## Generate Dynamic vars off of settings ##
echo "Generating Dynamic Variable files"
sudo python3 /vagrant/scripts/multi_node/generate_dynamic_vars.py "control"
sudo python3 /vagrant/scripts/multi_node/generate_dynamic_vars.py "worker"

## Run Ansible ##
echo "Running ansible on the lead node"
/usr/local/bin/ansible-playbook /etc/ansible/playbook/lead-node-playbook.yml

echo "sleeping to let k3s start fully before provisioning other nodes"
sleep 30

echo "running ansible on the control nodes"
/usr/local/bin/ansible-playbook /etc/ansible/playbook/control-node-playbook.yml

echo "running ansible on the worker nodes"
/usr/local/bin/ansible-playbook /etc/ansible/playbook/worker-playbook.yml

## KubeConfig Stuff ##
echo "Setting kube config for the vagrant user on the lead node"
mkdir /home/vagrant/.kube
python3 /vagrant/scripts/correct_kubeconfig.py "multi-node"
chmod 777 /home/vagrant/.kube/config

echo "Copying kube config to shared folder"
chmod 777 /vagrant/kubeconfig/config.yaml

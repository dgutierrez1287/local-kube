#!/usr/bin/env bash
exec > /vagrant/logs/provision.txt 2>&1

# Args
output_type=$1

## copy ansible hosts file ##
echo "Copying ansible hosts file"
cp /vagrant/ansible-resources/hosts /etc/ansible/hosts
chmod 777 /etc/ansible/hosts

## Run Ansible ##
if [[ "${output_type}" == "debug" ]]; then
  echo "Running ansible in debug mode"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/playbook.yml -vvvv
else
  echo "Running ansible"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/playbook.yml
fi 

## Copy kubeconfig ##
cp /etc/rancher/k3s/k3s.yaml /vagrant/kubeconfig/k3s.yaml
chmod 777 /vagrant/kubeconfig/k3s.yaml

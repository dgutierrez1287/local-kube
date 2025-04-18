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
  echo "Running ansible in debug mode on the lead node"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/lead-node-playbook.yml -vvvv
else 
  echo "Running ansible on the lead node"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/lead-node-playbook.yml
fi

echo "sleeping to let k3s start fully before provisioning other nodes"
sleep 30

if [[ "${output_type}" == "debug" ]]; then
  echo "Running ansible in debug mode on the control nodes"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/control-node-playbook.yml -vvvv
else
  echo "Running ansible on the control nodes"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/control-node-playbook.yml
fi

if [[ "${output_type}" == "debug" ]]; then
  echo "Running ansible in debug mode on the worker nodes"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/worker-playbook.yml -vvvv
else
  echo "Running ansible on the worker nodes"
  /usr/local/bin/ansible-playbook /etc/ansible/playbook/worker-playbook.yml
fi

## Copy Kubeconfig ##
cp /etc/rancher/k3s/k3s.yaml /vagrant/kubeconfig/k3s.yaml
chmod 777 /vagrant/kubeconfig/k3s.yaml

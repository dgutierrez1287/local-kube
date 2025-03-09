import yaml 
import ast
import sys

def get_default_kubeconfig():
    with open('/etc/rancher/k3s/k3s.yaml', 'r') as stream:
        try:
            data = yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print(f"yaml error getting the default kubeconfig {e}")
            exit(123)

    return data

def get_cluster_settings():
    with open('/vagrant/cluster/settings.yaml', 'r') as stream:
        try:
            data = yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print(f"yaml error getting the cluster settings {e}")
            exit(123)

    return data

def get_ansible_user_vars():
    with open('/etc/ansible/vars/user/ansible-vars.yml', 'r') as stream:
        try:
            data = yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print(f"yaml error getting the ansible user vars {e}")
            exit(123)
    
    return data

def is_kubevip_enabled(ansible_user_vars):
    if ansible_user_vars is not None:
        if 'k3s_enable_kubevip' in ansible_user_vars.keys():
            if ansible_user_vars['k3s_enable_kubevip']:
                return True

    return False

def single_node_kube_url(cluster_settings, kubevip_enabled):

    if kubevip_enabled:
        vip = cluster_settings['cluster-vip']
        return f"https://{vip}:6443"
    else:
        vm_ip = cluster_settings['machine_settings']['ip']
        return f"https://{vm_ip}:6443"

def multi_node_kube_url(cluster_settings, kubevip_enabled):

    if kubevip_enabled:
        vip = cluster_settings['cluster-vip'] 
        return f"https://{vip}:6443"
    else:
        lead_vm_ip = cluster_settings['lead-control-node'][0]['ip']
        return f"https://{lead_vm_ip}:6443"

def update_write_out_new_kubeconfig(kubeconfig, kube_ip):

    kubeconfig['clusters'][0]['cluster']['server'] = kube_ip

    with open(f'/home/vagrant/.kube/config', 'w') as outfile:
        yaml.dump(kubeconfig, outfile, default_flow_style=False)

    with open(f'/vagrant/kubeconfig/config.yaml', 'w') as outfile2:
        yaml.dump(kubeconfig, outfile2, default_flow_style=False)

if __name__=="__main__":

    cluster_type = sys.argv[1]
    
    kubeconfig = get_default_kubeconfig()
    cluster_settings = get_cluster_settings()
    ansible_user_vars = get_ansible_user_vars()

    kubevip_enabled = is_kubevip_enabled(ansible_user_vars)

    if cluster_type == "single-node":
        kube_ip = single_node_kube_url(cluster_settings, kubevip_enabled)
        update_write_out_new_kubeconfig(kubeconfig, kube_ip)
    elif cluster_type == "multi-node":
        kube_ip = multi_node_kube_url(cluster_settings, kubevip_enabled)
        update_write_out_new_kubeconfig(kubeconfig, kube_ip)

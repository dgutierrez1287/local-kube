package settings

import (
	"fmt"

	"github.com/dgutierrez1287/local-kube/logger"
)

/*
  Cluster - Settings for a cluster
*/
type Cluster struct {
  KubeConfigName string             `json:"kubeconfigName,omitempty"`     // The name for the cluster in kubeconfig (if empty cluster name is used)
  Vip string                        `json:"vip,omitempty"`                // The vip for the kubernetes cluster
  ClusterType string                `json:"clusterType,omitempty"`        // (single or ha) the type of cluster
  ProviderName string               `json:"providerName,omitempty"`       // The name of the provider that the cluster uses
  Leaders []Machine                 `json:"leaders"`                      // a list of leader machines
  Workers []Machine                 `json:"workers"`                      // a list of worker machines
  ClusterFeatures *ClusterFeatures  `json:"clusterFeatures,omitempty"`    // feature configuration only used if autoConfigre is true
}

/*
  Machine - Settings for a given machine in the 
  cluster
*/
type Machine struct {
  Name string         `json:"name" yaml:"name,omitempty"`           // The name of the machine
  IpAddress string    `json:"ipAddress" yaml:"ip,omitempty"`         // The IP address for the machine
  Memory int          `json:"memory" yaml:"memory,omitempty"`       // The memory for the machine
  Cpu int             `json:"cpus" yaml:"cpus,omitempty"`             // The CPU setting for the machine
  DiskSize string     `json:"diskSize" yaml:"disk_size,omitempty"`  // The size of the primary disk
}

/*
Gets a list of all but the first control node and returns
a list of just the node names
*/
func (cluster Cluster) GetSecondaryControlNodeNames() []string {

  names := []string{}
  secondaryControlNodes := cluster.Leaders[1:]

  for _, node := range secondaryControlNodes {
    names = append(names, node.Name)
  }
  return names
}

/*
Gets a list of all the machine names in the cluster
*/
func (cluster Cluster) GetMachineNameList() []string {
  machineNameList := []string{}

  for _, machine := range(cluster.Leaders) {
    machineNameList = append(machineNameList, machine.Name)
  }

  for _, machine := range(cluster.Workers) {
    machineNameList = append(machineNameList, machine.Name)
  }
  
  return machineNameList
}

/*
Gets the server url for the cluster based on if 
kubevip is enabled and if the cluster is ha or single
*/
func (cluster Cluster) GetServerUrl() string {

  if cluster.ClusterFeatures.KubeVipEnable {
    logger.LogDebug("KubeVip is enabled returning the server url using kubevip")
    return fmt.Sprintf("https://%s:6443", cluster.Vip)
  }

  logger.LogDebug("Returning the lead node ip since kubevip is not enabled")
  leaderIp := cluster.Leaders[0].IpAddress
  return fmt.Sprintf("https://%s:6443", leaderIp)
}

/*
Gets the vagrant machine name for the node that is 
setup with ansible for provisioning, in multi-node this
is the first control node
*/
func (cluster Cluster) GetAnsibleNodeVagrantName() string {
  if cluster.ClusterType == "single" {
    return "default"
  } else {
    return cluster.Leaders[0].Name
  }
}  

/*
Gets a list of just the worker node names 
*/
func (cluster Cluster) GetWorkerNodeNames() []string {

  names := []string{}

  for _, node := range cluster.Workers {
    names = append(names, node.Name)
  }
  return names
}

/*
Gets a list of control node IPs 
*/
func (cluster Cluster) GetControlNodeIps() []string {
  
  ips := []string{}

  for _, node := range cluster.Leaders {
    ips = append(ips, node.IpAddress)
  }
  return ips
}

/*
check if the cluster is ha or not
*/
func (cluster Cluster) IsHA() bool {
  if cluster.ClusterType == "ha" {
    return true
  } else {
    return false
  }
}



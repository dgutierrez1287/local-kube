package settings

/*
  Cluster - Settings for a cluster
*/
type Cluster struct {
  Vip string                        `json:"vip,omitempty"`                // The vip for the kubernetes cluster
  ClusterType string                `json:"clusterType,omitempty"`        // (single or ha) the type of cluster
  ProviderName string               `json:"providerName,omitempty"`       // The name of the provider that the cluster uses
  Leaders []Machine                 `json:"leaders"`                      // a list of leader machines
  Workers []Machine                 `json:"workers"`                      // a list of worker machines
  ClusterFeatures *ClusterFeatures   `json:"clusterFeatures,omitempty"`    // feature configuration only used if autoConfigre is true
}

/*
  Machine - Settings for a given machine in the 
  cluster
*/
type Machine struct {
  Name string         `json:"name"`       // The name of the machine
  IpAddress string    `json:"ipAdress"`   // The IP address for the machine
  Memory int          `json:"memory"`     // The memory for the machine
  Cpu int             `json:"cpu"`        // The CPU setting for the machine
  DiskSize string     `json:"diskSize"`   // The size of the primary disk
}

// Gets a list of all but the first control node and returns
// a list of just the node names
func (cluster Cluster) GetSecondaryControlNodeNames() []string {

  names := []string{}
  secondaryControlNodes := cluster.Leaders[1:]

  for _, node := range secondaryControlNodes {
    names = append(names, node.Name)
  }
  
  return names
}

// Gets a list of just the worker node names 
func (cluster Cluster) GetWorkerNodeNames() []string {

  names := []string{}

  for _, node := range cluster.Workers {
    names = append(names, node.Name)
  }

  return names
}

// Gets a list of control node IPs 
func (cluster Cluster) GetControlNodeIps() []string {
  
  ips := []string{}

  for _, node := range cluster.Leaders {
    ips = append(ips, node.IpAddress)
  }

  return ips
}

// check if the cluster is ha or not
func (cluster Cluster) IsHA() bool {
  if cluster.ClusterType == "ha" {
    return true
  } else {
    return false
  }
}



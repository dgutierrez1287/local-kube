package settings

/*
  Cluster - Settings for a cluster
*/
type Cluster struct {
  Vip string          `json:"vip"`            // The vip for the kubernetes cluster
  ClusterType string  `json:"clusterType"`    // (single or ha) the type of cluster
  Leaders []Machine   `json:"leaders"`        // a list of leader machines
  Workers []Machine   `json:"workers"`        // a list of worker machines
}

/*
  Machine - Settings for a given machine in the 
  cluster
*/
type Machine struct {
  Name string         `json:"name"`       // The name of the machine
  IpAddress string    `json:"ipAdress"`   // The IP address for the machine
  Memory int          `json:"memory"`     // The memory for the machine
  Cpu string          `json:"cpu"`        // The CPU setting for the machine
}


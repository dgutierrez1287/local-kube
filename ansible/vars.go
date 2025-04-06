package ansible

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"gopkg.in/yaml.v3"
)

/*
  Structs for mapping items to yaml for Marshaling
  These structs cover different areas and will be
  merged to make the variables file
*/

// General Vars
type GeneralVars struct {
  KubeVersion string                      `yaml:"kube_kube_version,omitempty"`
  HaCluster bool                          `yaml:"kube_ha_cluster,omitempty"`
  IsLeadNode bool                         `yaml:"kube_is_lead_control_plane_node,omitempty"`
  NodePurpose string                      `yaml:"kube_node_purpose,omitempty"`
  TlsSanAddresses []string                `yaml:"kube_tls_san_addresses,omitempty"`
  ControlPlanList map[string]ControlNode  `yaml:"kube_control_plane_list,omitempty"`
  CniController string                    `yaml:"kube_cni_controller,omitempty"`
  IngressController string                `yaml:"kube_ingress_controller,omitempty"`
  StorageController string                `yaml:"kube_storage_controller,omitempty"`
  DisableDefaultMetrics bool              `yaml:"kube_disable_default_metrics,omitempty"`
  UseIps bool                             `yaml:"kube_use_ips,omitempty"`
} 

// control node type
type ControlNode struct {
  primary bool                            `yaml:"primary"`
  ip string                               `yaml:"ip"`
}

// KubeVip
type KubeVipVars struct {
  Vip string                        `yaml:"kube_kubevip_ip,omitempty"`
  Version string                    `yaml:"kube_kubevip_version,omitempty"`
  Enable bool                       `yaml:"kube_enable_kubevip,omitempty"`
}

// Cilium
type CiliumVars struct {
  Version string                    `yaml:"kube_cilium_version,omitempty"`
  CliVersion string                 `yaml:"kube_cilium_cli_version,omitempty"`
  Install bool                      `yaml:"kube_cilium_install_cilium,omitempty"`
  InstallHubble bool                `yaml:"kube_cilium_install_hubble,omitempty"`
}

// Calico
type CalicoVars struct {
  Version string                    `yaml:"kube_calico_version,omitempty"`
  Install bool                      `yaml:"kube_calico_install_calico,omitempty"`
}

// Longhorn
type LonghornVars struct {
  Version string                    `yaml:"kube_longhorn_version,omitempty"`
  Install bool                      `yaml:"kube_install_longhorn,omitempty"`
}

// Merged variables that will hold all variables
// to be Marshaled to a file
type MergedVars struct {
  GeneralVars     `yaml:",inline,omitempty"`
  KubeVipVars     `yaml:",inline,omitempty"`
  CiliumVars      `yaml:",inline,omitempty"`
  CalicoVars      `yaml:",inline,omitempty"`
  LonghornVars    `yaml:",inline,omitempty"`
}

func GenerateVarsFile(appDir string, clusterName string, clusterType string, 
  nodeType string, appSettings settings.Settings) error {

  varsFilePath := filepath.Join(appDir, clusterName, "ansible", "variables")
  
  var fileName string
  var haCluster bool
  // use ips should always be true for local clusters
  useIps := true

  // pull out the features out of settings to make function 
  // calls easier
  features := *appSettings.Clusters[clusterName].ClusterFeatures

  // set up settings used to build
  // the variables
  if clusterType == "ha" {
    logger.Logger.Debug("setting some variables for ha cluster")

    fileName = fmt.Sprintf("vars-%s.yml", nodeType)
    haCluster = true
  } else {
    logger.Logger.Debug("Setting some variables for single node cluster")

    fileName = "vars.yml"
    haCluster = false
  }

  logger.Logger.Debug("vars file name being rendered is", "filename", fileName)

  logger.Logger.Debug("Getting general variables")
  generalVars := getGeneralVars(appSettings, haCluster, useIps, nodeType, clusterName)
  logger.Logger.Debug("Getting kubevip variables")
  kubevipVars := getKubeVipVars(features, appSettings.Clusters[clusterName].Vip)
  logger.Logger.Debug("Getting cilium variables")
  ciliumVars := getCiliumVars(features)
  logger.Logger.Debug("Getting calico variables")
  calicoVars := getCalicoVars(features)
  logger.Logger.Debug("Getting longhorn variables")
  longhornVars := getLonghornVars(features)

  logger.Logger.Debug("Merging all variables for marshaling to file")
  vars := MergedVars{
    GeneralVars: generalVars,
    KubeVipVars: kubevipVars,
    CiliumVars: ciliumVars,
    CalicoVars: calicoVars,
    LonghornVars: longhornVars,
  }

  yamlData, err := yaml.Marshal(&vars)
  if err != nil {
    logger.Logger.Error("Error marshaling variable data to yaml")
    return err
  }

  path := filepath.Join(varsFilePath, fileName)
  file, err := os.Create(path)

  if err != nil {
    logger.Logger.Error("Error creating vars file", "path", path)
    return err
  }

  defer file.Close()

  _, err = file.WriteString("---\n")
  if err != nil {
    logger.Logger.Error("Error writing the yaml doc start marker")
    return err
  }

  _, err = file.Write(yamlData)
  if err != nil {
    logger.Logger.Error("Error writing out vars content to file")
    return err
  }
  return nil 
}

// Gets all the general settings set from cluster features
func getGeneralVars(appSettings settings.Settings, haCluster bool, 
  useIps bool, nodeType string, clusterName string) GeneralVars {
  var vars GeneralVars
  var isLeadNode bool
  var nodePurpose string

  features := appSettings.Clusters[clusterName].ClusterFeatures
  cpNodes := appSettings.Clusters[clusterName].Leaders
  vip := appSettings.Clusters[clusterName].Vip

  switch nodeType {
  case "lead":
    isLeadNode = true
    nodePurpose = "lead"
  case "control":
    isLeadNode = false
    nodePurpose = "control-plane"
  case "worker":
    isLeadNode = false
    nodePurpose = "worker"
  default:
    isLeadNode = false
    nodePurpose = ""
  }

  logger.Logger.Debug("Generating general variables")
  vars.KubeVersion = features.KubeVersion
  vars.HaCluster = haCluster
  vars.IsLeadNode = isLeadNode
  vars.NodePurpose = nodePurpose
  vars.TlsSanAddresses = GetTlsSanList(cpNodes, features.KubeVipEnable, vip)
  vars.ControlPlanList = GetControlPlaneList(cpNodes, appSettings.Clusters[clusterName].ClusterType)
  vars.CniController = features.CniController
  vars.IngressController = features.IngressController
  vars.StorageController = features.StorageController
  vars.DisableDefaultMetrics = features.DisableDefaultMetrics
  vars.UseIps = useIps

  return vars
}

// Gets KubeVip vars if kubevip is not set it will return empty
func getKubeVipVars(features settings.ClusterFeatures, vip string) KubeVipVars {
  var vars KubeVipVars 

  if !features.KubeVipEnable {
    logger.Logger.Debug("Kubevip is not enabled, not setting kubevip settings")
    return vars
  }

  logger.Logger.Debug("Kubevip is enabled, setting kubevip settings")
  vars.Vip = vip
  vars.Version = features.KubeVipVersion
  vars.Enable = true

  return vars
}

// Gets Cilium vars if not enabled it will return empty
func getCiliumVars(features settings.ClusterFeatures) CiliumVars {
  var vars CiliumVars
  
  if features.CniController != "cilium" {
    logger.Logger.Debug("CniController is not cilium, not setting cilium settings")
    return vars
  }

  logger.Logger.Debug("CniController is cilium, setting the cilum settings")
  vars.Version = features.CniControllerVersion
  vars.CliVersion = features.CiliumCliVersion
  vars.Install = features.ManagedCniController
  vars.InstallHubble = true

  return vars
}

// Gets calico vars if not enabled it will return empty
func getCalicoVars(features settings.ClusterFeatures) CalicoVars {
  var vars CalicoVars

  if features.CniController != "calico" {
    logger.Logger.Debug("CniController is not Calico, not setting calico settings")
    return vars
  }

  logger.Logger.Debug("CniController is calico, setting the calico settings")
  vars.Version = features.CniControllerVersion
  vars.Install = features.ManagedCniController

  return vars
}

// Gets longhorn vars if not enabled it will return empty
func getLonghornVars(features settings.ClusterFeatures) LonghornVars {
  var vars LonghornVars

  if features.StorageController != "longhorn" {
    logger.Logger.Debug("StorageController is not longhorn, not setting longhorn settings")
    return vars
  }

  logger.Logger.Debug("StorageController is Longhorn, setting the longhorn settings")
  vars.Version = features.StorageControllerVersion
  vars.Install = features.ManagedStorageController

  return vars
}
 
/*
  Gets a list of all the control node Ip addresses and the vip if
  kubevip is enabled, this is needed to add the TLS san setting in k3s
  since it manages its own cert
*/
func GetTlsSanList(leadNodes []settings.Machine, kubeVipEnabled bool, kubeVipIp string) []string {
  var tlsSanIps = []string{}
  
  logger.Logger.Debug("Getting list of ips of control nodes")
  for _, node := range leadNodes {
    logger.Logger.Debug("Processing node", "node", node)
    logger.Logger.Debug("Adding node ip to san list", "ip", node.IpAddress)
    tlsSanIps = append(tlsSanIps, node.IpAddress)
  }

  if kubeVipEnabled {
    logger.Logger.Debug("KubeVip is enabled adding to tls san ip list")
    tlsSanIps = append(tlsSanIps, kubeVipIp)
  }

  logger.Logger.Debug("Tls san ip list", "ips", tlsSanIps)
  return tlsSanIps
}

/*
  Gets a list of control plane nodes
*/
func GetControlPlaneList(leadNodes []settings.Machine, clusterType string) map[string]ControlNode {
  var cpList = make(map[string]ControlNode)

  if clusterType == "single" {
    return cpList
  }

  for index, node := range leadNodes {
    if index == 0 {
      logger.Logger.Debug("Adding lead node to control plane list", "name", node.Name)
      cpList[node.Name] = ControlNode{
        primary: true,
        ip: node.IpAddress,
      }
    }
    logger.Logger.Debug("Adding control plane node", "name", node.Name)
    cpList[node.Name] = ControlNode{
      primary: false,
      ip: node.IpAddress,
    }
  }
  return cpList
}





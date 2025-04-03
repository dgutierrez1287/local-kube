package cluster

import (
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"gopkg.in/yaml.v3"
)

/*
  ScriptSettings - Settings that will be written out
  in a yaml file that the machine(s) will use in local
  scripts
*/
type ScriptSettings struct {
  ClusterName string                  `yaml:"cluster-name,omitempty"`
  ClusterVip string                   `yaml:"cluster-vip,omitempty"`
  LeadNode []settings.Machine         `yaml:"lead-control-node,omitempty"`
  ControlNodes []settings.Machine     `yaml:"control-nodes,omitempty"`
  WorkerNodes []settings.Machine      `yaml:"workers,omitempty"`
  MachineSettings settings.Machine    `yaml:"machine_settings,omitempty"`
}
  
/*
GernerateScriptSettings - This generates a settings yaml file that will be used
by scripts in the vms to do settings that can only be done on the machine
*/
func GenerateScriptSettings(appDir string, clusterName string, appSettings settings.Settings) error {
  settingsFile := filepath.Join(appDir, clusterName, "settings", "settings.yaml")

  settings := ScriptSettings{}
  clusterSettings := appSettings.Clusters[clusterName]

  settings.ClusterName = clusterName

  if clusterSettings.ClusterFeatures.KubeVipEnable {
    logger.Logger.Debug("Kube vip enabled setting vip in settings")
    settings.ClusterVip = clusterSettings.Vip
  }

  if clusterSettings.ClusterType == "ha" {
    logger.Logger.Debug("Cluster is an ha cluster, setting machines for settings file")

    for index, node := range clusterSettings.Leaders {

      if index == 0 {
        logger.Logger.Debug("Adding lead node to settings", "name", node.Name)
        settings.LeadNode = append(settings.LeadNode, node)
      } else {
        logger.Logger.Debug("Adding control node to settings", "name", node.Name)
        settings.ControlNodes = append(settings.ControlNodes, node)
      }
    }

    for _, node := range clusterSettings.Workers {
      logger.Logger.Debug("Adding worker node to settings", "name", node.Name)
      settings.WorkerNodes = append(settings.WorkerNodes, node)
    }
    
  } else {
    logger.Logger.Debug("Cluster is a single node cluster adding the single machine to settings")
    settings.MachineSettings = clusterSettings.Leaders[0]
  }

  yamlData, err := yaml.Marshal(&settings)
  if err != nil {
    logger.Logger.Error("Error marshaling script settings to yaml")
    return err
  }

  err = os.WriteFile(settingsFile, yamlData, 0755)
  if err != nil {
    logger.Logger.Error("Error writing script settings data to file")
    return err
  }
  return nil
}

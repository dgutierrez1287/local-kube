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
  LeadNode []settings.Machine         `yaml:"lead-control-node"`
  ControlNodes []settings.Machine     `yaml:"control-nodes"`
  WorkerNodes []settings.Machine      `yaml:"workers"`
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
    logger.LogDebug("Kube vip enabled setting vip in settings")
    settings.ClusterVip = clusterSettings.Vip
  }

  if clusterSettings.ClusterType == "ha" {
    logger.LogDebug("Cluster is an ha cluster, setting machines for settings file")

    for index, node := range clusterSettings.Leaders {

      if index == 0 {
        logger.LogDebug("Adding lead node to settings", "name", node.Name)
        settings.LeadNode = append(settings.LeadNode, node)
      } else {
        logger.LogDebug("Adding control node to settings", "name", node.Name)
        settings.ControlNodes = append(settings.ControlNodes, node)
      }
    }

    for _, node := range clusterSettings.Workers {
      logger.LogDebug("Adding worker node to settings", "name", node.Name)
      settings.WorkerNodes = append(settings.WorkerNodes, node)
    }
    
  } else {
    logger.LogDebug("Cluster is a single node cluster adding the single machine to settings")
    settings.MachineSettings = clusterSettings.Leaders[0]
  }

  yamlData, err := yaml.Marshal(&settings)
  if err != nil {
    logger.LogError("Error marshaling script settings to yaml")
    return err
  }

  file, err := os.Create(settingsFile)
  if err != nil {
    logger.LogError("Error creating the script settings file")
    return err
  }

  defer file.Close()

  _, err = file.WriteString("---\n")
  if err != nil {
    logger.LogError("Error writing the yaml doc start marker")
    return err
  }

  _, err = file.Write(yamlData)
  if err != nil {
    logger.LogError("Error writing script settings content to file")
    return err
  }
  return nil
}

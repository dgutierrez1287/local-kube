package cluster

import (

	"github.com/dgutierrez1287/local-kube/ansible"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
)

func GenerateAnsibleResources(appDir string, clusterName string, appSettings settings.Settings) {

  // empty node name lists
  var secondaryControlNodeNames []string
  var workerNodeNames []string

  logger.Logger.Info("Generating Ansible resources")
  logger.Logger.Debug("Generating ansible hosts file")
  if appSettings.Clusters[clusterName].ClusterType == "single" {

    logger.Logger.Debug("Single Node creating empty node name arrays")

  } else {

    logger.Logger.Debug("Multi Node getting node name lists from settings")
    secondaryControlNodeNames = appSettings.Clusters[clusterName].GetSecondaryControlNodeNames()
    workerNodeNames = appSettings.Clusters[clusterName].GetWorkerNodeNames()

    logger.Logger.Debug("secondary control nodes", "names", secondaryControlNodeNames)
    logger.Logger.Debug("work control nodes", "names", workerNodeNames)
  }

  // Generate ansible hosts file
  err := ansible.GenerateAnsibleHostsFile(appDir, clusterName,
    appSettings.Clusters[clusterName].ClusterType,
    secondaryControlNodeNames,
    workerNodeNames)

  if err != nil {
    logger.Logger.Error("Error creating ansible hosts file")
    FailureCleanup(appDir, clusterName)
  }

  logger.Logger.Debug("Generating bootstrap.sh script")
  err = ansible.RenderBootstrapScript(appDir, clusterName,
    appSettings.ProvisionSettings.AnsibleVersion,
    appSettings.ProvisionSettings.AnsibleCollections)

  if err != nil {
    logger.Logger.Error("Error Rendering bootstrap script")
    FailureCleanup(appDir, clusterName)
  }
}

func renderVagrantFile(clusterName string, appSettings settings.Settings) {
  
}


package cluster

import (

	"github.com/dgutierrez1287/local-kube/ansible"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
)

/*
  Generates the ansible bootstrap script for the lead node, this 
  will have the desired version of ansible on it. Also it will
  generate the ansible hosts file for the cluster (either ha or single node)
*/
func GenerateAnsibleResources(appDir string, clusterName string, appSettings settings.Settings) error {

  // empty node name lists
  var secondaryControlNodeNames []string
  var workerNodeNames []string

  logger.Logger.Debug("Generating ansible hosts file")
  if appSettings.Clusters[clusterName].ClusterType == "single" {

    logger.Logger.Debug("Single Node creating empty node name arrays")
    secondaryControlNodeNames = []string{}
    workerNodeNames = []string{}
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
    return err
  }

  logger.Logger.Debug("Generating bootstrap.sh script")
  err = ansible.RenderBootstrapScript(appDir, clusterName,
    appSettings.ProvisionSettings.AnsibleVersion)

  if err != nil {
    logger.Logger.Error("Error Rendering bootstrap script")
    return err
  }
  return nil
}

/*
  This will generate all needed playbooks for the desired cluster. There
  will be 3 playbooks for an ha cluster and 1 playbook for a single node
  cluster
*/
func GenerateAnsiblePlaybooks(appDir string, clusterName string, appSettings settings.Settings) error {

  clusterType := appSettings.Clusters[clusterName].ClusterType

  if clusterType == "ha" {
    logger.Logger.Debug("Generating playbooks for ha cluster")

    logger.Logger.Debug("Rendering playbook for lead node")
    err := ansible.RenderPlaybook(appDir, clusterName, "localhost",
      clusterType, "lead")

    if err != nil {
      logger.Logger.Error("Error rendering lead node playbook")
      return err 
    }

    logger.Logger.Debug("Rendering playbook for control nodes")
    err = ansible.RenderPlaybook(appDir, clusterName, "control-nodes",
      clusterType, "control")

    if err != nil {
      logger.Logger.Error("Error rendering playbook for control nodes")
      return err
    }

    logger.Logger.Debug("Rendering playbook for worker nodes")
    err = ansible.RenderPlaybook(appDir, clusterName, "worker-nodes",
      clusterType, "worker")

    if err != nil {
      logger.Logger.Error("Error rendering playbook for worker nodes")
      return err
    }
  } else {
    logger.Logger.Debug("Generating playbooks for single node cluster")
    err := ansible.RenderPlaybook(appDir, clusterName, "localhost",
      clusterType, "")

    if err != nil {
      logger.Logger.Error("Error rendering playbook for single node cluster")
      return err
    }
  }
  return nil
}

/*
  This will generate all the variables files for the desired cluster. There 
  will be 3 variables files for an ha cluster and 1 for a single node cluster
*/
func GenerateAnsibleVariables(appDir string, clusterName string, appSettings settings.Settings) error {

  clusterType := appSettings.Clusters[clusterName].ClusterType

  if clusterType == "ha" {
    logger.Logger.Debug("Generating variables file for ha cluster")

    logger.Logger.Debug("Renderinng vars for lead node")
    err := ansible.GenerateVarsFile(appDir, clusterName, clusterType, "lead", appSettings)
    
    if err != nil {
      logger.Logger.Error("Error rendering vars file for lead node")
      return err
    }

    logger.Logger.Debug("Error Rendering vars for control nodes")
    err = ansible.GenerateVarsFile(appDir, clusterName, clusterType, "control", appSettings)

    if err != nil {
      logger.Logger.Error("Error rendering vars for control nodes")
      return err
    }

    logger.Logger.Debug("Rendering vars for worker nodes") 
    err = ansible.GenerateVarsFile(appDir, clusterName, clusterType, "worker", appSettings)

    if err != nil {
      logger.Logger.Error("Error rendering vars for worker nodes")
      return err
    }
  } else {
    logger.Logger.Debug("Generating variables for single node cluster")
    err := ansible.GenerateVarsFile(appDir, clusterName, clusterType, "", appSettings)

    if err != nil {
      logger.Logger.Error("Error rendering vars for single node cluster")
      return err
    }
  }
  return nil
}

func renderVagrantFile(clusterName string, appSettings settings.Settings) {
  
}


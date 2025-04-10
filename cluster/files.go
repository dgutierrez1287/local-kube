package cluster

import (
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/ansible"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/static"
	"github.com/dgutierrez1287/local-kube/template"
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

  logger.LogDebug("Generating ansible hosts file")
  if appSettings.Clusters[clusterName].ClusterType == "single" {

    logger.LogDebug("Single Node creating empty node name arrays")
    secondaryControlNodeNames = []string{}
    workerNodeNames = []string{}
  } else {

    logger.LogDebug("Multi Node getting node name lists from settings")
    secondaryControlNodeNames = appSettings.Clusters[clusterName].GetSecondaryControlNodeNames()
    workerNodeNames = appSettings.Clusters[clusterName].GetWorkerNodeNames()

    logger.LogDebug("secondary control nodes", "names", secondaryControlNodeNames)
    logger.LogDebug("work control nodes", "names", workerNodeNames)
  }

  // Generate ansible hosts file
  err := ansible.GenerateAnsibleHostsFile(appDir, clusterName,
    appSettings.Clusters[clusterName].ClusterType,
    secondaryControlNodeNames,
    workerNodeNames)

  if err != nil {
    logger.LogError("Error creating ansible hosts file")
    return err
  }

  logger.LogDebug("Generating bootstrap.sh script")
  err = ansible.RenderBootstrapScript(appDir, clusterName,
    appSettings.ProvisionSettings.AnsibleVersion)

  if err != nil {
    logger.LogError("Error Rendering bootstrap script")
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
    logger.LogDebug("Generating playbooks for ha cluster")

    logger.LogDebug("Rendering playbook for lead node")
    err := ansible.RenderPlaybook(appDir, clusterName, "localhost",
      clusterType, "lead")

    if err != nil {
      logger.LogError("Error rendering lead node playbook")
      return err 
    }

    logger.LogDebug("Rendering playbook for control nodes")
    err = ansible.RenderPlaybook(appDir, clusterName, "control-nodes",
      clusterType, "control")

    if err != nil {
      logger.LogError("Error rendering playbook for control nodes")
      return err
    }

    logger.LogDebug("Rendering playbook for worker nodes")
    err = ansible.RenderPlaybook(appDir, clusterName, "worker-nodes",
      clusterType, "worker")

    if err != nil {
      logger.LogError("Error rendering playbook for worker nodes")
      return err
    }
  } else {
    logger.LogDebug("Generating playbooks for single node cluster")
    err := ansible.RenderPlaybook(appDir, clusterName, "localhost",
      clusterType, "")

    if err != nil {
      logger.LogError("Error rendering playbook for single node cluster")
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
    logger.LogDebug("Generating variables file for ha cluster")

    logger.LogDebug("Renderinng vars for lead node")
    err := ansible.GenerateVarsFile(appDir, clusterName, clusterType, "lead", appSettings)
    
    if err != nil {
      logger.LogError("Error rendering vars file for lead node")
      return err
    }

    logger.LogDebug("Error Rendering vars for control nodes")
    err = ansible.GenerateVarsFile(appDir, clusterName, clusterType, "control", appSettings)

    if err != nil {
      logger.LogError("Error rendering vars for control nodes")
      return err
    }

    logger.LogDebug("Rendering vars for worker nodes") 
    err = ansible.GenerateVarsFile(appDir, clusterName, clusterType, "worker", appSettings)

    if err != nil {
      logger.LogError("Error rendering vars for worker nodes")
      return err
    }
  } else {
    logger.LogDebug("Generating variables for single node cluster")
    err := ansible.GenerateVarsFile(appDir, clusterName, clusterType, "", appSettings)

    if err != nil {
      logger.LogError("Error rendering vars for single node cluster")
      return err
    }
  }
  return nil
}

/*
   This will copy all needed scripts to the cluster scripts directories that 
   are needed for a given cluster
*/
func SetupStaticScripts(appDir string, clusterName string) error {
  provisionScriptPath := filepath.Join(appDir, clusterName, "scripts", "provision")
  remoteScriptPath := filepath.Join(appDir, clusterName, "scripts", "remote")

  logger.LogDebug("Getting list of provision scripts to copy")
  provisionScripts, err := static.ListProvisonScripts()

  if err != nil {
    logger.LogError("Error getting list of provision scripts")
    return err
  }

  for _, scriptName := range provisionScripts {
    logger.LogDebug("Copying provision script", "name", scriptName)
    scriptContent, err := static.ReadProvisionScriptFile(scriptName)

    if err != nil {
      logger.LogError("Error getting provision script content", "name", scriptName)
      return err
    }

    path := filepath.Join(provisionScriptPath, scriptName)
    err = os.WriteFile(path, []byte(scriptContent), 0755)

    if err != nil {
      logger.LogError("Error writing script file", "name", scriptName)
      return err
    }
  }

  logger.LogDebug("Getting list of remote scripts to copy")
  remoteScripts, err := static.ListRemoteScripts()

  if err != nil {
    logger.LogError("Error getting list of remote scripts")
    return err
  }

  for _, scriptName := range remoteScripts {
    logger.LogDebug("Copying remote script", "name", scriptName)
    scriptContent, err := static.ReadRemoteScriptFile(scriptName)

    if err != nil {
      logger.LogError("Error getting remote script content", "name", scriptName)
      return err
    }

    path := filepath.Join(remoteScriptPath, scriptName)
    err = os.WriteFile(path, []byte(scriptContent), 0755)

    if err != nil {
      logger.LogError("Error writing script file", "name", scriptName)
      return err
    }
  } 
  return nil 
}

/*
  This will render gather needed data for a vagrant file template and will render 
  the template and will write the result out to the cluster directory 
*/
func RenderVagrantFile(appDir string, clusterName string, appSettings settings.Settings) error {
  providerName := appSettings.Clusters[clusterName].ProviderName
  providerType := appSettings.Providers[providerName].ProviderType
  clusterType := appSettings.Clusters[clusterName].ClusterType

  vagrantFilePath := filepath.Join(appDir, clusterName, "VagrantFile")

  logger.LogDebug("Getting data for VagrantFile rendering", "providerName", providerName, "providerType", providerType)
  logger.LogDebug("VagrantFile path", "path", vagrantFilePath)

  data := make(map[string]interface{})

  data["Provider"] = appSettings.Providers[providerName]

  logger.LogDebug("Provider settings", "settings", data["Provider"])

  if clusterType == "ha" {
    logger.LogDebug("Setting up vagrant template data for ha cluster")

    var leadControlNode []settings.Machine 
    leadControlNode = append(leadControlNode, appSettings.Clusters[clusterName].Leaders[0])

    controlNodes := appSettings.Clusters[clusterName].Leaders[1:]
    workerNodes := appSettings.Clusters[clusterName].Workers

    data["LeadControlNode"] = leadControlNode
    data["ControlNodes"] = controlNodes
    data["WorkerNodes"] = workerNodes

  } else {
    logger.LogDebug("Setting up vagrant template data for single node cluster")

    data["Node"] = appSettings.Clusters[clusterName].Leaders[0]
    logger.LogDebug("Node values are", "node", data["Node"])
  }

  renderedVagrantFile, err := template.RenderVagrantfileTemplate(providerType, clusterType, data)

  if err != nil {
    logger.LogError("Error rendering Vagrantfile")
    return err
  }

  logger.LogDebug("Writing out vagrantfile to cluster directory")
  err = os.WriteFile(vagrantFilePath, []byte(renderedVagrantFile), 0755)

  if err != nil {
    logger.LogError("Error writing vagrantfile to location")
    return err
  }

  return nil
}



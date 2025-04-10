package ansible

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/template"
  "github.com/otiai10/copy"
)

/*
   This will generate the ansible hosts file lead or single node will
   always be localhost
*/
func GenerateAnsibleHostsFile(appDir string, clusterName string, clusterType string, secondaryControlNodes []string, workerNodes []string) error {
  var ansibleHostsContent []string
  clusterDir := filepath.Join(appDir, clusterName)
  ansibleHostsFilePath := filepath.Join(clusterDir, "ansible", "resources", "hosts")

  if clusterType == "single" {
    // since its a single node we only have the lead and only node that 
    // is also the ansible master
    logger.LogDebug("Single node cluster adding only localhost to ansible hosts")
    ansibleHostsContent = append(ansibleHostsContent, "localhost ansible_connection=local")
  } else {
    logger.LogDebug("creating ansible hosts for multi-node cluster")

    // lead node that also acts as the ansible master
    logger.LogDebug("adding lead node as localhost since its ansible master")
    ansibleHostsContent = append(ansibleHostsContent, "[lead-node]")
    ansibleHostsContent = append(ansibleHostsContent, "localhost ansible_connection=local")
    ansibleHostsContent = append(ansibleHostsContent, "")

    // other control nodes
    logger.LogDebug("adding other control nodes")
    ansibleHostsContent = append(ansibleHostsContent, "[control-nodes]")
    for _, cn := range secondaryControlNodes {
      ansibleHostString := cn + " ansible_connection=ssh ansible_user=vagrant"
      
      logger.LogDebug("adding control node", "line", ansibleHostString)
      ansibleHostsContent = append(ansibleHostsContent, ansibleHostString)
    }
    ansibleHostsContent = append(ansibleHostsContent, "")

    // worker nodes
    if len(workerNodes) > 0 {
      logger.LogDebug("adding worker nodes")
      ansibleHostsContent = append(ansibleHostsContent, "[worker-nodes]")
      for _, wn := range workerNodes {
        ansibleHostString := wn + " ansible_connection=ssh ansible_user=vagrant"

        logger.LogDebug("adding worker node", "line", ansibleHostString)
        ansibleHostsContent = append(ansibleHostsContent, ansibleHostString)
      }
    } else {
      logger.LogDebug("no worker nodes to add to the cluster")
    }
  }

  logger.LogDebug("Writing ansible hosts file to path", "path", ansibleHostsFilePath)
  err := os.WriteFile(ansibleHostsFilePath, []byte(strings.Join(ansibleHostsContent, "\n")), 0644)
  if err != nil {
    logger.LogError("Error writing ansible hosts file", "error", err)
    return err
  }
  return nil
}

/*
  This will render the bootstrap.sh script with the desired version of ansible,
  this is the script that will be run to provision ansible on the lead or single node
*/
func RenderBootstrapScript(appDir string, clusterName string, ansibleVersion string) error {
  scriptsDir := filepath.Join(appDir, clusterName, "scripts", "provision")
  boostrapScriptPath := filepath.Join(scriptsDir, "bootstrap.sh")

  templateData := map[string]interface{}{
    "ansibleVersion": ansibleVersion,
  }

  templateContent, err := template.RenderProvisionTemplate("bootstrap", templateData) 
  if err != nil {
    logger.LogError("Error rendering bootstrap script")
    return err
  }

  err = os.WriteFile(boostrapScriptPath, []byte(templateContent), 0644)
  if err != nil {
    logger.LogError("Error writing out bootstrap script")
    return err
  }
  return nil 
}

/*
  This will copy ansible roles from the app wide role repository to the 
  roles directory for the cluster for drive mapping to the ansible 
  machine
*/
func CopyAnsibleRoles(appDir string, clusterName string, rolesNames []string) error {
  appAnsibleRoleDir := filepath.Join(appDir, "ansible-roles")
  clusterAnsibleRoleDir := filepath.Join(appDir, clusterName, "ansible", "roles")

  logger.LogDebug("Copying ansible roles to cluster dir")
  for _, roleName := range rolesNames {
    logger.LogDebug("Copying role", "roleName", roleName)

    src := filepath.Join(appAnsibleRoleDir, roleName)
    dest := filepath.Join(clusterAnsibleRoleDir, roleName)
    err := copy.Copy(src, dest)

    if err != nil {
      logger.LogError("Error copying role to cluster", "role", roleName)
      return err
    }
  }
  return nil
}

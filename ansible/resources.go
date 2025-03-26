package ansible

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/template"
)

func GenerateAnsibleHostsFile(appDir string, clusterName string, clusterType string, secondaryControlNodes []string, workerNodes []string) error {
  var ansibleHostsContent []string
  clusterDir := filepath.Join(appDir, clusterName)
  ansibleHostsFilePath := filepath.Join(clusterDir, "ansible", "resources", "ansible_hosts")

  if clusterType == "single" {
    // since its a single node we only have the lead and only node that 
    // is also the ansible master
    logger.Logger.Debug("Single node cluster adding only localhost to ansible hosts")
    ansibleHostsContent = append(ansibleHostsContent, "localhost ansible_connection=local")
  } else {
    logger.Logger.Debug("creating ansible hosts for multi-node cluster")

    // lead node that also acts as the ansible master
    logger.Logger.Debug("adding lead node as localhost since its ansible master")
    ansibleHostsContent = append(ansibleHostsContent, "[lead-node]")
    ansibleHostsContent = append(ansibleHostsContent, "localhost ansible_connection=local")
    ansibleHostsContent = append(ansibleHostsContent, "")

    // other control nodes
    logger.Logger.Debug("adding other control nodes")
    ansibleHostsContent = append(ansibleHostsContent, "[control-nodes]")
    for _, cn := range secondaryControlNodes {
      ansibleHostString := cn + " ansible_connection=ssh ansible_user=vagrant"
      
      logger.Logger.Debug("adding control node", "line", ansibleHostString)
      ansibleHostsContent = append(ansibleHostsContent, ansibleHostString)
    }
    ansibleHostsContent = append(ansibleHostsContent, "")

    // worker nodes
    if len(workerNodes) > 0 {
      logger.Logger.Debug("adding worker nodes")
      ansibleHostsContent = append(ansibleHostsContent, "[worker-nodes]")
      for _, wn := range workerNodes {
        ansibleHostString := wn + " ansible_connection=ssh ansible_user=vagrant"

        logger.Logger.Debug("adding worker node", "line", ansibleHostString)
        ansibleHostsContent = append(ansibleHostsContent, ansibleHostString)
      }
    } else {
      logger.Logger.Debug("no worker nodes to add to the cluster")
    }
  }

  logger.Logger.Debug("Writing ansible hosts file to path", "path", ansibleHostsFilePath)
  err := os.WriteFile(ansibleHostsFilePath, []byte(strings.Join(ansibleHostsContent, "\n")), 0644)
  if err != nil {
    logger.Logger.Error("Error writing ansible hosts file", "error", err)
    return err
  }
  return nil
}

func RenderBootstrapScript(appDir string, clusterName string, ansibleVersion string, additionalCollections []string) error {
  scriptsDir := filepath.Join(appDir, clusterName, "scripts", "provision")
  boostrapScriptPath := filepath.Join(scriptsDir, "bootstrap.sh")

  templateData := map[string]interface{}{
    "ansibleVersion": ansibleVersion,
    "additionalAnsibleModules": additionalCollections,
  }

  templateContent, err := template.RenderProvisionTemplate("bootstrap", templateData) 
  if err != nil {
    logger.Logger.Error("Error rendering bootstrap script")
    return err
  }

  err = os.WriteFile(boostrapScriptPath, []byte(templateContent), 0644)
  if err != nil {
    logger.Logger.Error("Error writing out bootstrap script")
    return err
  }
  return nil 
}

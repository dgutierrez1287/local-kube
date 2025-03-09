package ansible

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dgutierrez1287/local-kube/cluster"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
)

func generateAnsibleHosts(clusterName string, clusterType string, secondaryControlNodes []string, workerNodes []string) {
  var ansibleHostsContent []string
  appDir := settings.GetAppDirPath()
  clusterDir := filepath.Join(appDir, clusterName)
  ansibleHostsFilePath := filepath.Join(clusterDir, "ansible", "resources", "ansible_hosts")

  if clusterType == "single-node" {
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

    // cleanup
    cluster.FailureCleanup(clusterName)
  }
}

func renderBootstrapScript(clusterName string, ansibleVersion string, additionalCollections []string) {
  appDir := settings.GetAppDirPath()
  scriptsDir := filepath.Join(appDir, clusterName, "scripts", "provision")
  boostrapScriptPath := filepath.Join(scriptsDir, "bootstrap.sh")

  templateData := map[string]interface{}{
    "ansibleVersion": ansibleVersion,
    "additionalAnsibleModules": additionalCollections,
  }

  tmpl, err := template.ParseFiles("bootstrap.tmpl")
  if err != nil {
    logger.Logger.Error("Error parsing bootstrap script template", "error", err)

    // cleanup
    cluster.FailureCleanup(clusterName)
  }

  outfile, err := os.Create(boostrapScriptPath)
  if err != nil {
    logger.Logger.Error("Error creating bootstrap.sh file", "path", boostrapScriptPath, "error", err)

    // cleanup
    cluster.FailureCleanup(clusterName)
  }

  err = tmpl.Execute(outfile, templateData)
  if err != nil {
    logger.Logger.Error("Error rendering bootstrap.sh template", "error", err)

    //cleanup
    cluster.FailureCleanup(clusterName)
  }
}

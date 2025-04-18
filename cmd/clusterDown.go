package cmd

import (
	"fmt"
	"os"

	"github.com/dgutierrez1287/local-kube/cluster"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/output"
	"github.com/dgutierrez1287/local-kube/settings"
  "github.com/dgutierrez1287/local-kube/kubeconfig"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)


var clusterDownCmd = &cobra.Command {
  Use: "cluster-down",
  Short: "Brings a cluster down",
  Long: "Brings a cluster down if it is currently up",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput output.MachineOutput

    if machineOutput && debug {
      logger.Logger.Error("Error you can't have machine output set and debug set")
      os.Exit(20)
    }

    if !machineOutput {
      fmt.Println(util.TitleText)
    }

    appDir := settings.GetAppDirPath()

    logger.LogInfo("Reading settings file")
    // read the settings.json file
    appSettings, err := settings.ReadSettingsFile(appDir)
    if err != nil {
      logger.LogErrorExit("Error reading settings", 200, err)
    }

    // run a check to make sure the cluster is there and figure out how much action is 
    // needed
    created, createdStatus, err := cluster.CheckForExistingCluster(appDir, clusterName, machineOutput)
    if err != nil {
      logger.LogErrorExit("Error checking cluster status", 110, err)
    }

    // If just the directory is present just clear the directory for the cluster
    if created && createdStatus == "directory" {
      logger.LogInfo("Only cluster directory exists, clearing the cluster directory")

      // delete the cluster directory
      err := cluster.DeleteClusterDir(appDir, clusterName)
      if err != nil {
        logger.LogErrorExit("Error deleting the cluster directory", 110, err)
      }

      if !machineOutput {
        logger.LogInfo("Cluster directory deleted, cluster destroyed")
        os.Exit(0)
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.ClusterStatus = "destroyed"
        machineReadableOutput.StatusMessage = "cluster directory deleted, cluster destroyed"
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }

    // If there are machines created run vagrant destroy first and then clear cluster 
    // directory 
    if created && createdStatus == "created" {
        logger.LogInfo("Machines for cluster exit, destroying the cluster")

      // destroy the cluster
      err := cluster.ClusterDown(appDir, clusterName, machineOutput)
      if err != nil {
        logger.LogErrorExit("Error destroying the cluster machines", 110, err)
      }

      // delete the cluster directory
      err = cluster.DeleteClusterDir(appDir, clusterName)
      if err != nil {
        logger.LogErrorExit("Error deleting cluster directory", 110, err)
      }

      logger.LogInfo("Starting kubeconfig update")
      kubeConfigPath := appSettings.KubeConfigPath
      var kubeConfigClusterName string

      if appSettings.Clusters[clusterName].KubeConfigName == "" {
        kubeConfigClusterName = clusterName
      } else {
        kubeConfigClusterName = appSettings.Clusters[clusterName].KubeConfigName
      }

      logger.LogInfo("Backing up current kubeconfig")
      err = kubeconfig.BackupKubeConfig(appDir, kubeConfigPath)
      if err != nil {
        logger.LogErrorExit("Error backing up current kubeconfig", 200, err)
      }

      logger.LogInfo("Reading current kubeconfig")
      destKubeconfig, err := kubeconfig.ReadKubeConfig(kubeConfigPath)
      if err != nil {
        logger.LogErrorExit("Error reading the kubeconfig", 200, err)
      }

      logger.LogInfo("Removing cluster from kubeconfig")
      destKubeconfig.RemoveCluster(kubeConfigClusterName)

      logger.LogInfo("Writing kubeconfig")
      err = kubeconfig.WriteKubeConfig(kubeConfigPath, destKubeconfig)
      if err != nil {
        logger.LogErrorExit("Error writing out new kubeconfig", 200, err)
      }

      logger.LogInfo("Cleaning up kubeconfig backup")
      kubeconfig.CleanKubeConfigBackup(appDir)
      if err != nil {
        logger.LogErrorExit("Error cleaning up kubeconfig backup", 200, err)
      }

      if !machineOutput {
        logger.LogInfo("Cluster Destroyed and cluster directory deleted, cluster is successfully deleted")
        os.Exit(0)
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.ClusterStatus = "destroyed"
        machineReadableOutput.StatusMessage = "cluster destroyed and cluster directory deleted"
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }

    // If there is not part of the cluster present write output to infrom and 
    // exit
    if !created {
      if !machineOutput {
        logger.LogInfo("Cluster does not exist, nothing to destroy")
        os.Exit(0)
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.ClusterStatus = "not created"
        machineReadableOutput.StatusMessage = "Cluster does not exist, nothing to destroy"
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }
    logger.LogInfo("Cluster deleted successfully")
  },
}

func init() {
  // required args for this command
  clusterDownCmd.MarkFlagRequired("cluster")

  // add command
  RootCmd.AddCommand(clusterDownCmd)
}


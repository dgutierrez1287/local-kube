package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/cluster"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var machineSshCmd = &cobra.Command {
  Use: "machine-ssh",
  Short: "Opens ssh session to machine",
  Long: "Opens menu to pick which machine to open ssh to",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(util.TitleText)

    appDir := settings.GetAppDirPath()
    clusterDir := filepath.Join(appDir, clusterName)
    machineName := ""

    // check if cluster exists
    exists, existsType, err := cluster.CheckForExistingCluster(appDir, clusterName,
    machineOutput)
    if err != nil {
      logger.LogErrorExit("Error checking if the cluster exists", 100, err)
    }

    if !exists {
      logger.LogInfo("Cluster does not exists")
      os.Exit(0)
    }

    if existsType == "directory" {
      logger.LogInfo("Cluster exists but not machines are created")
      os.Exit(0)
    }

    // read settings file
    logger.LogInfo("Reading settings file")
    appSettings, err := settings.ReadSettingsFile(appDir)
    if err != nil {
      logger.LogError("Error reading settings file", "error", err)
    }

    if appSettings.Clusters[clusterName].ClusterType == "single" {
      logger.LogInfo("Single node cluster connecting you to the default machine")
      machineName = "default"

    } else {
      machineList := appSettings.Clusters[clusterName].GetMachineNameList()

      prompt := promptui.Select{
        Label: "Select a machine to connect to",
        Items: machineList,
      }

      _, result, err := prompt.Run()
      if err != nil {
        logger.LogErrorExit("Error running prompt for machine selection", 100, err)
      }
      machineName = result
    }

    configs, err := cluster.GetSshConfigs(clusterDir, machineName)
    if err != nil {
      logger.LogErrorExit("Error getting ssh configs", 100, err)
    }
    
    signer, err := cluster.LoadPrivateKey(configs.IdentityFile)
    if err != nil {
      logger.LogErrorExit("Error parsing the ssh key", 100, err)
    }

    err = cluster.OpenSshSession(configs, signer)
    if err != nil {
      logger.LogErrorExit("Error opening ssh session", 100, err)
    }
  },
}

func init() {
  // command specific args

  // required args for this command
  machineSshCmd.MarkFlagRequired("cluster")

  // add command
  RootCmd.AddCommand(machineSshCmd)
}


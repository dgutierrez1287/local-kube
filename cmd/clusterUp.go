package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/ansible"
	"github.com/dgutierrez1287/local-kube/cluster"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/output"
	"github.com/dgutierrez1287/local-kube/settings"
  "github.com/dgutierrez1287/local-kube/kubeconfig"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)

/*
  to generate all files but do not run vagrant up to create VMs,
  This should really only be used for debugging since the cluster
  directory and files are cleared and regenerated everytime you
  use the cluster-up command
*/
var noUp bool 

/*
  to generate all the configs and run vagrant up but do not 
  run the provisioning script. This is mainly used for debugging
  or verification
*/
var noProvision bool

var clusterUpCmd = &cobra.Command {
  Use: "cluster-up",
  Short: "Brings a cluster up",
  Long: "Brings a cluster up if it is not currently up",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput output.MachineOutput

    if !machineOutput {
      fmt.Println(util.TitleText)
    }

    appDir := settings.GetAppDirPath()

    logger.LogInfo("Bringing up cluster", "name", clusterName)
    logger.LogInfo("Running preflight checks")

    // preflight just makes sure app directory is present
    // and there is a settings file present 
    preflight := settings.PreflightCheck(appDir)
    if !preflight {
      logger.LogErrorExit("Error preflight checks failed", 200, nil)
    }
    logger.LogInfo("Preflight Checks passed!")

    logger.LogInfo("Reading settings file")
    // read the settings.json file
    appSettings, err := settings.ReadSettingsFile(appDir)
    if err != nil {
      logger.LogErrorExit("Error reading settings", 200, err)
    }

    logger.LogInfo("Validating settings")
    validSettings := appSettings.SettingsValid(clusterName) 
    if !validSettings {
      logger.LogErrorExit("Error settings could not be validated", 200, nil)
    }

    logger.LogDebug("setting defaults for cluster features")
    appSettings.Clusters[clusterName].ClusterFeatures.SetDefaults(appSettings.Clusters[clusterName].ClusterType,
      appSettings.Clusters[clusterName].Vip)

    logger.LogInfo("Checking to make sure a cluster isn't already present")
    clusterExists, existsType, err := cluster.CheckForExistingCluster(appDir, clusterName, false)
    if err != nil {
      logger.LogErrorExit("Error checking if the cluster exists", 200, err)
    }

    if clusterExists {
      if existsType == "directory" {
        logger.LogInfo("It appears that a cluster directory already exists but no machines exists")
        logger.LogInfo("Clearing out the current cluster directory")
        cluster.DeleteClusterDir(appDir, clusterName)
      } else {
        logger.LogErrorExit("Cluster machines already exist", 100, nil)
      }
    }

    logger.LogInfo("Creating cluster directory and all subdirectories")
    cluster.CreateClusterDirs(appDir, clusterName)

    logger.LogInfo("Generating files for the cluster", "cluster", clusterName)

    // Ansible Resources
    logger.LogInfo("Generating ansible resources")
    err = cluster.GenerateAnsibleResources(appDir, clusterName, appSettings)

    if err != nil {
      logger.LogErrorExit("Error generating ansible resources", 200, err)
    }

    // Roles
    logger.LogInfo("Copying ansible roles to cluster dir")
    roleNames := []string{"kube"}
    err = ansible.CopyAnsibleRoles(appDir, clusterName, roleNames)
    if err != nil {
      logger.LogErrorExit("Error copying ansible roles", 100, err)
    }

    // Playbooks
    logger.LogInfo("Generating ansible playbooks")
    err = cluster.GenerateAnsiblePlaybooks(appDir, clusterName, appSettings)
    if err != nil {
      logger.LogErrorExit("Error generating ansible playbooks", 100, err)
    }

    // Variables
    logger.LogInfo("Generating ansible variables")
    err = cluster.GenerateAnsibleVariables(appDir, clusterName, appSettings)
    if err != nil {
      logger.LogErrorExit("Error generating ansible variables", 100, err)
    }

    // Static Resources
    logger.LogInfo("Copying static scripts to the cluster directories")
    err = cluster.SetupStaticScripts(appDir, clusterName)
    if err != nil {
      logger.LogErrorExit("Error copying static scripts to cluster directory", 100, err)
    }

    // script settings (settings yaml that the machine gets for configuration)
    logger.LogInfo("Generating script settings")
    err = cluster.GenerateScriptSettings(appDir, clusterName, appSettings)
    if err != nil {
      logger.LogErrorExit("Error generating script settings", 100, err)
    }

    // Vagrant file 
    logger.LogInfo("Generating vagrantFile")
    err = cluster.RenderVagrantFile(appDir, clusterName, appSettings)
    if err != nil {
      logger.LogErrorExit("Error generating vagrantfile", 100, err)
    }

    if err != nil {
      logger.LogErrorExit("Error creating vagrant client", 200, err)
    }

    if noUp {
      if !machineOutput {
        logger.LogInfo("noup was set, exiting now that everything has been generated")
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.ClusterStatus = "files generated"
        machineReadableOutput.StatusMessage = "noup was set, exiting since everything is generated"
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }

    logger.LogInfo("Bringing up the cluster")
    _, err = cluster.ClusterUp(appDir, clusterName, machineOutput)
    if err != nil {
      logger.LogErrorExit("Error bringing up the cluster", 100, err)
    }

    if noProvision {
      if !machineOutput {
        logger.LogInfo("no-provision was set, VMs should be up but stopping before provision")
        os.Exit(0)
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.ClusterStatus = "VMs up"
        machineReadableOutput.StatusMessage = "no-provision was set, stopping before provision"
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }

    logger.LogInfo("Provisioning the VMs in the cluster")
    err = cluster.ClusterProvision(appDir, clusterName, appSettings, machineOutput, debug)
    if err != nil {
      logger.LogErrorExit("Error provisioning the cluster machines", 100, err)
    }

    logger.LogInfo("Cluster provisioning complete")

    logger.LogInfo("Starting kubeconfig update")
    kubeConfigPath := appSettings.KubeConfigPath
    var kubeConfigClusterName string

    if appSettings.Clusters[clusterName].KubeConfigName == "" {
      kubeConfigClusterName = clusterName
    } else {
      kubeConfigClusterName = appSettings.Clusters[clusterName].KubeConfigName
    }

    serverUrl := appSettings.Clusters[clusterName].GetServerUrl()
    sourceKubeConfigPath := filepath.Join(appDir, clusterName, "kubeconfig", "k3s.yaml")

    logger.LogDebug("kube config path", "path", kubeConfigPath)
    logger.LogDebug("source kube config path", "path", sourceKubeConfigPath)
    logger.LogDebug("kube cluster name", "name", kubeConfigClusterName)
    logger.LogDebug("server url", "url", serverUrl)

    _, err = os.Stat(sourceKubeConfigPath)
    if err != nil {
      logger.LogErrorExit("Error new cluster kubeconfig does not exist", 200, err)
    }

    logger.LogInfo("Backing up current kubeconfig")
    err = kubeconfig.BackupKubeConfig(appDir, kubeConfigPath)
    if err != nil {
      logger.LogErrorExit("Error backing up current kubeconfig", 200, err)
    }

    logger.LogInfo("Reading destination and source kubeconfig")
    sourceKubeConfig, err := kubeconfig.ReadKubeConfig(sourceKubeConfigPath)
    if err != nil {
      logger.LogErrorExit("Error reading the source kubeconfig", 200, err)
    }

    destKubeconfig, err := kubeconfig.ReadKubeConfig(kubeConfigPath)
    if err != nil {
      logger.LogErrorExit("Error reading the kubeconfig", 200, err)
    }

    logger.LogInfo("Updating server url for new cluster")
    err = sourceKubeConfig.UpdateServerUrl(serverUrl, "default")
    if err != nil {
      logger.LogErrorExit("Error updating server url in kube config", 200, err)
    }

    logger.LogDebug("Updated server url in source config", "url", sourceKubeConfig.Clusters[0].Cluster.Server)

    logger.LogInfo("Adding cluster to kubeconfig")
    err = destKubeconfig.AddCluster(sourceKubeConfig, "default", kubeConfigClusterName)
    if err != nil {
      logger.LogErrorExit("Error adding cluster to kubeconfig", 200, err)
    }

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

    logger.LogInfo("Cluster provisioning complete successfully")
  },
}

func init() {
  // command specific args
  clusterUpCmd.PersistentFlags().BoolVarP(&noUp, "noup", "", false, "Only Generate files but do not create vms")
  clusterUpCmd.PersistentFlags().BoolVarP(&noProvision, "no-provision", "", false, "Create VMs but do not run provision")

  // required args for this command
  clusterUpCmd.MarkFlagRequired("cluster")

  // Add command
  RootCmd.AddCommand(clusterUpCmd)
}

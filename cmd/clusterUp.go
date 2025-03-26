package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/cluster"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)


var clusterUpCmd = &cobra.Command {
  Use: "cluster-up",
  Short: "Brings a cluster up",
  Long: "Brings a cluster up if it is not currently up",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(util.TitleText)

    appDir := settings.GetAppDirPath()
    client, err := cluster.NewVagrantClient(filepath.Join(appDir, clusterName))
    
    if err != nil {
      logger.Logger.Error("Error creating vagrant client factory", "error", err)
    }

    logger.Logger.Info("Bringing up cluster", "name", clusterName)
    logger.Logger.Info("Running preflight checks")

    // preflight just makes sure app directory is present
    // and there is a settings file present 
    preflight := settings.PreflightCheck(appDir)
    if !preflight {
      os.Exit(200)
    }
    logger.Logger.Info("Preflight Checks passed!")

    logger.Logger.Info("Reading settings file")
    // read the settings.json file
    appSettings, err := settings.ReadSettingsFile(appDir)
    if err != nil {
      logger.Logger.Error("Error reading settings", "error", err)
      os.Exit(100)
    }

    logger.Logger.Info("Validating settings")
    validSettings := appSettings.SettingsValid(clusterName) 
    if !validSettings {
      logger.Logger.Error("Error settings could not be validated")
      os.Exit(100)
    }

    logger.Logger.Debug("setting defaults for cluster features")
    appSettings.Clusters[clusterName].ClusterFeatures.SetDefaults(appSettings.Clusters[clusterName].ClusterType,
      appSettings.Clusters[clusterName].Vip)

    logger.Logger.Info("Checking to make sure a cluster isn't already present")
    clusterExists, existsType, err := cluster.CheckForExistingCluster(appDir, clusterName, client)

    if err != nil {
      logger.Logger.Error("There was an error checking if the cluster exists", "error", err)
      os.Exit(100)
    }

    if clusterExists {
      if existsType == "directory" {
        logger.Logger.Info("It appears that a cluster directory already exists but no machines exists")
        logger.Logger.Info("Clearing out the current cluster directory")
        cluster.DeleteClusterDir(appDir, clusterName)
      } else {
        logger.Logger.Error("Cluster machines for this cluster already exist")
        logger.Logger.Info("Please run cluster-down to destroy")
        os.Exit(100)
      }
    }

    logger.Logger.Info("Creating cluster directory and all subdirectories")
    cluster.CreateClusterDirs(appDir, clusterName)

  },
}

func init() {
  // command specific args

  // required args for this command
  clusterUpCmd.MarkFlagRequired("cluster")
}

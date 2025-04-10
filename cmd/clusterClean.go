package cmd

import (
	"fmt"
	"os"

	"github.com/dgutierrez1287/local-kube/cluster"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)

var clusterCleanCmd = &cobra.Command{
  Use: "cluster-clean",
  Short: "Clears a cluster directory",
  Long: "Clears a cluster directory, only use this for testing or if you know no machines are up",
  Run: func(cmd *cobra.Command, args []string) {

    if !machineOutput {
      fmt.Println(util.TitleText)
    }

    appDir := settings.GetAppDirPath()

    logger.LogInfo("Removing the cluster directory")
    err := cluster.DeleteClusterDir(appDir, clusterName)

    if err != nil {
      logger.LogError("Error removing the cluster dir", "error", err)
      os.Exit(100)
    }

    logger.LogInfo("Successfully removed the cluster directory")
    os.Exit(0)
  },
}

func init() {
  // required args for this command
  clusterCleanCmd.MarkFlagRequired("cluster")

  RootCmd.AddCommand(clusterCleanCmd)
}

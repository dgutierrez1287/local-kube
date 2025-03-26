package cmd

import (
	"fmt"
	"os"

	"github.com/dgutierrez1287/local-kube/ansible"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)

var roleCleanCmd = &cobra.Command{
  Use: "roles-clean",
  Short: "Cleans all roles and resets cache",
  Long: "Cleans all roles and resets cache",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(util.TitleText)

    appDir := settings.GetAppDirPath()

    logger.Logger.Info("Cleaning all roles and reseting cache")
    cacheExists, err := ansible.RoleCacheFileExists(appDir)
    if err != nil {
      logger.Logger.Error("There was an error checking if the cache exists", "error", err)
      os.Exit(130)
    }

    if !cacheExists {
      logger.Logger.Info("No role cache is present roles should already be clean")
      os.Exit(130)
    }

    logger.Logger.Info("Deleted all downloaded ansible roles")

    err = ansible.ClearRoles(appDir)
    if err != nil {
      logger.Logger.Error("There was an error clearing all the downloaded roles", "error", err)
      os.Exit(130)
    }

    logger.Logger.Info("Roles deleted successfully")
    logger.Logger.Info("Deleting the roles cache file")

    err = ansible.RoleCacheFileDelete(appDir)
    if err != nil {
      logger.Logger.Error("There was an error deleting the cache file", "error", err)
      os.Exit(130)
    }

    logger.Logger.Info("Roles cache file was successfully deleted")
    os.Exit(0)
  },
}

func init() {
  RootCmd.AddCommand(roleCleanCmd)
}

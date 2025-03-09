package cmd

import (
	"fmt"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/util"
  "github.com/dgutierrez1287/local-kube/settings"
	"github.com/spf13/cobra"
)


var initCmd = &cobra.Command {
  Use: "init",
  Short: "Initializes the system to use local-kube",
  Long: "Initializes the system to use local-kube",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(util.TitleText)

    logger.Logger.Info("Running initialization")

    dirExists := settings.AppDirExists()
    settingsExists := settings.SettingsFileExists()
    var alreadyInit bool

    if dirExists {
      logger.Logger.Debug("App dir already exists")
      alreadyInit = true
    } else {
      logger.Logger.Info("Creating app dir which will be at ~/.local-kube")
      settings.CreateAppDir()
      alreadyInit = false
    }  

    if settingsExists {
      logger.Logger.Debug("Settings file already exists")
      alreadyInit = true
    } else {
      logger.Logger.Info("Creating default settings file")
      settings.CreateDefaultSettingsFile()
      alreadyInit = false
    }

    if alreadyInit {
      logger.Logger.Info("Local-Kube already initialized!")
    } else {
      logger.Logger.Info("Local-Kube initialized!")
    }
  },
}

func init() {
  RootCmd.AddCommand(initCmd)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)


var initCmd = &cobra.Command {
  Use: "init",
  Short: "Initializes the system to use local-kube",
  Long: "Initializes the system to use local-kube",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(util.TitleText)

    logger.Logger.Info("Running initialization")

    logger.Logger.Debug("Getting app directory paths")
    appDir := settings.GetAppDirPath()
    ansibleRoleDir := filepath.Join(appDir, "ansible-roles")

    logger.Logger.Debug("Checking if the directories and settings exist")
    appDirExists := settings.DirectoryExists(appDir)
    ansibleRoleDirExists := settings.DirectoryExists(ansibleRoleDir)
    var alreadyInit bool

    // Top application directory
    if appDirExists {
      logger.Logger.Debug("App directory already exists")
      alreadyInit = true
    } else {
      logger.Logger.Info("Creating app directory which will be at ~/.local-kube")
      settings.CreateDirectory(appDir)
      alreadyInit = false
    }  

    // ansible role directory
    if ansibleRoleDirExists {
      logger.Logger.Debug("Ansible role directory already exists")
      alreadyInit = true
    } else {
      logger.Logger.Info("Creating ansible role directory")
      settings.CreateDirectory(ansibleRoleDir)
      alreadyInit = false
    }

    // Settings file
    settingsExists, err := settings.SettingsFileExists()
    if err != nil {
      logger.Logger.Error("Error checking for the settings file")
      os.Exit(100)
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

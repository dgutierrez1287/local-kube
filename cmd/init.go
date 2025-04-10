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
    if !machineOutput {
      fmt.Println(util.TitleText)
    }

    logger.LogInfo("Running initialization")

    logger.LogDebug("Getting app directory paths")
    appDir := settings.GetAppDirPath()
    ansibleRoleDir := filepath.Join(appDir, "ansible-roles")

    logger.LogDebug("Checking if the directories and settings exist")
    appDirExists := settings.DirectoryExists(appDir)
    ansibleRoleDirExists := settings.DirectoryExists(ansibleRoleDir)
    var alreadyInit bool

    // Top application directory
    if appDirExists {
      logger.LogDebug("App directory already exists")
      alreadyInit = true
    } else {
      logger.LogInfo("Creating app directory which will be at ~/.local-kube")
      settings.CreateDirectory(appDir)
      alreadyInit = false
    }  

    // ansible role directory
    if ansibleRoleDirExists {
      logger.LogDebug("Ansible role directory already exists")
      alreadyInit = true
    } else {
      logger.LogInfo("Creating ansible role directory")
      settings.CreateDirectory(ansibleRoleDir)
      alreadyInit = false
    }

    // Settings file
    settingsExists, err := settings.SettingsFileExists(appDir)
    if err != nil {
      logger.LogError("Error checking for the settings file")
      os.Exit(100)
    }

    if settingsExists {
      logger.LogDebug("Settings file already exists")
      alreadyInit = true
    } else {
      logger.LogInfo("Creating default settings file")
      settings.CreateDefaultSettingsFile(appDir)
      alreadyInit = false
    }

    if alreadyInit {
      logger.LogInfo("Local-Kube already initialized!")
    } else {
      logger.LogInfo("Local-Kube initialized!")
    }
  },
}

func init() {
  RootCmd.AddCommand(initCmd)
}

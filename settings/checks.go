package settings

import (
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
)

// Run basic checks to make sure init has been run and needed files
// are there
func PreflightCheck(appDir string) bool {
  appDirExist := DirectoryExists(appDir)
  ansibleRoleDirExist := DirectoryExists(filepath.Join(appDir, "ansible-roles"))
  settingsExist, err := SettingsFileExists(appDir)

  if err != nil {
    logger.LogError("Error checking if settings file exists", "error", err)
    os.Exit(100)
  }

  if !appDirExist {
    logger.LogError("Preflight check failed, app directory does not exist")
    return false
  }

  if !ansibleRoleDirExist {
    logger.LogError("Preflight check failed, ansible-role directory does not exist")
    return false
  }

  if !settingsExist {
    logger.LogError("Preflight check failed, settings file does not exist")
    return false 
  }
  return true
}



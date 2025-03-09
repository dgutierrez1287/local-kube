package settings

import (
	"github.com/dgutierrez1287/local-kube/logger"
)



// Run basic checks to make sure init has been run and needed files 
// are there 
func PreflightCheck() bool {
  dirExist := AppDirExists()
  settingsExist := SettingsFileExists()

  if !dirExist || !settingsExist {
    logger.Logger.Error("Preflight check failed, app directory or settings file is missing, have you run Init?")
    return false
  }
  return true
}



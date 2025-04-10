package settings

import (
  "os"
  "path/filepath"
  "errors"

  "github.com/dgutierrez1287/local-kube/logger"
)

/*
These are functions for managing application level directories
this is the main application directory (home/username/.local-kube)
and also any other application directories which live longer and are 
independent from clusters
*/

/* 
Get the path for the application directory
this will be $HOME/$Username/.local-kube and 
the equilvalent on windows
*/
func GetAppDirPath() string {
  userDir, err := os.UserHomeDir()
  if err != nil {
    logger.LogError("Error getting user home dir", "error", err)
    os.Exit(120)
  }

  return filepath.Join(userDir, ".local-kube")
}

/*
This will check if an application level directory 
exists
*/
func DirectoryExists(dirPath string) bool {
  var dirExists bool

  if _, err := os.Stat(dirPath); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.LogDebug("directory does not exist", "path", dirPath)
      dirExists = false
    } else {
      logger.LogError("Error checking if directory exists", "path", dirPath, "error", err)
      os.Exit(120)
    }
  } else {
    logger.LogDebug("Directory exists")
    dirExists = true
  }

  return dirExists
}

/*
Creates an application level directory
*/
func CreateDirectory(dirPath string) {

  err := os.Mkdir(dirPath, 0750)
  if err != nil {
    logger.LogError("Error creating directory", "path", dirPath,"error", err)
    os.Exit(123)
  }
  logger.LogInfo("Directory created")
}


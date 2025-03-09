package settings

import (
  "os"
  "path/filepath"
  "errors"

  "github.com/dgutierrez1287/local-kube/logger"
)

// Get the path for the app dir
func GetAppDirPath() string {
  userDir, err := os.UserHomeDir()
  if err != nil {
    logger.Logger.Error("Error getting user home dir", "error", err)
    os.Exit(120)
  }

  return filepath.Join(userDir, ".local-kube")
}

// Check if app dir exists
func AppDirExists() bool {
  appDir := GetAppDirPath()
  var dirExists bool

  if _, err := os.Stat(appDir); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.Logger.Debug("App dir does not exist")
      dirExists = false
    } else {
      logger.Logger.Error("Error checking if app dir exists", "error", err)
      os.Exit(120)
    }
  } else {
    logger.Logger.Debug("App Dir exists")
    dirExists = true
  }

  return dirExists
}

// Create the application dir
func CreateAppDir() {
  appDir := GetAppDirPath()

  err := os.Mkdir(appDir, 0750)
  if err != nil {
    logger.Logger.Error("Error creating app dir", "error", err)
    os.Exit(123)
  }
  logger.Logger.Info("App dir created")
}


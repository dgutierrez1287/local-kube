package cluster

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
)

var ansibleDirs = []string {
  "roles",
  "playbooks",
  "variables",
  "resources",
}

var scriptsDirs = []string {
  "provision",
  "remote",
}

var mainDirs = []string {
  "ansible",
  "scripts",
  "kubeconfig",
  "settings",
}

// check if the cluster dir exists
func ClusterDirExists(appDir string, clusterName string) (bool, error) {
  clusterDir := filepath.Join(appDir, clusterName)
  var dirExists bool

  if _, err := os.Stat(clusterDir); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.Logger.Debug("Cluster dir does not exists", "cluster", clusterName)
      dirExists = false
    } else {
      logger.Logger.Error("Error checking if cluster dir exists", "cluster", clusterName, "error", err)
      return false, err
    }
  } else {
    logger.Logger.Debug("Cluster dir exists", "cluster", clusterName)
    dirExists = true
  }
  return dirExists, nil
}

// create all the directories needed for a cluster and 
// provisioning 
func CreateClusterDirs(appDir string, clusterName string) error {
  clusterDir := filepath.Join(appDir, clusterName)

  // create the cluster directory
  err := os.Mkdir(clusterDir, 0750)
  if err != nil {
    logger.Logger.Error("Error creating cluster dir")
    return err
  }
  logger.Logger.Info("Cluster directory created", "cluster", clusterName)

  // create all the top level dirs in the cluster dir
  logger.Logger.Debug("Creating all main cluster directories")
  err = createClusterSubDirs(clusterDir, mainDirs)
  if err != nil {
    logger.Logger.Error("Error creating main directories for cluster", "cluster", clusterName)
    return err
  }

  // create all the ansible directories
  logger.Logger.Debug("Creating ansible sub directories")
  err = createClusterSubDirs(filepath.Join(clusterDir, "ansible"), ansibleDirs)
  if err != nil {
    logger.Logger.Error("Error creating ansible directories for cluster", "cluster", clusterName)
    return err
  }

  // create all the script directories
  logger.Logger.Debug("Creating script sub directories")
  err = createClusterSubDirs(filepath.Join(clusterDir, "scripts"), scriptsDirs)
  if err != nil {
    logger.Logger.Error("Error creating script directories for cluster", "cluster", clusterName)
    return err
  }

  return nil
}

func createClusterSubDirs(parentDir string, subDirs []string) error {
  for _, dir := range subDirs {
    logger.Logger.Debug("creating sub dir", "dir", dir, "parentDir", parentDir)
    dirPath := filepath.Join(parentDir, dir)
    
    err := os.Mkdir(dirPath, 0750)
    if err != nil {
      logger.Logger.Error("Error creating sub dir", "path", dirPath)
      return err
    }
    logger.Logger.Debug("Sub directory created", "path", dirPath)
  }
  return nil 
}

// Deletes the cluster directory, this
// clears the cluster for the next time its run
// all files needed for the cluster is created when its needed
func DeleteClusterDir(appDir string, clusterName string) error {
  clusterDir := filepath.Join(appDir, clusterName)

  err := os.RemoveAll(clusterDir)
  if err != nil {
    logger.Logger.Error("Error removing the cluster dir", "cluster", clusterName, "error", err)
    return err
  }
  logger.Logger.Info("Cluster directory has been deleted", "cluster", clusterName)
  return nil
}

// clean up after an error 
// logs the cleanup, cleanup directories
// exits
func FailureCleanup(appDir, clusterName string) {
  logger.Logger.Info("Cleaning up cluster directory", "cluster", clusterName)
  DeleteClusterDir(appDir, clusterName)
  os.Exit(200)
}

package cluster

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
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
  "resources",
}

// check if the cluster dir exists
func ClusterDirExists(clusterName string) bool {
  appDir := settings.GetAppDirPath()
  clusterDir := filepath.Join(appDir, clusterName)
  var dirExists bool

  if _, err := os.Stat(clusterDir); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.Logger.Debug("Cluster dir does not exists", "cluster", clusterName)
      dirExists = false
    } else {
      logger.Logger.Error("Error checking if cluster dir exists", "cluster", clusterName, "error", err)
      os.Exit(120)
    }
  } else {
    logger.Logger.Debug("Cluster dir exists", "cluster", clusterName)
    dirExists = true
  }
  return dirExists
}

// create all the directories needed for a cluster and 
// provisioning 
func CreateClusterDirs(clusterName string) {
  appDir := settings.GetAppDirPath()
  clusterDir := filepath.Join(appDir, clusterName)

  // create the cluster directory
  err := os.Mkdir(clusterDir, 0750)
  if err != nil {
    logger.Logger.Error("Error creating cluster dir")
    os.Exit(120)
  }
  logger.Logger.Info("Cluster directory created", "cluster", clusterName)

  // create all the top level dirs in the cluster dir
  for _, dir := range mainDirs {
    dirPath := filepath.Join(clusterDir, dir)
    err := os.Mkdir(dirPath, 0750)
    if err != nil {
      logger.Logger.Error("Error creating main directory", "directory", dirPath)
      os.Exit(120)
    }
    logger.Logger.Debug("Main directory created", "cluster", clusterName, "directory", dirPath)
  }
  logger.Logger.Info("Main directories created", "cluster", clusterName)

  // create all the ansible directories
  for _, dir := range ansibleDirs {
    dirPath := filepath.Join(clusterDir, "ansible", dir)
    err := os.Mkdir(dirPath, 0750)
    if err != nil {
      logger.Logger.Error("Error creating ansible directory", "directory", dirPath)
      os.Exit(120)
    }
    logger.Logger.Debug("Ansible directory created", "cluster", clusterName, "directory", dirPath)
  }
  logger.Logger.Info("Ansible directories created", "cluster", clusterName)

  // create all the script directories
  for _, dir := range scriptsDirs {
    dirPath := filepath.Join(clusterDir, "scripts", dir)
    err := os.Mkdir(dirPath, 0750)
    if err != nil {
      logger.Logger.Error("Error creating script directory", "directory", dirPath)
      os.Exit(120)
    }
    logger.Logger.Debug("Script directory created", "cluster", clusterName, "directory", dirPath)
  }
  logger.Logger.Info("Script directories created", "cluster", clusterName)
}

// Deletes the cluster directory, this
// clears the cluster for the next time its run
// all files needed for the cluster is created when its needed
func DeleteClusterDir(clusterName string) {
  appDir := settings.GetAppDirPath()
  clusterDir := filepath.Join(appDir, clusterName)

  err := os.RemoveAll(clusterDir)
  if err != nil {
    logger.Logger.Error("Error removing the cluster dir", "cluster", clusterName, "error", err)
    os.Exit(120)
  }
  logger.Logger.Info("Cluster directory has been deleted", "cluster", clusterName)
}

// clean up after an error 
// logs the cleanup, cleanup directories
// exits
func FailureCleanup(clusterName string) {
  logger.Logger.Info("Cleaning up cluster directory", "cluster", clusterName)
  DeleteClusterDir(clusterName)
  os.Exit(200)
}

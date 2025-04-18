package cluster

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
)

// ansible directories
var ansibleDirs = []string {
  "roles",
  "playbooks",
  "variables",
  "resources",
}

// script directories
var scriptsDirs = []string {
  "provision",
  "remote",
}

// main directories for the cluster
var mainDirs = []string {
  "ansible",
  "scripts",
  "logs",
  "kubeconfig",
  "settings",
}

/*
check if the cluster dir exists
*/
func ClusterDirExists(appDir string, clusterName string) (bool, error) {
  clusterDir := filepath.Join(appDir, clusterName)
  var dirExists bool

  if _, err := os.Stat(clusterDir); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.LogDebug("Cluster dir does not exists", "cluster", clusterName)
      dirExists = false
    } else {
      logger.LogError("Error checking if the cluster directory exists")
      return false, err
    }
  } else {
    logger.LogDebug("Cluster dir exists", "cluster", clusterName)
    dirExists = true
  }
  return dirExists, nil
}

/*
create all the directories needed for a cluster and 
provisioning 
*/
func CreateClusterDirs(appDir string, clusterName string) error {
  clusterDir := filepath.Join(appDir, clusterName)

  // create the cluster directory
  err := os.Mkdir(clusterDir, 0750)
  if err != nil {
    logger.LogError("Error creating the cluster directory")
    return err
  }
  logger.LogInfo("Cluster directory created", "cluster", clusterName)

  // create all the top level dirs in the cluster dir
  logger.LogDebug("Creating all main cluster directories")
  err = createClusterSubDirs(clusterDir, mainDirs)
  if err != nil {
    logger.LogError("Error creating the main directories for the cluster")
    return err
  }

  // create all the ansible directories
  logger.LogDebug("Creating ansible sub directories")
  err = createClusterSubDirs(filepath.Join(clusterDir, "ansible"), ansibleDirs)
  if err != nil {
    logger.LogError("Error creating ansible directories for cluster")
    return err
  }

  // create all the script directories
  logger.LogDebug("Creating script sub directories")
  err = createClusterSubDirs(filepath.Join(clusterDir, "scripts"), scriptsDirs)
  if err != nil {
    logger.LogError("Error creating script directories for cluster")
    return err
  }

  return nil
}

/*
creates sub directories for a given parent directory
*/
func createClusterSubDirs(parentDir string, subDirs []string) error {
  for _, dir := range subDirs {
    logger.LogDebug("creating sub dir", "dir", dir, "parentDir", parentDir)
    dirPath := filepath.Join(parentDir, dir)
    
    err := os.Mkdir(dirPath, 0750)
    if err != nil {
      logger.LogError("Error creating sub dir")
      return err
    }
    logger.LogDebug("Sub directory created", "path", dirPath)
  }
  return nil 
}

/*
Deletes the cluster directory, this
clears the cluster for the next time its run
all files needed for the cluster is created when its needed
*/
func DeleteClusterDir(appDir string, clusterName string) error {
  clusterDir := filepath.Join(appDir, clusterName)

  err := os.RemoveAll(clusterDir)
  if err != nil {
    logger.LogError("Error removing the cluster directory")
    return err
  }
  logger.LogInfo("Cluster directory has been deleted", "cluster", clusterName)
  return nil
}

/*
clean up after an error 
logs the cleanup, cleanup directories
exits
*/
func FailureCleanup(appDir, clusterName string) {
  logger.LogInfo("Cleaning up cluster directory", "cluster", clusterName)
  DeleteClusterDir(appDir, clusterName)
  os.Exit(200)
}

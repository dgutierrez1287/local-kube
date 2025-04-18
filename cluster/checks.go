package cluster

import (
	"errors"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
)

/* This is going to be used by a few different
top level commands, it will check if the
cluster directory exists and if so it will
then check if the vagrant machines are already
created and if they are running

it has 2 return values

a boolean if anything for the cluster exists and
a string that has two possible values:
directory - if the cluster directory exists but no machine
exists for the cluster

created - if the cluster has any machines from it that are
in any state
*/
func CheckForExistingCluster(appDir string, clusterName string, 
machineOutput bool) (bool, string, error){
  var dirExists bool
  machineCreated := false
  clusterDir := filepath.Join(appDir, clusterName)

  logger.LogDebug("checking if any cluster currently exists")

  dirExists, err := ClusterDirExists(appDir, clusterName)
  if err != nil {
    logger.LogError("Error checking for the cluster directory")
    return false, "", err
  }
  
  // if the cluster dir doesnt exist stopping looking 
  //now and return
  if !dirExists {
    logger.LogDebug("cluster directory doesnt exist there is no cluster present")
    return false, "", nil
  }

  logger.LogDebug("Getting the vagrant client")
  client, err := NewVagrantClient(clusterDir)
  if err != nil {
    logger.LogError("Error getting vagrant client")
    return false, "", err
  }

  logger.LogDebug("checking for current cluster that has been created")

  statusCmd := client.Status()

  logger.LogDebug("statusCmd", statusCmd)
  
  if statusCmd == nil {
    logger.LogError("Error status command is nil")
    return false, "", errors.New("status command is nil")
  }
  
  if err := statusCmd.Start(); err != nil {
    logger.LogError("Error running the vagrant status command")
    return false, "", err
  }

  if err := statusCmd.Wait(); err != nil {
    logger.LogError("Error waiting for the status command to return")
    return false, "", err
  }

  resp := statusCmd.StatusResponse
  logger.LogDebug("response", resp)
  respErrors := resp.ErrorResponse

  if respErrors.Error != nil {
    logger.LogError("Error getting the vagrant status")
    return false, "", respErrors.Error
  }

  statuses := resp.Status

  for name, status := range statuses {
    logger.LogDebug("machine status", "name", name, "status", status)
    if status != "not_created" {
      logger.LogDebug("machine exists", "name", name)
      machineCreated = true
    }
  }

  if machineCreated {
    logger.LogDebug("Directory and machine(s) are created")
    return true, "created", nil
  }

  logger.LogDebug("Directory for cluster is created, but no machines are created")
  return true, "directory", nil
}

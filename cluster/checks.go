package cluster

import (
	"errors"

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
client VagrantClientInterface, machineOutput bool) (bool, string, error){
  var dirExists bool
  machineCreated := false

  logger.Logger.Debug("checking if any cluster currently exists")

  dirExists, err := ClusterDirExists(appDir, clusterName)
  if err != nil {
    if !machineOutput {
      logger.Logger.Error("Error checking for cluster dir", "cluster", clusterName)
    }
    return false, "", err
  }
  
  // if the cluster dir doesnt exist stopping looking 
  //now and return
  if !dirExists {
    logger.Logger.Debug("cluster directory doesnt exist there is no cluster present")
    return false, "", nil
  }

  logger.Logger.Debug("checking for current cluster that has been created")

  statusCmd := client.Status()

  logger.Logger.Debug("statusCmd", statusCmd)
  
  if statusCmd == nil {
    if !machineOutput {
      logger.Logger.Error("Error status command is nil")
    }
    return false, "", errors.New("status command is nil")
  }
  
  if err := statusCmd.Start(); err != nil {
    if !machineOutput {
      logger.Logger.Error("Error running the vagrant status command", "error", err)
    }
    return false, "", err
  }

  if err := statusCmd.Wait(); err != nil {
    if !machineOutput {
      logger.Logger.Error("Error waiting for the status command to return", "error", err)
    }
    return false, "", err
  }

  resp := statusCmd.StatusResponse
  logger.Logger.Debug("response", resp)
  respErrors := resp.ErrorResponse

  if respErrors.Error != nil {
    if !machineOutput {
      logger.Logger.Error("Error getting the vagrant status", "error", respErrors.Error)
    }
    return false, "", respErrors.Error
  }

  statuses := resp.Status

  for name, status := range statuses {
    logger.Logger.Debug("machine status", "name", name, "status", status)
    if status != "not_created" {
      logger.Logger.Debug("machine exists", "name", name)
      machineCreated = true
    }
  }

  if machineCreated {
    logger.Logger.Debug("Directory and machine(s) are created")
    return true, "created", nil
  }

  logger.Logger.Debug("Directory for cluster is created, but no machines are created")
  return true, "directory", nil
}

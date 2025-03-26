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
client VagrantClientInterface) (bool, string, error){
  var dirExists bool

  logger.Logger.Debug("checking if any cluster currently exists")

  dirExists, err := ClusterDirExists(appDir, clusterName)
  if err != nil {
    logger.Logger.Error("Error checking for cluster dir", "cluster", clusterName)
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
    logger.Logger.Error("Error status command is nil")
    return false, "", errors.New("status command is nil")
  }
  
  if err := statusCmd.Start(); err != nil {
    logger.Logger.Error("Error running the vagrant status command", "error", err)
    return false, "", err
  }

  if err := statusCmd.Wait(); err != nil {
    logger.Logger.Error("Error waiting for the status command to return", "error", err)
    return false, "", err
  }

  resp := statusCmd.StatusResponse
  respErrors := resp.ErrorResponse

  if respErrors.Error != nil {
    logger.Logger.Error("Error getting the vagrant status", "error", respErrors.Error)
    return false, "", respErrors.Error
  }

  statuses := resp.Status

  if len(statuses) == 0 {
    logger.Logger.Debug("Directory for cluster exists but no vms for cluster are present")
    return true, "directory", nil
  }

  logger.Logger.Debug("Directory for cluster and vms are present, punt this back to the user")
  return true, "created", nil
}

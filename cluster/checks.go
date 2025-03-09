package cluster

import (
	"os"
	"path/filepath"

	vagrant "github.com/bmatcuk/go-vagrant"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
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
func CheckForExistingCluster(clusterName string) (bool, string){
  var dirExists bool

  logger.Logger.Debug("checking if any cluster currently exists")

  appDir := settings.GetAppDirPath()
  clusterDirPath := filepath.Join(appDir, clusterName)

  dirExists = ClusterDirExists(clusterName)
  
  // if the cluster dir doesnt exist stopping looking 
  //now and return
  if !dirExists {
    logger.Logger.Debug("cluster directory doesnt exist there is no cluster present")
    return false, ""
  }

  logger.Logger.Debug("checking for current cluster that has been created")

  // get a vagrant client for that cluster
  client, err := vagrant.NewVagrantClient(clusterDirPath)
  if err != nil {
    logger.Logger.Error("Error getting vagrant client", "error", err)
    os.Exit(400)
  }

  statusCmd := client.Status()
  
  if err := statusCmd.Start(); err != nil {
    logger.Logger.Error("Error running the vagrant status command", "error", err)
    os.Exit(400)
  }

  if err := statusCmd.Wait(); err != nil {
    logger.Logger.Error("Error waiting for the status command to return", "error", err)
    os.Exit(400)
  }

  resp := statusCmd.StatusResponse
  respErrors := resp.ErrorResponse

  if respErrors.Error != nil {
    logger.Logger.Error("Error getting the vagrant status", "error", statusCmd.Error)
    os.Exit(400)
  }

  statuses := resp.Status

  if len(statuses) == 0 {
    logger.Logger.Debug("Directory for cluster exists but no vms for cluster are present")
    return true, "directory"
  }

  logger.Logger.Debug("Directory for cluster and vms are present, punt this back to the user")
  return true, "created"
}

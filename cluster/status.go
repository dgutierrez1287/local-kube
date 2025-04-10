package cluster

import (
	"errors"

	"github.com/dgutierrez1287/local-kube/logger"
)


/*
This will return a more detailed status of all the machines in the cluster
and also a cluster status, ex if any machines are running the cluster status
would be running
*/
func GetDetailedClusterStatus(appDir string, clusterName string, 
client VagrantClientInterface, machineOutput bool) (string, map[string]string, error) {
  machineCount := 0
  machineRunningCount := 0
  machinePauseCount := 0
  machineStoppedCount := 0
  var clusterStatus string

  logger.LogDebug("Checking status of the cluster")

  statusCmd := client.Status()

  logger.LogDebug("statusCmd", statusCmd)

  if statusCmd == nil {
    logger.LogError("Error status command is nil")
    return "", nil, errors.New("status command is nil")
  }

  err := statusCmd.Start()
  if err != nil {
    logger.LogError("Error running the vagrant status command")
    return "", nil, err
  }

  err = statusCmd.Wait()
  if err != nil {
    logger.LogError("Error waiting for vagrant status command")
    return "", nil, err
  }

  resp := statusCmd.StatusResponse
  logger.LogDebug("response", resp)

  respErrors := resp.ErrorResponse
  
  if respErrors.Error != nil {
    logger.LogError("Error getting the vagrant status")
    return "", nil, respErrors.Error
  }

  machineCount = len(resp.Status)

  for name, status := range resp.Status {
    logger.LogDebug("machine status", "name", name, "status", status)
    switch status {
    case "running":
      machineRunningCount++
    case "paused":
      machinePauseCount++
    case "poweroff":
      machineStoppedCount++
    }
  }

  logger.LogDebug("machineCount", machineCount, "running", machineRunningCount, "paused", machinePauseCount, "poweroff", machineStoppedCount)
  if machineCount == machineRunningCount {
    clusterStatus = "running"
  } else if machineCount == machinePauseCount {
    clusterStatus = "paused"  
  } else if machineCount == machineStoppedCount {
    clusterStatus = "poweroff"
  } else {
    clusterStatus = "patially_running"
  }

  logger.LogDebug("clusterStatus", clusterStatus)
  return clusterStatus, resp.Status, nil
}



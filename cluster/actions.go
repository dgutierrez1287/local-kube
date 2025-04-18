package cluster

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	vagrant "github.com/bmatcuk/go-vagrant"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
)

/*
This will spin up a Cluster but will not run any ansible
provisioning
*/
func ClusterUp(appDir string, clusterName string, 
machineOutput bool) (map[string]*vagrant.VMInfo ,error) {
  clusterDir := filepath.Join(appDir, clusterName)

  logger.LogDebug("Getting vagrant client")
  client, err := NewVagrantClient(clusterDir)
  if err != nil {
    logger.LogError("Error getting vagrant client")
    return nil, err
  }

  logger.LogDebug("Bringing up the cluster", "name", clusterName)

  upCmd := client.Up()
  upCmd.DestroyOnError = false
  upCmd.InstallProvider = true

  logger.LogDebug("upCmd", upCmd)

  if upCmd == nil {
    logger.LogError("Error up command is nil")
    return nil, errors.New("up command is nil")
  }

  err = upCmd.Start()
  if err != nil {
    logger.LogError("Error running the vagrant up command")
    return nil, err
  }

  err = upCmd.Wait()
  if err != nil {
    logger.LogError("Error waiting for vagrant up command")
    return nil, err
  }

  resp := upCmd.UpResponse
  logger.LogDebug("response", resp)

  respErrors := resp.ErrorResponse

  if respErrors.Error != nil {
    logger.LogError("Error bringing up the vagrant stack")
    return nil, respErrors.Error
  }

  return resp.VMInfo, nil
}

/*
This will ssh to the ansible(lead) node in the cluster and run a provision script
that will run ansible for a given cluster type
*/
func ClusterProvision(appDir string, clusterName string,
appSettings settings.Settings, machineOutput bool, debug bool) error {
  clusterDir := filepath.Join(appDir, clusterName)
  cmdStr := ""

  clusterType := appSettings.Clusters[clusterName].ClusterType
  vagrantNodeName := appSettings.Clusters[clusterName].GetAnsibleNodeVagrantName()
  if debug {
    cmdStr = fmt.Sprintf("bash /scripts/%s-provision.sh debug", clusterType)
  } else {
    cmdStr = fmt.Sprintf("bash /scripts/%s-provision.sh info", clusterType)
  }

  sshConfig, err := GetSshConfigs(clusterDir, vagrantNodeName)
  if err != nil {
    logger.LogError("Error getting vagrant ssh config")
    return err
  }

  sshArgs := []string {
    "-i", sshConfig.IdentityFile,
    "-p", fmt.Sprintf("%d", sshConfig.Port),
    "-o", "StrictHostKeyChecking=no",
    "-o", "UserKnownHostsFile=/dev/null",
    fmt.Sprintf("%s@%s", sshConfig.User, sshConfig.HostName),
    cmdStr,
  }

  sshArgs = append(sshArgs, cmdStr)

  logger.LogDebug("ssh args", "args", sshArgs)

  cmd := exec.Command("ssh", sshArgs...)

  output, err := cmd.CombinedOutput()
  if err != nil {
    logger.LogError("Provision command failed")
    return err
  }

  logger.LogDebug("Ssh output", "output", output)
  return nil
}

/*
This will destroy a given Cluster
*/
func ClusterDown(appDir string, clusterName string, 
machineOutput bool) error {
  clusterDir := filepath.Join(appDir, clusterName)

  logger.LogDebug("Getting vagrant client")
  client, err := NewVagrantClient(clusterDir)
  if err != nil {
    logger.LogError("Error getting vagrant client")
    return err
  }
  
  logger.LogDebug("Destroying the cluster", "name", clusterName)

  destroyCmd := client.Destroy()

  logger.LogDebug("destroyCmd", destroyCmd)

  if destroyCmd == nil {
    logger.LogError("Error destroy command is nil")
    return errors.New("destroy command is nil")
  }

  err = destroyCmd.Start()
  if err != nil {
    logger.LogError("Error running the vagrant destroy command")
    return err
  }

  err = destroyCmd.Wait()
  if err != nil {
    logger.LogError("Error waiting for the vagrant destroy command")
    return err
  }

  respErrors := destroyCmd.ErrorResponse

  if respErrors.Error != nil {
    logger.LogError("Error destroying the vagrant stack")
    return respErrors.Error
  }
  return nil
}



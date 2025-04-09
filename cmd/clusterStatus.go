package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/cluster"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/output"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)

var clusterStatusCmd = &cobra.Command{
  Use: "cluster-status",
  Short: "Gets the status of a cluster",
  Long: "Gets the status of a cluster",
  Run: func(cmd *cobra.Command, args []string) {
    var machineReadableOutput output.MachineOutput

    logger.Logger.Debug("machineOutput flag is", "value", machineOutput)

    if machineOutput && debug {
      logger.Logger.Error("Error you can't have machine output set and debug set")
      os.Exit(20)
    }

    if !machineOutput {
      fmt.Println(util.TitleText)
    }

    appDir := settings.GetAppDirPath()
    clusterDir := filepath.Join(appDir, clusterName)

    if !machineOutput {
      logger.Logger.Info("Getting Vagrant cluster for cluster", "name", clusterName, "dir", clusterDir)
    }

    // get vagrant client and check for error 
    vagrantClient, err := cluster.NewVagrantClient(clusterDir)
    if err != nil {
      if !machineOutput {
        logger.Logger.Error("Error getting vagrant client for cluster", "name", clusterName, "error", err)
        os.Exit(100)
      } else {
        machineReadableOutput.ExitCode = 100
        machineReadableOutput.ErrorMessage = fmt.Sprintf("Error getting vagrant client %v", err)
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }
    
    // Run an initial check if the cluster directory exists and see if any machines are present
    created, createdStatus, err := cluster.CheckForExistingCluster(appDir, clusterName, vagrantClient, machineOutput)
    if err != nil {
      if !machineOutput {
        logger.Logger.Error("Error checking for the initial cluster status", "error", err)
        os.Exit(110)
      } else {
        machineReadableOutput.ExitCode = 110
        machineReadableOutput.ErrorMessage = fmt.Sprintf("Error checking inital cluster status %v", err)
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }

    // Just the directory exists but no machines are present
    if created && createdStatus == "directory" {
      if !machineOutput {
        logger.Logger.Info("Cluster directory exists but no machines are created")
        os.Exit(0)
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.DirectoryCreated = true
        machineReadableOutput.ClusterStatus = "directory created"
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }

    // machines are present, call to get a more detailed status
    if created && createdStatus == "created" {
      if !machineOutput {
        logger.Logger.Info("Getting detailed cluster and machine status")
      }

      // Get detailed cluster and machine status for output
      clusterStatus, statuses, err := cluster.GetDetailedClusterStatus(appDir, clusterName, vagrantClient, machineOutput)
      if err != nil {
        if !machineOutput {
          logger.Logger.Error("Error getting detailed cluster status", "error", err)
          os.Exit(110)
        } else {
          machineReadableOutput.ExitCode = 110
          machineReadableOutput.ErrorMessage = fmt.Sprintf("Error getting detailed machine status %v", err)
          output, eCode := machineReadableOutput.GetMachineOutputJson()
          fmt.Println(output)
          os.Exit(eCode)
        }
      }

      // output status in the desired format 
      if !machineOutput {
        logger.Logger.Info("Cluster status is", "status", clusterStatus)
        logger.Logger.Info("detailedStatuses", statuses)
        os.Exit(0)
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.DirectoryCreated = true
        machineReadableOutput.ClusterStatus = clusterStatus
        machineReadableOutput.DetailedMachineStatus = statuses
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }

    // No part of the cluster exists, the directory is not created and machines are 
    // not present
    if !created {
      if !machineOutput {
        logger.Logger.Info("Cluster does not exist")
        os.Exit(0)
      } else {
        machineReadableOutput.ExitCode = 0
        machineReadableOutput.DirectoryCreated = false
        machineReadableOutput.ClusterStatus = "not created"
        output, eCode := machineReadableOutput.GetMachineOutputJson()
        fmt.Println(output)
        os.Exit(eCode)
      }
    }
  },
}

func init() {
  // required args for this command
  clusterStatusCmd.MarkFlagRequired("cluster")

  // add command
  RootCmd.AddCommand(clusterStatusCmd)
}


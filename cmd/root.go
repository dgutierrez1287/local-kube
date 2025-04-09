package cmd

import (
	"fmt"

  "github.com/dgutierrez1287/local-kube/logger"
  "github.com/dgutierrez1287/local-kube/util"

	"github.com/spf13/cobra"
)

var debug bool
var logColorize bool
var clusterName string
var machineOutput bool

var RootCmd = &cobra.Command{
  Use: "local-kube",
  Short: "A program to create and manage a local kube cluster using vagrant",
  Long: "A program to create and manage a local kube cluster using vagrant",
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    logger.InitLogging(debug, logColorize, machineOutput)
  }, 
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(util.TitleText)
    fmt.Println("local-kube, Use --help for help")
  },
}

func Execute() error {
  return RootCmd.Execute()
}

func init() {
  // cluster name
  RootCmd.PersistentFlags().StringVarP(&clusterName, "cluster", "c", "", "The cluster to run the action on")

  // debugging flag
  RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")

  // output colorization flag
  RootCmd.PersistentFlags().BoolVarP(&logColorize, "colorize", "", true, "Enable/Disable output colorization")

  // machine only output flag
  // This will suppess all other output and only output json that 
  // is machine readable
  RootCmd.PersistentFlags().BoolVarP(&machineOutput, "machine-output", "m", false, "Enables machine only output, json that can be used by executing script")
}



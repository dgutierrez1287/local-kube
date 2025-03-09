package cmd

import "github.com/spf13/cobra"


var clusterDownCmd = &cobra.Command {
  Use: "cluster-down",
  Short: "Brings a cluster down",
  Long: "Brings a cluster down if it is currently up",
  Run: func(cmd *cobra.Command, args []string) {

  },
}

func init() {
  // command specific args 

  // required args for this command
  clusterDownCmd.MarkFlagRequired("cluster")
}


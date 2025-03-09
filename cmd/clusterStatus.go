package cmd

import "github.com/spf13/cobra"

var clusterStatusCmd = &cobra.Command{
  Use: "cluster-status",
  Short: "Gets the status of a cluster",
  Long: "Gets the status of a cluster",
  Run: func(cmd *cobra.Command, args []string) {

  },
}

func init() {}


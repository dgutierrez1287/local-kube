package cmd

import "github.com/spf13/cobra"

var machineSshCmd = &cobra.Command {
  Use: "machine-ssh",
  Short: "Opens ssh session to machine",
  Long: "Opens menu to pick which machine to open ssh to",
  Run: func(cmd *cobra.Command, args []string) {

  },
}

func init() {}


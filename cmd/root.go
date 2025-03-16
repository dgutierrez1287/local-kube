package cmd

import (
	"fmt"

  "github.com/dgutierrez1287/local-kube/logger"
  "github.com/dgutierrez1287/local-kube/util"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

var debug bool
var logColorize bool

var RootCmd = &cobra.Command{
  Use: "local-kube",
  Short: "A program to create and manage a local kube cluster using vagrant",
  Long: "A program to create and manage a local kube cluster using vagrant",
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    setupLogging()
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
  // debugging flag
  RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")

  // output colorization flag
  RootCmd.PersistentFlags().BoolVarP(&logColorize, "colorize", "", true, "Enable/Disable output colorization")
}

func setupLogging() {
  // set up logger 

  var colorOpt hclog.ColorOption
  
  // debug output setup
  if debug {
    logger.LogLevel = "DEBUG"
    fmt.Println("Debugging enabled")
  } else {
    logger.LogLevel = "INFO"
  }

  // colorization setup
  if logColorize {
    colorOpt = hclog.ColorOption(hclog.AutoColor)
  } else {
    fmt.Println("Output colorization disabled")
    colorOpt = hclog.ColorOption(hclog.ColorOff)
  }

  fmt.Printf("loglevel %s \n", logger.LogLevel)

  // create global logger
  logger.Logger = hclog.New(&hclog.LoggerOptions{
    Name: "local-kube",
    Level: hclog.LevelFromString(logger.LogLevel),
    Color: colorOpt,
  })
}

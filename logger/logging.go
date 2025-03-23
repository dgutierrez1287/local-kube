package logger

import (
  "fmt"

  "github.com/hashicorp/go-hclog"
)

var Logger hclog.Logger
var LogLevel string

func InitLogging(debug bool, colorize bool) {
  // set up logger 

  var colorOpt hclog.ColorOption
  
  // debug output setup
  if debug {
    LogLevel = "DEBUG"
    fmt.Println("Debugging enabled")
  } else {
    LogLevel = "INFO"
  }

  // colorization setup
  if colorize {
    colorOpt = hclog.ColorOption(hclog.AutoColor)
  } else {
    fmt.Println("Output colorization disabled")
    colorOpt = hclog.ColorOption(hclog.ColorOff)
  }

  fmt.Printf("loglevel %s \n", LogLevel)

  // create global logger
  Logger = hclog.New(&hclog.LoggerOptions{
    Name: "local-kube",
    Level: hclog.LevelFromString(LogLevel),
    Color: colorOpt,
  })
}

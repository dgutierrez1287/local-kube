package logger

import (
	"fmt"
	"os"

	"github.com/dgutierrez1287/local-kube/output"
	"github.com/hashicorp/go-hclog"
)

var Logger hclog.Logger
var LogLevel string
var machineOutput bool

/*
Initialize logging, this will set level, colorization and 
if machine output is set
*/
func InitLogging(debug bool, colorize bool, machineOnlyOutput bool) {
  // set up logger 
  machineOutput = machineOnlyOutput
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
    if !machineOnlyOutput {
      fmt.Println("Output colorization disabled")
    }
    colorOpt = hclog.ColorOption(hclog.ColorOff)
  }

  if !machineOnlyOutput {
    fmt.Printf("loglevel %s \n", LogLevel)
  }

  // create global logger
  Logger = hclog.New(&hclog.LoggerOptions{
    Name: "local-kube",
    Level: hclog.LevelFromString(LogLevel),
    Color: colorOpt,
  })
}

/*
This will wrap info logging to handle if it should be 
written to console based on machine output setting
*/
func LogInfo(message string, args ...interface{}) {
  if !machineOutput {
    Logger.Info(message, args...)
  }
}

/*
This will wrap debug logging to handle any special 
caes
*/
func LogDebug(message string, args ...interface{}) {
  Logger.Debug(message, args...)
}

/*
This will wrap error logging (where there is no exit)
and will handle based on the machine output setting
*/
func LogError(message string, args ...interface{}) {
  if !machineOutput {
    Logger.Error(message)
  }
}

/*
This will wrap error logging with an exit and setting 
exit code with it and will act based on machine output 
setting
*/
func LogErrorExit(meesage string, exitCode int, err error) {
  var machineReadableOutput output.MachineOutput

  if !machineOutput {
    Logger.Error(meesage, "error", err)
    os.Exit(exitCode)
  } else {
    machineReadableOutput.ExitCode = exitCode
    machineReadableOutput.ErrorMessage = fmt.Sprintf("%s: %v", meesage, err)
    output, _ := machineReadableOutput.GetMachineOutputJson()
    fmt.Println(output)
    os.Exit(exitCode)
  }
}

package output

import (
  "encoding/json"
)

/*
  StatusMachineOutput - The json structure for machine
  readable output of the status command
*/
type MachineOutput struct {
  ExitCode int                              `json:"exitCode"`
  ErrorMessage string                       `json:"errorMessage,omitempty"`
  StatusMessage string                      `json:"statusMessage,omitempty"`
  DirectoryCreated bool                     `json:"directoryCreated,omitempty"`
  ClusterStatus string                      `json:"clusterStatus,omitempty"`
  DetailedMachineStatus map[string]string   `json:"machineStatus,omitempty"`
}

/*
This will get the json string of machine readable output
*/
func (status MachineOutput) GetMachineOutputJson() (string, int){
  
  jsonBytes, err := json.Marshal(status)
  if err != nil {
    return "{\"exitCode\": 100, \"errorMessage\": \"Error unmarshaling machine output\"}", 100
  }
  return string(jsonBytes), 0
}

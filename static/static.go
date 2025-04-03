package static

import (
	"embed"
	"io/fs"

	"github.com/dgutierrez1287/local-kube/logger"
)

//embed directory for remote scripts
//(scripts that run on the master node to provision
//go:embed remote_scripts/*
var remoteScriptsFs embed.FS

//embed directory for provision scripts
//scripts that are present on every node that 
//are used for provision tasks during vagrant up
//go:embed provision_scripts/*
var provisionScriptsFs embed.FS

func GetProvisionScriptsFS() fs.FS {
  return provisionScriptsFs
}

func GetRemoteScriptsFS() fs.FS {
  return remoteScriptsFs
}

/*
  This will return an array of all the files in the 
  provision scripts embeded directory
*/
func ListProvisonScripts() ([]string, error) {
  fileNames := []string{}

  files, err := fs.ReadDir(provisionScriptsFs, "provision_scripts")
  if err != nil {
    logger.Logger.Error("Error reading provision_scripts directory")
    return fileNames, err
  }

  for _, file := range files {
    logger.Logger.Debug("Adding script file to list", "file", file.Name())
    fileNames = append(fileNames, file.Name())
  }
  return fileNames, nil
}

/*
  This will return an array of all the file in the
  remote scripts embeded directory
*/
func ListRemoteScripts() ([]string, error) {
  fileNames := []string{}

  files, err := fs.ReadDir(remoteScriptsFs, "remote_scripts")
  if err != nil {
    logger.Logger.Error("Error reading remote_scripts directory")
    return fileNames, err
  }

  for _, file := range files {
    logger.Logger.Debug("Adding script file to list", "file", file.Name())
    fileNames = append(fileNames, file.Name())
  }
  return fileNames, nil
}

/*
  This will return a string of the provision script 
  file content
*/
func ReadProvisionScriptFile(scriptName string) (string, error) {
  path := "provision_scripts/" + scriptName

  data, err := provisionScriptsFs.ReadFile(path)
  if err != nil {
    logger.Logger.Error("Error reading static provision script", "path", path)
    return "", err
  }
  return string(data), nil
}

/*
  This will return a string of the remote script 
  file content
*/
func ReadRemoteScriptFile(scriptName string) (string, error){
  path := "remote_scripts/" + scriptName

  data, err := remoteScriptsFs.ReadFile(path)
  if err != nil {
    logger.Logger.Error("Error reading static remote script", "path", path)
    return "", err
  }
  return string(data), nil
}



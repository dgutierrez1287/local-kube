package static

import (
	"embed"

	"github.com/dgutierrez1287/local-kube/logger"
)

//embed directory for remote scripts
//(scripts that run on the master node to provision
//go:embed remote_scripts
var rsScripts embed.FS

//embed directory for provision scripts
//scripts that are present on every node that 
//are used for provision tasks during vagrant up
var psScripts embed.FS


func ReadProvisionScriptFile(path string) (string, error) {
  data, err := psScripts.ReadFile(path)
  if err != nil {
    logger.Logger.Error("Error reading static provision script", "path", path)
    return "", err
  }
  return string(data), nil
}

func ReadRemoteScriptFile(path string) (string, error){
  data, err := rsScripts.ReadFile(path)
  if err != nil {
    logger.Logger.Error("Error reading static remote script", "path", path)
    return "", err
  }
  return string(data), nil
}



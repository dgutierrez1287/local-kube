package ansible

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"gopkg.in/yaml.v3"
)

/*
  Playbook - Ansible playbook type
*/
type Playbook struct {
  Name string         `json:"name" yaml:"name"`
  Hosts string        `json:"hosts" yaml:"hosts"`
  Become bool         `json:"become" yaml:"become"`
  BecomeUser string   `json:"become_user" yaml:"become_user"`
  VarsFiles []string  `json:"vars_files" yaml:"vars_files"`
  Roles []string      `json:"roles" yaml:"roles"`
}

/*
RenderPlaybook - This will render a playbook for either a single node cluster or 
a playbook for a node type in an HA cluster
*/
func RenderPlaybook(appDir string, clusterName string, hosts string, clusterType string, nodeType string) error {

  var playName string
  var dynamicVarsFileName string
  var playbookFileName string

  if clusterType == "ha" {
    logger.Logger.Debug("cluster is an ha cluster writing playbook for node type", "nodeType", nodeType)

    playName = fmt.Sprintf("multi node %s node playbook", nodeType)
    dynamicVarsFileName = fmt.Sprintf("vars-dynamic-%s.yml", nodeType)
    playbookFileName = fmt.Sprintf("%s-playbook.yml", nodeType)
  } else {
    logger.Logger.Debug("cluster is a single node writing playbook for single node cluster")

    playName = "single node cluster playbook"
    dynamicVarsFileName = "vars-dynamic.yml"
    playbookFileName = "playbook.yml"
  }

  logger.Logger.Debug("Generating playbook data")
  var playbookData = Playbook{
    Name: playName,
    Hosts: hosts,
    Become: true,
    BecomeUser: "root",
    VarsFiles: []string {
      "/etc/ansible/vars/static/vars-static.yml",
      fmt.Sprintf("/etc/ansible/vars/dynamic/%s", dynamicVarsFileName),
    },
    Roles: []string {
      "kube",
    },
  }

  yamlData, err := yaml.Marshal(&playbookData)
  if err != nil {
    logger.Logger.Error("Error marshaling playbook data to yaml")
    return err
  }

  path := filepath.Join(appDir, clusterName, "playbooks", playbookFileName)
  err = os.WriteFile(path, yamlData, 0755)
  if err != nil {
    logger.Logger.Error("Error writing playbook data to file")
    return err
  }
  return nil
}


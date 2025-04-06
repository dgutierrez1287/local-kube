package ansible

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"gopkg.in/yaml.v3"
)

/*
  Play - Ansible play type
*/
type Play struct {
  Name string         `json:"name" yaml:"name"`
  Hosts string        `json:"hosts" yaml:"hosts"`
  Become bool         `json:"become" yaml:"become"`
  BecomeUser string   `json:"become_user" yaml:"become_user"`
  VarsFiles []string  `json:"vars_files" yaml:"vars_files"`
  Roles []string      `json:"roles" yaml:"roles"`
}

/*
  PlayBook - An Ansible playbook (consists of one or multiple plays)
*/
type Playbook []Play

/*
RenderPlaybook - This will render a playbook for either a single node cluster or 
a playbook for a node type in an HA cluster
*/
func RenderPlaybook(appDir string, clusterName string, hosts string, clusterType string, nodeType string) error {

  var playName string
  var varsFileName string
  var playbookFileName string

  var playbook Playbook

  if clusterType == "ha" {
    logger.Logger.Debug("cluster is an ha cluster writing playbook for node type", "nodeType", nodeType)

    playName = fmt.Sprintf("multi node %s node playbook", nodeType)
    varsFileName = fmt.Sprintf("vars-%s.yml", nodeType)
    playbookFileName = fmt.Sprintf("%s-playbook.yml", nodeType)
  } else {
    logger.Logger.Debug("cluster is a single node writing playbook for single node cluster")

    playName = "single node cluster playbook"
    varsFileName = "vars.yml"
    playbookFileName = "playbook.yml"
  }

  logger.Logger.Debug("Generating playbook data")
  var playData = Play{
    Name: playName,
    Hosts: hosts,
    Become: true,
    BecomeUser: "root",
    VarsFiles: []string {
      fmt.Sprintf("/etc/ansible/vars/%s", varsFileName),
    },
    Roles: []string {
      "kube",
    },
  }

  playbook = append(playbook, playData)

  yamlData, err := yaml.Marshal(&playbook)
  if err != nil {
    logger.Logger.Error("Error marshaling playbook data to yaml")
    return err
  }

  path := filepath.Join(appDir, clusterName, "ansible" ,"playbooks", playbookFileName)
  file, err := os.Create(path)

  if err != nil {
    logger.Logger.Error("Error creating the playbook file", "name", playbookFileName)
    return err
  }

  defer file.Close()

  _, err = file.WriteString("---\n")
  if err != nil {
    logger.Logger.Error("Error writing the yaml doc start marker")
    return err
  }

  _, err = file.Write(yamlData)
  if err != nil {
    logger.Logger.Error("Error writing playbook content to file")
    return err
  }
  return nil
}


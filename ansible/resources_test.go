package ansible

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dgutierrez1287/local-kube/template"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
)

func TestEmbeddedFile(t *testing.T) {
  fs := template.GetProvisionFS()

  _, err := fs.Open("provision/bootstrap.tmpl")
  assert.NoError(t, err)
}


/*
      Tests for GenerateAnsibleHostsFile
*/
func TestGenerateAnsibleHostsFileMultiNode(t *testing.T) {
  clusterName := "test"
  clusterType := "ha"
  secondaryControlNodes := []string{
    "cp2",
    "cp3",
  }
  workerNodes := []string{
    "worker1",
    "worker2",
  }
  filePath := filepath.Join(util.MockAppDir, clusterName, "ansible", "resources", "hosts")

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "ansible", "resources"), 0755)
  assert.NoError(t, err)

  err = GenerateAnsibleHostsFile(util.MockAppDir, clusterName, clusterType, secondaryControlNodes, workerNodes)
  assert.NoError(t, err)

  // verify the file was created
  assert.FileExists(t, filePath)

  // start verification of content
  // triggers to make sure each group exists
  leadNodeFound := false
  controlNodesFound := false
  workerNodesFound := false

  // arrays to compare control and worker ips in groups
  leadNodeLine := ""
  controlNodesTestArray := []string{}
  workerNodesTestArray := []string{}

  fileLines, err := util.ReadFileToStringArray(filePath)
  assert.NoError(t, err)

  assert.Len(t, fileLines, 10)

  //Loop through the file lines an set group triggers and verify content
  for index, line := range fileLines {
    switch line {
    case "[lead-node]":
      leadNodeFound = true //set trigger to true since its found 
      leadNodeLine = fileLines[index + 1] // get the next line for verification
    case "[control-nodes]":
      controlNodesFound = true // set trigger to true
      // get the next two lines for verification
      controlNodesTestArray = append(controlNodesTestArray, fileLines[index + 1])
      controlNodesTestArray = append(controlNodesTestArray, fileLines[index + 2])
    case "[worker-nodes]":
      workerNodesFound = true
      // get the next two lines for verifcation
      workerNodesTestArray = append(workerNodesTestArray, fileLines[index + 1])
      workerNodesTestArray = append(workerNodesTestArray, fileLines[index + 2])
    default: 
      //skip
    }
  }

  // all triggers should be true
  assert.True(t, leadNodeFound)
  assert.True(t, controlNodesFound)
  assert.True(t, workerNodesFound)

  // verify lines
  assert.Contains(t, leadNodeLine, "localhost ansible_connection=local")

  secondaryControlVerifyArray := []string {
    "cp2 ansible_connection=ssh ansible_user=vagrant",
    "cp3 ansible_connection=ssh ansible_user=vagrant",
  }

  workerNodeVerifyArray := []string {
    "worker1 ansible_connection=ssh ansible_user=vagrant",
    "worker2 ansible_connection=ssh ansible_user=vagrant",
  }

  assert.Equal(t, controlNodesTestArray, secondaryControlVerifyArray)
  assert.Equal(t, workerNodesTestArray, workerNodeVerifyArray)
}

func TestGenerateAnsibleHostsFileSingleNode(t *testing.T) {
  clusterName := "test"
  clusterType := "single"
  secondaryControlNodes := []string{}
  workerNodes := []string{}
  filePath := filepath.Join(util.MockAppDir, clusterName, "ansible", "resources", "hosts")

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "ansible", "resources"), 0755)
  assert.NoError(t, err)

  err = GenerateAnsibleHostsFile(util.MockAppDir, clusterName, clusterType, secondaryControlNodes, workerNodes)
  assert.NoError(t, err)

  // verify the file was created 
  assert.FileExists(t, filePath)

  // start verification of the content
  fileLines, err := util.ReadFileToStringArray(filePath) // read generated file to array of strings
  assert.NoError(t, err)

  assert.Len(t, fileLines, 1) // file should only have one line
  line := fileLines[0] // read that one line
  assert.Contains(t, line, "localhost ansible_connection=local") // validate the line for content
}

func TestGenerateAnsibleHostsFileError(t *testing.T) {
  clusterName := "test"
  clusterType := "single"
  secondaryControlNodes := []string{}
  workerNodes := []string{}
  
  err := GenerateAnsibleHostsFile(util.MockAppDir, clusterName, clusterType, secondaryControlNodes, workerNodes)

  assert.Error(t, err)
}

/*
      Tests for RenderBootstrapScript
*/
func TestRenderBootstrapScript(t *testing.T) {
  clusterName := "test"
  ansibleVersion := "2.17.5"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "scripts", "provision"), 0755)
  assert.NoError(t, err)

  err = RenderBootstrapScript(util.MockAppDir, clusterName, ansibleVersion)
  assert.NoError(t, err)

  //verification
  boostrapScriptPath := filepath.Join(util.MockAppDir, clusterName, "scripts", "provision", "bootstrap.sh")
  
  fileLines, err := util.ReadFileToStringArray(boostrapScriptPath)
  assert.NoError(t, err)
  verificationLine := ""

  for _, line := range fileLines {
    if strings.Contains(line, "ansible_version=") {
      verificationLine = line
    }
  }

  assert.Equal(t, verificationLine, "ansible_version=2.17.5")
}

func TestRenderBootstrapScriptErrorWrite(t *testing.T) {
  clusterName := "test"
  ansibleVersion := "2.17.5"

  err := RenderBootstrapScript(util.MockAppDir, clusterName, ansibleVersion)
  assert.Error(t, err)
}

/*
      Tests for Copy Ansible Roles
*/
func CopyAnsibleRolesSuccess(t *testing.T) {
  clusterName := "test"
  roleNames := []string{"kube"}

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  // create directories for tests
  err = os.MkdirAll(filepath.Join(util.MockAppDir, "ansible-roles", "kube"), 0755)
  assert.NoError(t, err)

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "ansible", "roles"), 0755)
  assert.NoError(t, err)

  //Run test
  err = CopyAnsibleRoles(util.MockAppDir, clusterName, roleNames)
  assert.NoError(t, err)
  assert.FileExists(t, filepath.Join(util.MockAppDir, clusterName, "ansible", "roles", "kube"))
}

func CopyAnsibleRolesError(t *testing.T) {
  clusterName := "test"
  roleNames := []string{"kube", "doesNotExist"}

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  // create directories for tests
  err = os.MkdirAll(filepath.Join(util.MockAppDir, "ansible-roles", "kube"), 0755)
  assert.NoError(t, err)

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "ansible", "roles"), 0755)
  assert.NoError(t, err)

  //Run test
  err = CopyAnsibleRoles(util.MockAppDir, clusterName, roleNames)
  assert.Error(t, err)
}


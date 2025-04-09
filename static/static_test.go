package static

import (
	"os"
	"strings"
	"testing"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/stretchr/testify/assert"
)

// TestMain is executed before running any tests
func TestMain(m *testing.M) {
	// Initialize the logger before running any tests
	logger.InitLogging(false, true, false)
	os.Exit(m.Run())
}

func TestListProvisionScripts(t *testing.T) {
  fileNames, err := ListProvisonScripts()

  assert.NoError(t, err)

  // verify file names in the list
  assert.Contains(t, fileNames, "install-yq.sh")
  assert.Contains(t, fileNames, "disk-expand.sh")
  assert.Contains(t, fileNames, "setup-hostsfile.sh")
  assert.Contains(t, fileNames, "multi-node-provision.sh")
}

func TestListRemoteScripts(t *testing.T) {
  fileNames, err := ListRemoteScripts()

  assert.NoError(t, err)

  // verify file names in the list
  assert.Contains(t, fileNames, "ha-provision.sh")
  assert.Contains(t, fileNames, "single-provision.sh")
  assert.Contains(t, fileNames, "correct_kubeconfig.py")
}

func TestReadProvisionScriptFile(t *testing.T) {
  fileContent, err := ReadProvisionScriptFile("disk-expand.sh")

  assert.NoError(t, err)
  
  // verify first line is the shebang
  firstLine := strings.Split(fileContent, "\n")[0]
  assert.Equal(t, firstLine, "#!/usr/bin/env bash")
}

func TestReadProvisionScriptError(t *testing.T) {
  _, err := ReadProvisionScriptFile("does-not-exist.sh")

  // should return an error
  assert.Error(t, err)
}

func TestReadRemoteScriptFile(t *testing.T) {
  fileContent, err := ReadRemoteScriptFile("single-provision.sh")

  assert.NoError(t, err)

  // verify first line is the shebang
  firstLine := strings.Split(fileContent, "\n")[0]
  assert.Equal(t, firstLine, "#!/usr/bin/env bash")
}

func TestReadRemoteScriptError(t *testing.T) {
  _, err := ReadRemoteScriptFile("does-not-exist.sh")

  //should be an error 
  assert.Error(t, err)
}

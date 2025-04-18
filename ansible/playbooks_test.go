package ansible

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

/*
      Tests for RenderPlaybook
*/

func TestRenderPlaybookWorkerNode(t *testing.T) {
  clusterName := "test"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "ansible", "playbooks"), 0755)
  assert.NoError(t, err)

  err = RenderPlaybook(util.MockAppDir, clusterName, "worker-nodes", "ha", "worker")
  assert.NoError(t, err)
  assert.FileExists(t, filepath.Join(util.MockAppDir, clusterName, "ansible", "playbooks", "worker-playbook.yml"))

  yamlFile, err := os.ReadFile(filepath.Join(util.MockAppDir, clusterName, "ansible", "playbooks", "worker-playbook.yml"))
  assert.NoError(t, err)

  var resultPlaybook Playbook
  err = yaml.Unmarshal(yamlFile, &resultPlaybook)
  play := resultPlaybook[0]
  assert.NoError(t, err)

  // verify result
  resultVarsFiles := []string {
    "/etc/ansible/vars/vars-worker.yml",
  }

  resultPlayName := "multi node worker node playbook"

  assert.Equal(t, play.Hosts, "worker-nodes")
  assert.Equal(t, play.VarsFiles, resultVarsFiles)
  assert.Equal(t, play.Name, resultPlayName)
  assert.True(t, play.Become)
  assert.Equal(t, play.BecomeUser, "root")
}

func TestRenderPlaybookSingleNodeCluster(t *testing.T) {
  clusterName := "test"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "ansible", "playbooks"), 0755)
  assert.NoError(t, err)

  err = RenderPlaybook(util.MockAppDir, clusterName, "localhost", "single", "")
  assert.NoError(t, err)
  assert.FileExists(t, filepath.Join(util.MockAppDir, clusterName, "ansible", "playbooks", "playbook.yml"))

  yamlFile, err := os.ReadFile(filepath.Join(util.MockAppDir, clusterName, "ansible", "playbooks", "playbook.yml"))
  assert.NoError(t, err)

  var resultPlaybook Playbook
  err = yaml.Unmarshal(yamlFile, &resultPlaybook)
  play := resultPlaybook[0]
  assert.NoError(t, err)

  // verify result
  resultVarsFiles := []string {
    "/etc/ansible/vars/vars.yml",
  }

  resultPlayName := "single node cluster playbook"

  assert.Equal(t, play.Hosts, "localhost")
  assert.Equal(t, play.VarsFiles, resultVarsFiles)
  assert.Equal(t, play.Name, resultPlayName)
  assert.True(t, play.Become)
  assert.Equal(t, play.BecomeUser, "root")
}

func TestRenderPlaybookSingleNodeError(t *testing.T) {
  clusterName := "test"

  err := RenderPlaybook(util.MockAppDir, clusterName, "localhost", "single", "")
  assert.Error(t, err)
}

package ansible

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestRenderPlaybookWorkerNode(t *testing.T) {
  clusterName := "test"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "playbooks"), 0755)
  assert.NoError(t, err)

  err = RenderPlaybook(util.MockAppDir, clusterName, "worker-nodes", "ha", "worker")
  assert.NoError(t, err)
  assert.FileExists(t, filepath.Join(util.MockAppDir, clusterName, "playbooks", "worker-playbook.yml"))

  yamlFile, err := os.ReadFile(filepath.Join(util.MockAppDir, clusterName, "playbooks", "worker-playbook.yml"))
  assert.NoError(t, err)

  var resultPlaybook Playbook
  err = yaml.Unmarshal(yamlFile, &resultPlaybook)
  assert.NoError(t, err)

  // verify result
  resultVarsFiles := []string {
    "/etc/ansible/vars/static/vars-static.yml",
    "/etc/ansible/vars/dynamic/vars-dynamic-worker.yml",
  }

  resultPlayName := "multi node worker node playbook"

  assert.Equal(t, resultPlaybook.Hosts, "worker-nodes")
  assert.Equal(t, resultPlaybook.VarsFiles, resultVarsFiles)
  assert.Equal(t, resultPlaybook.Name, resultPlayName)
  assert.True(t, resultPlaybook.Become)
  assert.Equal(t, resultPlaybook.BecomeUser, "root")
}

func TestRenderPlaybookSingleNodeCluster(t *testing.T) {
  clusterName := "test"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.MkdirAll(filepath.Join(util.MockAppDir, clusterName, "playbooks"), 0755)
  assert.NoError(t, err)

  err = RenderPlaybook(util.MockAppDir, clusterName, "localhost", "single", "")
  assert.NoError(t, err)
  assert.FileExists(t, filepath.Join(util.MockAppDir, clusterName, "playbooks", "playbook.yml"))

  yamlFile, err := os.ReadFile(filepath.Join(util.MockAppDir, clusterName, "playbooks", "playbook.yml"))
  assert.NoError(t, err)

  var resultPlaybook Playbook
  err = yaml.Unmarshal(yamlFile, &resultPlaybook)
  assert.NoError(t, err)

  // verify result
  resultVarsFiles := []string {
    "/etc/ansible/vars/static/vars-static.yml",
    "/etc/ansible/vars/dynamic/vars-dynamic.yml",
  }

  resultPlayName := "single node cluster playbook"

  assert.Equal(t, resultPlaybook.Hosts, "localhost")
  assert.Equal(t, resultPlaybook.VarsFiles, resultVarsFiles)
  assert.Equal(t, resultPlaybook.Name, resultPlayName)
  assert.True(t, resultPlaybook.Become)
  assert.Equal(t, resultPlaybook.BecomeUser, "root")
}

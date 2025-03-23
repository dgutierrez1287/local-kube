package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
)

func TestSettingsValid(t *testing.T) {

  clusterName := "test-cluster"

  clusters := make(map[string]Cluster)
  clusters[clusterName] = Cluster{}

  settings := Settings{
    ProvisionSettings: ProvisionSettings{},
    Providers: map[string]Provider{},
    Clusters: clusters,
  }

  valid := settings.SettingsValid(clusterName)

  assert.True(t, valid)
}

func TestSettingsNotValid(t *testing.T) {

  clusterName := "test-cluster"

  clusters := make(map[string]Cluster)
  clusters["test-cluster4"] = Cluster{}

  settings := Settings{
    ProvisionSettings: ProvisionSettings{},
    Providers: map[string]Provider{},
    Clusters: clusters,
  }

  notValid := settings.SettingsValid(clusterName)

  assert.False(t, notValid)
}

func TestCreateDefaultSettingsAndRead(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  err = CreateDefaultSettingsFile(util.MockAppDir)
  assert.NoError(t, err)

  settings, err := ReadSettingsFile(util.MockAppDir)
  assert.NoError(t, err)

  assert.Contains(t, settings.ProvisionSettings.AnsibleRoles, "kube")
  assert.Equal(t, settings.ProvisionSettings.AnsibleRoles["kube"].LocationType, "git")
  assert.Equal(t, settings.ProvisionSettings.AnsibleRoles["kube"].RefType, "branch")
  assert.Equal(t, settings.ProvisionSettings.AnsibleRoles["kube"].GitRef, "master")
  assert.Equal(t, settings.ProvisionSettings.AnsibleVersion, "2.17.6")

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestSettingsExists(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  file, err := os.Create(filepath.Join(util.MockAppDir, "settings.json"))
  assert.NoError(t, err)
  
  file.Close()

  exists, err := SettingsFileExists(util.MockAppDir)
  assert.NoError(t, err)

  assert.True(t, exists)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestSettingsNotExist(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  notExists, err := SettingsFileExists(util.MockAppDir)
  assert.NoError(t, err)
  
  assert.False(t, notExists)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

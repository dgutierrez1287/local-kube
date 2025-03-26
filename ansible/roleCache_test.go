package ansible

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgutierrez1287/local-kube/util"
  "github.com/dgutierrez1287/local-kube/settings"
	"github.com/stretchr/testify/assert"
)

func TestRoleCacheExists(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  file, err := os.Create(filepath.Join(util.MockAnsibleRoleDir, ".role-cache.json"))
  assert.NoError(t, err)
  
  file.Close()

  exists, err := RoleCacheFileExists(util.MockAppDir)
  assert.NoError(t, err)

  assert.True(t, exists)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestRoleCacheNotExist(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  notExists, err := RoleCacheFileExists(util.MockAppDir)
  assert.NoError(t, err)

  assert.False(t, notExists)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestDeleteRoleCache(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  file, err := os.Create(filepath.Join(util.MockAnsibleRoleDir, ".role-cache.json"))
  assert.NoError(t, err)

  file.Close()

  err = RoleCacheFileDelete(util.MockAppDir)
  assert.NoError(t, err)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestWriteAndReadRoleCache(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  cacheToWrite := RoleCache{
    Roles: make(map[string]settings.AnsibleRole),
  }

  cacheToWrite.Roles["kube"] = settings.AnsibleRole{
    LocationType: "git",
    Location: "https://github.com/dgutierrez1287/ansible-role-kube",
    RefType: "branch",
    GitRef: "master",
  }

  err = WriteRoleCache(util.MockAppDir, cacheToWrite)
  assert.NoError(t, err)

  readCache, err := ReadRoleCache(util.MockAppDir)
  assert.NoError(t, err)

  assert.Contains(t, readCache.Roles, "kube")
  assert.Equal(t, readCache.Roles["kube"].LocationType, "git")
  assert.Equal(t, readCache.Roles["kube"].Location, "https://github.com/dgutierrez1287/ansible-role-kube")
  assert.Equal(t, readCache.Roles["kube"].RefType, "branch")
  assert.Equal(t, readCache.Roles["kube"].GitRef, "master")

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}


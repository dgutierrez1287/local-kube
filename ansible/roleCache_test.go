package ansible

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
)

/*
   Tests for RoleCacheFileExists
*/
func TestRoleCacheExists(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  file, err := os.Create(filepath.Join(util.MockAnsibleRoleDir, ".role-cache.json"))
  assert.NoError(t, err)
  
  file.Close()

  exists, err := RoleCacheFileExists(util.MockAppDir)
  assert.NoError(t, err)

  assert.True(t, exists)
}

func TestRoleCacheNotExist(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  notExists, err := RoleCacheFileExists(util.MockAppDir)
  assert.NoError(t, err)

  assert.False(t, notExists)
}


/*
      Tests for RoleCacheFileDelete
*/
func TestDeleteRoleCache(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  file, err := os.Create(filepath.Join(util.MockAnsibleRoleDir, ".role-cache.json"))
  assert.NoError(t, err)

  file.Close()

  err = RoleCacheFileDelete(util.MockAppDir)
  assert.NoError(t, err)
}

func TestDeleteRoleCacheError(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = RoleCacheFileDelete(util.MockAppDir)
  assert.Error(t, err)
}


/*
      Tests for ReadRoleCache and WriteRoleCache 
      Success
*/
func TestWriteAndReadRoleCache(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

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
}

/*
      Tests for ReadRoleCache Failures
*/
func TestReadRoleCacheReadError(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  _, err = ReadRoleCache(util.MockAppDir)
  assert.Error(t, err)
}

func TestReadRoleCacheUnmarshalError(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  badJson := `
    "roles": {
      kube: {
        "locationType"; local
      }
    }
  `

  err = os.WriteFile(filepath.Join(util.MockAnsibleRoleDir, ".role-cache.json"), []byte(badJson), 0755)
  assert.NoError(t, err)

  _, err = ReadRoleCache(util.MockAppDir)
  assert.Error(t, err)
}

/*
        Tests for WriteRoleCache Failures
*/
func TestWriteRoleCacheError(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  err = os.Remove(filepath.Join(util.MockAppDir, "ansible-roles"))
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
  assert.Error(t, err)
}


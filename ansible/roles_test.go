package ansible

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
)

// TestMain is executed before running any tests
func TestMain(m *testing.M) {
	// Initialize the logger before running any tests
	logger.InitLogging(false, true)
	os.Exit(m.Run())
}

func TestClearRoles(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

	err = os.Mkdir(filepath.Join(util.MockAnsibleRoleDir, "test-role"), 0755)
	assert.NoError(t, err)

	err = ClearRoles(util.MockAppDir)
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(util.MockAnsibleRoleDir, "test-role"))
	assert.True(t, os.IsNotExist(err))

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestClearRole(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

	rolePath := filepath.Join(util.MockAnsibleRoleDir, "test-role")

	err = os.Mkdir(rolePath, 0755)
	assert.NoError(t, err)

	err = ClearRole(util.MockAppDir, "test-role")
	assert.NoError(t, err)

	_, err = os.Stat(rolePath)
	assert.True(t, os.IsNotExist(err))

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestRoleReconcileLists(t *testing.T) {
  currentRoles := map[string]settings.AnsibleRole{
		"role1": {LocationType: "git", Location: "repo1"},
		"role2": {LocationType: "local", Location: "/local/role2"},
    "role4": {LocationType: "local", Location: "/local/role4"},
	}
	desiredRoles := map[string]settings.AnsibleRole{
		"role1": {LocationType: "git", Location: "repo1"},
		"role3": {LocationType: "git", Location: "repo3"},
    "role4": {LocationType: "local", Location: "/local/role4"},
	}

	rolesToAdd, rolesToUpdate, rolesToReAdd, rolesToRemove, err := RoleReconcileLists(currentRoles, desiredRoles)
	assert.NoError(t, err)
  assert.Contains(t, rolesToReAdd, "role4")
	assert.Contains(t, rolesToAdd, "role3")
	assert.Contains(t, rolesToRemove, "role2")
	assert.Contains(t, rolesToUpdate, "role1")
}

func TestGitUpdateableRole(t *testing.T) {
	currentRole := settings.AnsibleRole{LocationType: "git", Location: "repo1"}
	newRole := settings.AnsibleRole{LocationType: "git", Location: "repo1"}

	assert.True(t, gitUpdateableRole(currentRole, newRole))

	newRole.Location = "repo2"
	assert.False(t, gitUpdateableRole(currentRole, newRole))
}

func TestCreateGitRole(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  roleName := "kube"
  roleData := settings.AnsibleRole{
    LocationType: "git", 
    Location: "https://github.com/dgutierrez1287/ansible-role-kube",
    RefType: "branch",
    GitRef: "master",
  }

  err = CreateGitRole(util.MockAppDir, roleName, roleData)
  assert.NoError(t, err)

  assert.DirExists(t, filepath.Join(util.MockAnsibleRoleDir, "kube"))
  assert.DirExists(t, filepath.Join(util.MockAnsibleRoleDir, "kube", ".git"))

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestUpdateGitRole(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  roleName := "kube"
  roleData := settings.AnsibleRole{
    LocationType: "git", 
    Location: "https://github.com/dgutierrez1287/ansible-role-kube",
    RefType: "branch",
    GitRef: "master",
  }

  err = CreateGitRole(util.MockAppDir, roleName, roleData)
  assert.NoError(t, err)

  assert.DirExists(t, filepath.Join(util.MockAnsibleRoleDir, "kube"))
  assert.DirExists(t, filepath.Join(util.MockAnsibleRoleDir, "kube", ".git"))

  err = UpdateGitRole(util.MockAppDir, roleName, roleData, roleData)
  assert.NoError(t, err)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestCreateLocalRole(t *testing.T) {
  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  roleName := "kube"
  roleData := settings.AnsibleRole{
    LocationType: "local",
    Location: "./mock/source/kube",
  }

  err = os.MkdirAll("./mock/source/kube/tasks", 0755)
  assert.NoError(t, err)
  
  err = CreateLocalRole(util.MockAppDir, roleName, roleData)
  assert.NoError(t, err)

  assert.DirExists(t, filepath.Join(util.MockAnsibleRoleDir, "kube"))
  assert.DirExists(t, filepath.Join(util.MockAnsibleRoleDir, "kube", "tasks"))

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestGetReferenceNameBranch(t *testing.T) {
	ref, spec, err := getReferenceName("branch", "main")
	assert.NoError(t, err)
	assert.Equal(t, "refs/heads/main", string(ref))
	assert.Equal(t, "+refs/heads/main:refs/remotes/origin/main", string(spec))
}

func TestGetReferenceNameTag(t *testing.T) {
  ref, spec, err := getReferenceName("tag", "v1.0.0")
  assert.NoError(t, err)
  assert.Equal(t, "refs/tags/v1.0.0", string(ref))
  assert.Equal(t, "+refs/tags/v1.0.0:refs/tags/v1.0.0", string(spec))
}


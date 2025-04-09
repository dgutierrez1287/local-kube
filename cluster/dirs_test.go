package cluster

import (
	"os"
	"path/filepath"
	"testing"

  "github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
)

// TestMain is executed before running any tests
func TestMain(m *testing.M) {
	// Initialize the logger before running any tests
	logger.InitLogging(false, true, false)
	os.Exit(m.Run())
}

func TestClusterDirExists(t *testing.T) {
  clusterName := "test-cluster"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  err = os.Mkdir(filepath.Join(util.MockAppDir, clusterName), 0755)
  assert.NoError(t, err)

  exists, err := ClusterDirExists(util.MockAppDir, clusterName)
  assert.NoError(t, err)
  assert.True(t, exists)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestClusterDirNotExists(t *testing.T) {
  clusterName := "test-cluster"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  notExists, err := ClusterDirExists(util.MockAppDir, clusterName)
  assert.NoError(t, err)
  assert.False(t, notExists)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestCreateClusterDirs(t *testing.T) {
  clusterName := "test-cluster"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  err = CreateClusterDirs(util.MockAppDir, clusterName)
  assert.NoError(t, err)

  // check main directories
  for _, dir := range mainDirs {
    dirPath := filepath.Join(util.MockAppDir, clusterName, dir)
    assert.DirExists(t, dirPath)
  }

  // check ansible directories 
  for _, dir := range ansibleDirs {
    dirPath := filepath.Join(util.MockAppDir, clusterName, "ansible", dir)
    assert.DirExists(t, dirPath)
  }

  // check script directories
  for _, dir := range scriptsDirs {
    dirPath := filepath.Join(util.MockAppDir, clusterName, "scripts", dir)
    assert.DirExists(t, dirPath)
  }

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

func TestDeleteClusterDir(t *testing.T) {
  clusterName := "test-cluster"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  err = os.Mkdir(filepath.Join(util.MockAppDir, clusterName), 0755)
  assert.NoError(t, err)

  err = DeleteClusterDir(util.MockAppDir, clusterName)
  assert.NoError(t, err)

  err = util.MockAppDirCleanup()
  assert.NoError(t, err)
}

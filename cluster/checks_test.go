package cluster

import (
	"testing"

	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
)

func TestCheckForExistingClusterNoDir(t *testing.T) {
  mockStatus := map[string]string{}
  clusterName := "test-cluster"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  vagrantClientMock := NewMockVagrantClientStatus(mockStatus)

  status, state, err := CheckForExistingCluster(util.MockAppDir, clusterName, vagrantClientMock, false)
  
  assert.False(t, status)
  assert.Equal(t, state, "")
  assert.NoError(t, err)
}

// func TestCheckForExistingClusterDirOnly(t *testing.T) {
//   mockStatus := map[string]string{}
//   clusterName := "test-cluster"
//
//   err := util.MockAppDirSetup()
//   assert.NoError(t, err)
//
//   err = os.Mkdir(filepath.Join(util.MockAppDir, clusterName), 0755)
//   assert.NoError(t, err)
//
//   defer util.MockAppDirCleanup()
//
//   vagrantClientMock := NewMockVagrantClientStatus(mockStatus)
//
//   status, state, err := CheckForExistingCluster(util.MockAppDir, clusterName, vagrantClientMock)
//
//   assert.True(t, status)
//   assert.Equal(t, state, "directory")
//   assert.NoError(t, err)
// }

// func TestCheckForExistingClusterWithVms(t *testing.T) {
//   mockStatus := map[string]string{
//     "vm-1": "running",
//     "vm-2": "running",
//   }
//   clusterName := "test-cluster"
//
//   err := util.MockAppDirSetup()
//   assert.NoError(t, err)
//
//   defer util.MockAppDirCleanup()
//
//   err = os.Mkdir(filepath.Join(util.MockAppDir, clusterName), 0755)
//   assert.NoError(t, err)
//
//   vagrantClientMock := NewMockVagrantClientStatus(mockStatus)
//  
//   status, state, err := CheckForExistingCluster(util.MockAppDir, clusterName, vagrantClientMock)
//
//   assert.True(t, status)
//   assert.Equal(t, state, "created")
//   assert.NoError(t, err)
// }


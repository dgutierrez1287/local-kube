package cluster

import (
	"testing"
  //"os"
  //"path/filepath"

	"github.com/dgutierrez1287/local-kube/util"
	"github.com/stretchr/testify/assert"
  //vagrant "github.com/bmatcuk/go-vagrant"
)

func TestCheckForExistingClusterNoDir(t *testing.T) {
  clusterName := "test-cluster"

  err := util.MockAppDirSetup()
  assert.NoError(t, err)

  defer util.MockAppDirCleanup()

  status, state, err := CheckForExistingCluster(util.MockAppDir, clusterName, false)
  
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
//   mockClient := new(MockVagrantClient)
// 	mockStatus := new(MockStatusCommand)
//
// 	mockStatus.On("Start").Return(nil)
// 	mockStatus.On("Wait").Return(nil)
// 	mockStatus.StatusResponse = vagrant.StatusResponse{
// 		Status: map[string]string{"default": "running"},
// 	}
//
// 	mockClient.On("Status").Return(mockStatus)
//
// 	originalNewClient := NewVagrantClient
// 	defer func() { NewVagrantClient = originalNewClient }()
// 	NewVagrantClient = func(dir string) (VagrantClientInterface, error) {
// 		return mockClient, nil
// 	}
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
//  
//   status, state, err := CheckForExistingCluster(util.MockAppDir, clusterName, vagrantClientMock)
//
//   assert.True(t, status)
//   assert.Equal(t, state, "created")
//   assert.NoError(t, err)
// }


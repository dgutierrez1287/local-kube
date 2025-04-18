package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterFeaturesInitialization(t *testing.T) {
  features := ClusterFeatures {
    CniController: "cilium",
    StorageController: "longhorn",
    ManagedStorageController: false,
  }

  assert.False(t, features.ManagedStorageController)
  assert.Equal(t, "cilium", features.CniController)
  assert.Equal(t, "longhorn", features.StorageController)
}


/*
      Tests for SetDefaults
*/
func TestClusterFeaturesDefaultValues(t *testing.T) {
  features := ClusterFeatures{}

  err := features.SetDefaults("single", "")
  assert.NoError(t, err)

  // Verification
  assert.Equal(t, features.KubeVersion, "1.31.4")
  assert.Equal(t, features.CniController, "flannel")
  assert.Equal(t, features.IngressController, "native-traefik")
  assert.Equal(t, features.StorageController, "local-storage")
}

func TestClusterFeaturesDefaultVersions(t *testing.T) {
  features := ClusterFeatures{
    CniController: "cilium",
    StorageController: "longhorn",
    KubeVipEnable: true,
  }

  err := features.SetDefaults("single", "2.2.2.2")
  assert.NoError(t, err)

  // Verification
  assert.Equal(t, features.KubeVersion, "1.31.4")
  assert.Equal(t, features.CniControllerVersion, "1.16.4")
  assert.Equal(t, features.CiliumCliVersion, "0.16.22")
  assert.Equal(t, features.StorageControllerVersion, "1.8.0")
  assert.Equal(t, features.KubeVipVersion, "0.5.0")
}

func TestClusterFeaturesDefaultsKubeVipError(t *testing.T) {
  features := ClusterFeatures{}

  err := features.SetDefaults("ha", "2.2.2.2")
  assert.Error(t, err)
}

func TestClusterFeaturesDefaultsMissingVip(t *testing.T) {
  features := ClusterFeatures{
    KubeVipEnable: true,
  }

  err := features.SetDefaults("single", "")
  assert.Error(t, err)
}

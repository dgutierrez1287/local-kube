package ansible

import (
	"testing"

	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/stretchr/testify/assert"
)

var leadNodes []settings.Machine = []settings.Machine{
  {Name: "lead01", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "200GB"},
  {Name: "lead02", IpAddress: "192.168.1.2", Memory: 4, Cpu: 2, DiskSize: "200GB"},
  {Name: "lead03", IpAddress: "192.168.1.3", Memory: 4, Cpu: 2, DiskSize: "200GB"},
}

/*
      Tests for GenerateVarsFile
*/

func TestGenerateVarsFileHa(t *testing.T) {

}

func TestGenerateVarsFileSingle(t *testing.T) {

}

func TestGenerateVarsFileUnmarshalError(t *testing.T) {

}

func TestGenerateVarsFileError(t *testing.T) {

}


/*
      Tests for getGeneralVars
*/
func TestGetGeneralVars(t *testing.T) {

}

/*
      Tests for getKubeVipVars
*/
func TestGetKubeVipVarsEnabled(t *testing.T) {
  vip := "192.168.40.40"

  features := settings.ClusterFeatures{
    KubeVipEnable: true,
    KubeVipVersion: "0.4.0",
  }

  expected := KubeVipVars{
    Version: "0.4.0",
    Vip: vip,
    Enable: true,
  }

  actual := getKubeVipVars(features, vip)

  assert.Equal(t, actual, expected)
}

func TestGetKubeVipVarsNotEnabled(t *testing.T) {
  features := settings.ClusterFeatures{
    KubeVipEnable: false,
  }

  expected := KubeVipVars{}

  actual := getKubeVipVars(features, "")

  assert.Equal(t, actual, expected)
}


/*
      Tests for getCiliumVars
*/
func TestGetCiliumVarsEnabled(t *testing.T) {
  features := settings.ClusterFeatures{
    CniController: "cilium",
    CniControllerVersion: "1.16.3",
    CiliumCliVersion: "0.16.0",
    ManagedCniController: true,
  }

  expected := CiliumVars{
    Version: "1.16.3",
    CliVersion: "0.16.0",
    Install: true,
    InstallHubble: true,
  }

  actual := getCiliumVars(features)

  assert.Equal(t, actual, expected)
}

func TestGetCiliumVarsNotEnabled(t *testing.T) {
  features := settings.ClusterFeatures{
    CniController: "flannel",
  }

  expected := CiliumVars{}

  actual := getCiliumVars(features)

  assert.Equal(t, actual, expected)
}


/*
      Tests for getCalicoVars
*/
func TestGetCalicoVarsEnabled(t *testing.T) {
  features := settings.ClusterFeatures{
    CniController: "calico",
    CniControllerVersion: "3.0.0",
    ManagedCniController: true,
  }

  expected := CalicoVars{
    Version: "3.0.0",
    Install: true,
  }

  actual := getCalicoVars(features)

  assert.Equal(t, actual, expected)
}

func TestGetCalicoVarsNotEnabled(t *testing.T) {
  features := settings.ClusterFeatures{
    CniController: "flannel",
  }

  expected := CalicoVars{}

  actual := getCalicoVars(features)

  assert.Equal(t, actual, expected)
}


/*
      Tests for getLonghornVars
*/
func TestGetLonghornVarsEnabled(t *testing.T) {
  features := settings.ClusterFeatures{
    StorageController: "longhorn",
    StorageControllerVersion: "0.5.0",
    ManagedStorageController: true,
  }

  expected := LonghornVars{
    Version: "0.5.0",
    Install: true,
  }

  actual := getLonghornVars(features)

  assert.Equal(t, actual, expected)
}

func TestGetLonghornVarsNotEnabled(t *testing.T) {
  features := settings.ClusterFeatures{
    StorageController: "local-storage",
  }

  expected := LonghornVars{}

  actual := getLonghornVars(features)

  assert.Equal(t, actual, expected)
}


/*
      Tests for GetTlsSanList
*/
func TestGetTlsSanListVipEnabled(t *testing.T) {
  kubeVipIp := "192.168.40.40"

  actual := GetTlsSanList(leadNodes, true, kubeVipIp)

  expected := []string{
    "192.168.1.1",
    "192.168.1.2",
    "192.168.1.3",
    "192.168.40.40",
  }

  assert.Equal(t, actual, expected)
}

func TestGetTlsSanListVipNotEnabled(t *testing.T) {

  actual := GetTlsSanList(leadNodes, false, "")

  expected := []string{
    "192.168.1.1",
    "192.168.1.2",
    "192.168.1.3",
  }

  assert.Equal(t, actual, expected)
}

/*
      Tests for GetControlPlaneList 
*/
func TestGetControlPlaneListHa(t *testing.T) {
  controlNodes := GetControlPlaneList(leadNodes, "ha")

  primaryNode := controlNodes["lead01"]
  secondControl := controlNodes["lead02"]
  thirdControl := controlNodes["lead03"]

  assert.True(t, primaryNode.primary)
  assert.False(t, secondControl.primary)
  assert.False(t, thirdControl.primary)

  assert.Equal(t, primaryNode.ip, "192.168.1.1")
  assert.Equal(t, secondControl.ip, "192.168.1.2")
  assert.Equal(t, thirdControl.ip, "192.168.1.3")
}

func TestGetControlPlaneListSingle(t *testing.T) {
  controlNodes := GetControlPlaneList(leadNodes, "single")

  assert.Empty(t, controlNodes)
}

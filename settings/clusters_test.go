package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
   Tests for GetSecondaryControlNodeNames
*/
func TestGetSecondaryControlNodeNames(t *testing.T) {
	// Test case where there are more than one leader
	cluster := Cluster{
		Leaders: []Machine{
			{Name: "leader1", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "100GB"},
			{Name: "leader2", IpAddress: "192.168.1.2", Memory: 4, Cpu: 2, DiskSize: "100GB"},
			{Name: "leader3", IpAddress: "192.168.1.3", Memory: 4, Cpu: 2, DiskSize: "100GB"},
		},
	}

	// We expect the function to return names of the secondary leaders
	expected := []string{"leader2", "leader3"}
	actual := cluster.GetSecondaryControlNodeNames()

	assert.Equal(t, expected, actual)
}

func TestGetSecondaryControlNodeNames_NoSecondary(t *testing.T) {
	// Test case where there's only one leader (no secondary nodes)
	cluster := Cluster{
		Leaders: []Machine{
			{Name: "leader1", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "100GB"},
		},
	}

	// We expect an empty list since there are no secondary nodes
	expected := []string{}
	actual := cluster.GetSecondaryControlNodeNames()

	assert.Equal(t, expected, actual)
}

/*
        Tests for GetMachineNameList
*/
func TestGetMachineNameList(t *testing.T) {
  cluster := Cluster{
    Leaders: []Machine{
      {Name: "leader1", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "100GB"},
    },
    Workers: []Machine{
      {Name: "worker1", IpAddress: "192.168.1.2", Memory: 4, Cpu: 2, DiskSize: "50GB"},
    },
  }

  expected := []string{"leader1", "worker1"}
  actual := cluster.GetMachineNameList()

  assert.Equal(t, expected, actual)
}


/*
        Tests for GetServerUrl
*/
func TestGetServerUrlVipEnabled(t *testing.T) {
  cluster := Cluster{
    Vip: "192.168.1.40",
    Leaders: []Machine{
      {Name: "leader1", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "100GB"},
    },
    ClusterFeatures: &ClusterFeatures{
      KubeVipEnable: true,
    },
  }

  expected := "https://192.168.1.40:6443"
  actual := cluster.GetServerUrl()

  assert.Equal(t, expected, actual)
} 

func TestGetServerUrlNoVip(t *testing.T) {
  cluster := Cluster{
    Leaders: []Machine{
      {Name: "leader1", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "100GB"},
    },
    ClusterFeatures: &ClusterFeatures{
      KubeVipEnable: false,
    },
  }

  expected := "https://192.168.1.1:6443"
  actual := cluster.GetServerUrl()

  assert.Equal(t, expected, actual)
}


/*
        Tests for GetAnsibleNodeVagrantName
*/
func TestGetAnsibleNodeVagrantNameSingleNode(t *testing.T) {
  cluster := Cluster{
    ClusterType: "single",
  }

  expected := "default"
  actual := cluster.GetAnsibleNodeVagrantName()

  assert.Equal(t, expected, actual)
}

func TestGetAnsibleNodeVagrantNameHaCluster(t *testing.T) {
  cluster := Cluster{
    ClusterType: "ha",
    Leaders: []Machine{
      {Name: "leader1", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "100GB"},
      {Name: "leader2", IpAddress: "192.168.1.2", Memory: 4, Cpu: 2, DiskSize: "50GB"},
    },
  }

  expected := "leader1"
  actual := cluster.GetAnsibleNodeVagrantName()

  assert.Equal(t, expected, actual)
}


/*
        Tests for GetWorkerNodeNames
*/
func TestGetWorkerNodeNames(t *testing.T) {
	// Test case where there are worker nodes
	cluster := Cluster{
		Workers: []Machine{
			{Name: "worker1", IpAddress: "192.168.2.1", Memory: 8, Cpu: 4, DiskSize: "200GB"},
			{Name: "worker2", IpAddress: "192.168.2.2", Memory: 8, Cpu: 4, DiskSize: "200GB"},
		},
	}

	// We expect the function to return names of the worker nodes
	expected := []string{"worker1", "worker2"}
	actual := cluster.GetWorkerNodeNames()

	assert.Equal(t, expected, actual)
}

func TestGetWorkerNodeNames_NoWorkers(t *testing.T) {
	// Test case where there are no worker nodes
	cluster := Cluster{
		Workers: []Machine{},
	}

	// We expect an empty list since there are no worker nodes
	expected := []string{}
	actual := cluster.GetWorkerNodeNames()

	assert.Equal(t, expected, actual)
}

/*
        Tests for GetControlNodeIps
*/
func TestGetControlNodeIps(t *testing.T) {
  cluster := Cluster {
    Leaders: []Machine{
			{Name: "leader1", IpAddress: "192.168.1.1", Memory: 4, Cpu: 2, DiskSize: "100GB"},
			{Name: "leader2", IpAddress: "192.168.1.2", Memory: 4, Cpu: 2, DiskSize: "100GB"},
			{Name: "leader3", IpAddress: "192.168.1.3", Memory: 4, Cpu: 2, DiskSize: "100GB"},
		},
  }

  ips := cluster.GetControlNodeIps()

  verificationArray := []string{
    "192.168.1.1",
    "192.168.1.2",
    "192.168.1.3",
  }

  assert.Equal(t, ips, verificationArray)
}

/*
        Tests for isHA
*/
func TestClusterHa(t *testing.T) {
  cluster := Cluster {
    ClusterType: "ha",
  }

  isHa := cluster.IsHA()
  assert.True(t, isHa)
}

func TestClusterNotHa(t *testing.T) {
  cluster := Cluster {
    ClusterType: "single",
  }

  isHa := cluster.IsHA()
  assert.False(t, isHa)
}

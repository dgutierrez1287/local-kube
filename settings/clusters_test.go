package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

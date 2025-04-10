package cluster

import (
	"context"

	vagrant "github.com/bmatcuk/go-vagrant"
)

type MockVagrantClient struct {
  mockStatus map[string]string 
}

var _ VagrantClientInterface = (*MockVagrantClient)(nil)

func NewMockVagrantClientStatus(status map[string]string) *MockVagrantClient {
  return &MockVagrantClient{mockStatus: status}
}

func (m *MockVagrantClient) Status() *vagrant.StatusCommand {
  return &vagrant.StatusCommand{
    BaseCommand: vagrant.BaseCommand{
      Context: context.Background(),
    },
    StatusResponse: vagrant.StatusResponse{
      Status: m.mockStatus,
    },
    
  }
}

func (m *MockVagrantClient) Up() *vagrant.UpCommand {
  return &vagrant.UpCommand{}
}

func (m *MockVagrantClient) Destroy() *vagrant.DestroyCommand {
  return &vagrant.DestroyCommand{}
}

func (m *MockVagrantClient) SshConfig() *vagrant.SSHConfigCommand {
  return &vagrant.SSHConfigCommand{}
}

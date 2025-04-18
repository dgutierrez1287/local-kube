package cluster

import (

  "github.com/stretchr/testify/mock"
	vagrant "github.com/bmatcuk/go-vagrant"
)

type MockVagrantClient struct {
  mock.Mock
}

func (m *MockVagrantClient) Status() *vagrant.StatusCommand {
	args := m.Called()
	return args.Get(0).(*vagrant.StatusCommand)
}

func (m *MockVagrantClient) Up() *vagrant.UpCommand {
	args := m.Called()
	return args.Get(0).(*vagrant.UpCommand)
}

func (m *MockVagrantClient) Destroy() *vagrant.DestroyCommand {
	args := m.Called()
	return args.Get(0).(*vagrant.DestroyCommand)
}

func (m *MockVagrantClient) SshConfig() *vagrant.SSHConfigCommand {
	args := m.Called()
	return args.Get(0).(*vagrant.SSHConfigCommand)
}


type MockStatusCommand struct {
	mock.Mock
	StatusResponse vagrant.StatusResponse
}

func (m *MockStatusCommand) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStatusCommand) Wait() error {
	args := m.Called()
	return args.Error(0)
}

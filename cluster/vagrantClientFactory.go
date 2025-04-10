package cluster

import (
  vagrant "github.com/bmatcuk/go-vagrant"
)

type VagrantClientInterface interface {
  Status() *vagrant.StatusCommand
  Up() *vagrant.UpCommand
  Destroy() *vagrant.DestroyCommand
  SshConfig() *vagrant.SSHConfigCommand
}

//type VagrantClientFactory func(vagrantDirPath string) (VagrantClientInterface, error)

type DefaultVagrantClient struct {
  client *vagrant.VagrantClient
}

func (v *DefaultVagrantClient) Status() *vagrant.StatusCommand {
  return v.client.Status()
}

func (v *DefaultVagrantClient) Up() *vagrant.UpCommand {
  return v.client.Up()
}

func (v *DefaultVagrantClient) SshConfig() *vagrant.SSHConfigCommand {
  return v.client.SSHConfig()
}

func (v *DefaultVagrantClient) Destroy() *vagrant.DestroyCommand {
  return v.client.Destroy()
}

func NewVagrantClient(vagrantDirPath string) (VagrantClientInterface, error) {
  client, err := vagrant.NewVagrantClient(vagrantDirPath)
  if err != nil {
    return nil, err
  }
  return &DefaultVagrantClient{client: client}, nil
}

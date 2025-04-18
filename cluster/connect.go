package cluster

import (
	"errors"
	"os"
  "fmt"

	vagrant "github.com/bmatcuk/go-vagrant"
	"github.com/dgutierrez1287/local-kube/logger"
	"golang.org/x/crypto/ssh"
)

/*
This will return the ssh command to connect to a given machine
in a cluster
*/
func GetSshConfigs(clusterDir string, nodeName string) (vagrant.SSHConfig, error) {
  var sshConfig vagrant.SSHConfig

  logger.LogDebug("Getting vagrant client")
  client, err := NewVagrantClient(clusterDir)
  if err != nil {
    logger.LogError("Error getting vagrant client")
    return sshConfig, err
  }

  sshCmd := client.SshConfig()
  sshCmd.Host = nodeName

  if sshCmd == nil {
    logger.LogError("Ssh config command is nil")
    return sshConfig, errors.New("ssh config command is nil")
  }

  err = sshCmd.Run()
  if err != nil {
    logger.LogError("Error running the ssh config command")
    return sshConfig, err
  }

  configs := sshCmd.Configs
  if len(configs) == 0 {
    logger.LogError("Error ssh configs are empty")
    return sshConfig, errors.New("ssh configs are empty")
  }

  sshConfig = configs[nodeName]

  return sshConfig, nil

}

/*
opens an an ssh session
*/
func OpenSshSession(configs vagrant.SSHConfig, signer ssh.Signer) error {
  
  sshConfig := &ssh.ClientConfig{
		User: configs.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

  addr := fmt.Sprintf("%s:%d", configs.HostName, configs.Port)
	clientConn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
    logger.LogError("Failed to dial SSH")
    return err
	}
	defer clientConn.Close()

	session, err := clientConn.NewSession()
	if err != nil {
    logger.LogError("failed to create SSH session")
    return err
	}
	defer session.Close()

	// Set up terminal
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
    logger.LogError("Request for pseudo terminal failed")
    return err
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
    logger.LogError("failed to start shell")
    return err
	}
	return session.Wait()
}

/*
Loads a private key to connect to a vm
*/
func LoadPrivateKey(path string) (ssh.Signer, error) {
  key, err := os.ReadFile(path)
	if err != nil {
    logger.LogError("Error loading private key", "path", path)
    return nil, err
	}
	return ssh.ParsePrivateKey(key)
}

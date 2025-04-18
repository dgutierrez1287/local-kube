package kubeconfig

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v3"
)

type KubeConfig struct {
	APIVersion     string           `yaml:"apiVersion"`
	Kind           string           `yaml:"kind"`
	Clusters       []NamedCluster   `yaml:"clusters"`
	Contexts       []NamedContext   `yaml:"contexts"`
	CurrentContext string           `yaml:"current-context"`
	Users          []NamedUser      `yaml:"users"`
	Preferences    Preferences      `yaml:"preferences"`
}

type NamedCluster struct {
	Name    string  `yaml:"name"`
	Cluster Cluster `yaml:"cluster"`
}

type Cluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthority     string `yaml:"certificate-authority,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

type NamedContext struct {
	Name    string  `yaml:"name"`
	Context Context `yaml:"context"`
}

type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

type NamedUser struct {
	Name string `yaml:"name"`
	User User   `yaml:"user"`
}

type User struct {
	ClientCertificate     string `yaml:"client-certificate,omitempty"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKey             string `yaml:"client-key,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
	Token                 string `yaml:"token,omitempty"`
	Username              string `yaml:"username,omitempty"`
	Password              string `yaml:"password,omitempty"`
	AuthProvider          map[string]interface{} `yaml:"auth-provider,omitempty"`
	Exec                  map[string]interface{} `yaml:"exec,omitempty"`
}

type Preferences struct {
	Colors bool `yaml:"colors,omitempty"`
}

/*
This will update the server url for a given cluster in a kubeconfig. This
update will be done in place 
*/
func (kubeConfig *KubeConfig) UpdateServerUrl(newServerUrl string, clusterName string) error {
  clusterFound := false
  logger.LogDebug("getting the current server url for cluster", "name", clusterName)

  for i := range kubeConfig.Clusters {
    if kubeConfig.Clusters[i].Name == clusterName {
      clusterFound = true
      logger.LogDebug("updating server url", "current", kubeConfig.Clusters[i].Cluster.Server, "new", newServerUrl)
      kubeConfig.Clusters[i].Cluster.Server = newServerUrl
    }
  }

  if !clusterFound {
    logger.LogError("Error cluster was not found in kubeconfig")
    return errors.New("Cluster doesn't exist")
  }
  return nil
}

/*
This will add a cluster to a kubeconfig, this assumes that the cluster name,
context name, and user name all match the sourceClusterName and they will
match the destClusterName in the resulting file
*/
func (kubeConfig *KubeConfig) AddCluster(sourceKubeConfig KubeConfig, sourceClusterName string, destClusterName string) error {
  sourceClusterFound := false
  sourceClusterData := Cluster{}
  sourceUserFound := false
  sourceUserData := User{}
  
  // get the cluster information
  logger.LogDebug("Looking for cluster in kube config", "name", sourceClusterName)
  for _, cluster := range(sourceKubeConfig.Clusters) {
    if cluster.Name == sourceClusterName {
      logger.LogDebug("Cluster found")
      sourceClusterFound = true
      sourceClusterData = cluster.Cluster
    }
  }

  // get the user information
  logger.LogDebug("Looking for user in kube config", "name", sourceClusterName)
  for _, user := range(sourceKubeConfig.Users) {
    if user.Name == sourceClusterName {
      logger.LogDebug("User found")
      sourceUserFound = true
      sourceUserData = user.User
    }
  }

  if !sourceClusterFound {
    logger.LogError("Error source cluster not found in source kube config")
    return errors.New("source cluster missing")
  }

  if !sourceUserFound {
    logger.LogError("Error source user not found in source kube config")
    return errors.New("source user missing")
  }

  logger.LogDebug("Creating new context")
  newContext := Context{
    Cluster: destClusterName,
    User: destClusterName,
  }
  newNamedContext := NamedContext{
    Name: destClusterName,
    Context: newContext,
  }

  logger.LogDebug("Creating new cluster")
  newNamedCluster := NamedCluster{
    Name: destClusterName,
    Cluster: sourceClusterData,
  }

  logger.LogDebug("Creating new User")
  newNamedUser := NamedUser{
    Name: destClusterName,
    User: sourceUserData,
  }

  logger.LogDebug("Adding new cluster to kubeconfig")
  kubeConfig.Clusters = append(kubeConfig.Clusters, newNamedCluster)
  kubeConfig.Contexts = append(kubeConfig.Contexts, newNamedContext)
  kubeConfig.Users = append(kubeConfig.Users, newNamedUser)

  return nil
}

/*
This will remove a cluster from a kubeConfig
*/
func (kubeConfig *KubeConfig) RemoveCluster(clusterName string) {
  clusterFound := false
  var clusterIndex int

  contextFound := false
  var contextIndex int

  userFound := false
  var userIndex int
  
  // find cluster index
  logger.LogDebug("Getting index for cluster", "name", clusterName)
  for index, cluster := range(kubeConfig.Clusters) {
    if cluster.Name == clusterName {
      logger.LogDebug("Cluster found")
      clusterFound = true
      clusterIndex = index
    }
  }
  
  // remove cluster if found
  if clusterFound {
    logger.LogDebug("Cluster found, Removing...", "name", clusterName)
    clusterResult := append(kubeConfig.Clusters[:clusterIndex], kubeConfig.Clusters[clusterIndex+1:]...)
    kubeConfig.Clusters = clusterResult
    logger.LogDebug("Cluster removed", "name", clusterName)
  } else {
    logger.LogDebug("No cluster with the name found", "name", clusterName)
  }

  // find context index
  logger.LogDebug("Getting index for context", "name", clusterName)
  for index, context := range(kubeConfig.Contexts) {
    if context.Name == clusterName {
      logger.LogDebug("Context found")
      contextFound = true
      contextIndex = index
    }
  }

  // remove context if found
  if contextFound {
    logger.LogDebug("Context found, Removing...", "name", clusterName)
    contextResult := append(kubeConfig.Contexts[:contextIndex], kubeConfig.Contexts[contextIndex+1:]...)
    kubeConfig.Contexts = contextResult
    logger.LogDebug("Context removed", "name", clusterName)
  } else {
    logger.LogDebug("No context with the name found", "name", clusterName)
  }

  // find user index
  logger.LogDebug("Getting index for the user", "name", clusterName)
  for index, user := range(kubeConfig.Users) {
    if user.Name == clusterName {
      logger.LogDebug("User found")
      userFound = true
      userIndex = index
    }
  }

  // remove user if found
  if userFound {
    logger.LogDebug("User found, Removing...", "name", clusterName)
    userResult := append(kubeConfig.Users[:userIndex], kubeConfig.Users[userIndex+1:]...)
    kubeConfig.Users = userResult
    logger.LogDebug("User removed", "name", clusterName)
  } else {
    logger.LogDebug("No user with the name found", "name", clusterName)
  }
}

/*
Reads a kubeconfig file 
*/
func ReadKubeConfig(filePath string) (KubeConfig, error) {
  
  logger.LogDebug("reading kubeconfig", "path", filePath)
  file, err := os.Open(filePath)
  if err != nil {
    logger.LogError("Error opening kubeconfig file")
    return KubeConfig{}, err
  }
  defer file.Close()

  bytes, err := io.ReadAll(file)
  if err != nil {
    logger.LogError("Error reading kubeconfig file")
    return KubeConfig{}, err
  }

  var kubeConfig KubeConfig
  err = yaml.Unmarshal(bytes, &kubeConfig)
  if err != nil {
    logger.LogError("Error unmarshaling kubeconfig yaml to struct")
    return KubeConfig{}, err
  }
  
  logger.LogDebug("Kubeconfig file is read successfully", "file", filePath)
  return kubeConfig, nil
}

/*
Writes a kubeconfig file
*/
func WriteKubeConfig(filePath string, kubeConfig KubeConfig) error {

  logger.LogDebug("Masharling kubeconfig to yaml")
  yamlData, err := yaml.Marshal(&kubeConfig)
  if err != nil {
    logger.LogError("Error marshaling kube config to yaml")
    return err
  }
  
  logger.LogDebug("Creating or truncating current kubeconfig for writing")
  file, err := os.Create(filePath)
  if err != nil {
    logger.LogError("Error creating or truncating kubeconfig")
    return err
  }
  defer file.Close()

  _, err = file.Write(yamlData)
  if err != nil {
    logger.LogError("Error writing kubeconfig to file")
  }
  return nil
}

/*
This will back up the current kubeconfig before changes are made
*/
func BackupKubeConfig(appDir string, filePath string) error {
  backupFilePath := filepath.Join(appDir, ".kubeconfig.bak")

  err := copy.Copy(filePath, backupFilePath)
  if err != nil {
    logger.LogError("Error backing up the current kubeconfig file")
    return err
  }
  return nil
}

/*
This will restore the backuo kube config and remove the current one 
incase anything goes wrong
*/
func RestoreKubeConfigBackup(appDir string, kubeConfigPath string) error {
  kubeConfigBackupPath := filepath.Join(appDir, ".kubeconfig.bak")

  logger.LogDebug("Removing the current kubeconfig")

  err := os.Remove(kubeConfigPath)
  if err != nil && !os.IsNotExist(err) {
    logger.LogError("Error removing the current kubeconfig")
    return err
  }

  logger.LogInfo("Copying kubeconfig backup to kubeconfig location")

  err = copy.Copy(kubeConfigBackupPath, kubeConfigPath)
  if err != nil {
    logger.LogError("Error copying the backup to the kubeconfig location")
    return err
  }
  return nil 
}

/*
This will clean up any backup of the kubeconfig once the 
action has been successful
*/
func CleanKubeConfigBackup(appDir string) error {
  backupFilePath := filepath.Join(appDir, ".kubeconfig.bak") 

  err := os.Remove(backupFilePath)
  if err != nil {
    logger.LogError("Error cleaning up the kubeconfig backup file")
    return err
  }
  return nil
}

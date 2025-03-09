package settings

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
)

/*
  Settings - Local kube settings
*/
type Settings struct {
  ProvisionSettings ProvisionSettings   `json:"provision"`  // Provision settings
  Providers map[string]Provider         `json:"providers"`  // Providers
  Clusters map[string]Cluster           `json:"clusters"`   // Clusters
}

//validate settings 
// TODO: right now this will only validate the cluster 
// name is in the settings more validation can be done later
func (settings *Settings) SettingsValid(clusterName string) bool {
  //validate cluster name is in settings
  if _, exists := settings.Clusters[clusterName]; exists {
    logger.Logger.Debug("cluster exists in settings", "cluster", clusterName)
    return true
  } else {
    logger.Logger.Error("cluster is not present in settings", "cluster", clusterName)
    return false
  }
}

// Check if setting file exist
func SettingsFileExists() bool {
  appDir := GetAppDirPath()
  settingsFilePath := filepath.Join(appDir, "settings.json")
  var settingsExists bool

  if _, err := os.Stat(settingsFilePath); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.Logger.Debug("Settings file does not exist")
      settingsExists = false
    } else {
      logger.Logger.Error("Error checking if settings file exists")
      os.Exit(120)
    }
  } else {
    logger.Logger.Debug("Settings file exists")
    settingsExists = true
  }

  return settingsExists
}

// create a default settings file (mostly blank) with 
// a couple known defaults filled in
func CreateDefaultSettingsFile() {
  appDir := GetAppDirPath()
  settingsFile := filepath.Join(appDir, "settings.json")

  defaultAnsibleRoles := make(map[string]AnsibleRole)
  defaultAnsibleRoles["kube"] = AnsibleRole{
    LocationType: "git",
    Location: "github.com/dgutierrez1287/ansible-role-kube",
    GitBranch: "master",
  }

  emptySettings := Settings{
    ProvisionSettings: ProvisionSettings{
      AnsibleVersion: "2.17.6",
      AnsibleRoles: defaultAnsibleRoles,
    },
    Providers: make(map[string]Provider),
    Clusters: make(map[string]Cluster),
  }
  logger.Logger.Debug("creating an empty settings with defauts", "settings", emptySettings)

  logger.Logger.Debug("Creating new settings file")
  file, err := os.Create(settingsFile)
  if err != nil {
    logger.Logger.Error("Error creating blank settings file", "error", err)
    os.Exit(120)
  }
  defer file.Close()

  encoder := json.NewEncoder(file)
  encoder.SetIndent("", " ")

  logger.Logger.Debug("Writing default settings to file")
  if err := encoder.Encode(emptySettings); err != nil {
    logger.Logger.Error("Error writing settings defaults to file", "error", err)
    os.Exit(120)
  }
}

// Read the settings file and return a settiings object
func ReadSettingsFile() (Settings, error){
  appDir := GetAppDirPath()
  settingsFile := filepath.Join(appDir, "settings.json")

  file, err := os.Open(settingsFile) 
  if err != nil {
    logger.Logger.Error("Error opening settings file", "error", err)
    return Settings{}, errors.New("error opening settings")
  }
  defer file.Close()

  bytes, err := io.ReadAll(file)
  if err != nil {
    logger.Logger.Error("Error reading settings file", "error", err)
    return Settings{}, errors.New("error Reading Settings")
  }

  var settings Settings
  err = json.Unmarshal(bytes, &settings)
  if err != nil {
    logger.Logger.Error("Error unmarshaling json to struct", "error", err)
    return Settings{}, errors.New("error unmarshaling settings")
  }

  return settings, nil
}



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

/* 
validate settings
TODO: right now this will only validate the cluster 
 name is in the settings more validation can be done later
*/
func (settings *Settings) SettingsValid(clusterName string) bool {
  //validate cluster name is in settings
  if _, exists := settings.Clusters[clusterName]; exists {
    logger.LogDebug("cluster exists in settings", "cluster", clusterName)
    return true
  } else {
    logger.LogError("cluster is not present in settings", "cluster", clusterName)
    return false
  }
}

/*
SettingsFileExists()
Checks if the settings file exists or not
*/
func SettingsFileExists(appDir string) (bool, error) {
  settingsFilePath := filepath.Join(appDir, "settings.json")

  _, err := os.Stat(settingsFilePath)

  if errors.Is(err, os.ErrNotExist) {
    logger.LogDebug("Settings file does not exist")
    return false, nil
  } else if err != nil {
    logger.LogError("Error checking if settings file exists")
    return false, err
  }
  logger.LogDebug("Settings file exists")
  return true,nil
}

/*
CreateDefaultSettingsFile()
create a default settings file (mostly blank) with 
a couple known defaults filled in
*/
func CreateDefaultSettingsFile(appDir string) error {
  settingsFile := filepath.Join(appDir, "settings.json")

  defaultAnsibleRoles := make(map[string]AnsibleRole)
  defaultAnsibleRoles["kube"] = AnsibleRole{
    LocationType: "git",
    Location: "https://github.com/dgutierrez1287/ansible-role-kube",
    RefType: "branch",
    GitRef: "master",
  }

  emptySettings := Settings{
    ProvisionSettings: ProvisionSettings{
      AnsibleVersion: "2.17.6",
      AnsibleRoles: defaultAnsibleRoles,
    },
    Providers: make(map[string]Provider),
    Clusters: make(map[string]Cluster),
  }
  logger.LogDebug("creating an empty settings with defauts", "settings", emptySettings)

  logger.LogDebug("Creating new settings file")
  file, err := os.Create(settingsFile)
  if err != nil {
    logger.LogError("Error creating blank settings file")
    return err
  }
  defer file.Close()

  encoder := json.NewEncoder(file)
  encoder.SetIndent("", " ")

  logger.LogDebug("Writing default settings to file")
  if err := encoder.Encode(emptySettings); err != nil {
    logger.LogError("Error writing settings defaults to file")
    return err
  }
  logger.LogDebug("Default settings file created successfully")
  return nil 
}

/*
ReadSettingsFile()
Reads the settings file and returns a settings object
this will contain all the settings in the settings file 
*/
func ReadSettingsFile(appDir string) (Settings, error){
  settingsFile := filepath.Join(appDir, "settings.json")

  file, err := os.Open(settingsFile) 
  if err != nil {
    logger.LogError("Error opening settings file")
    return Settings{}, err
  }
  defer file.Close()

  bytes, err := io.ReadAll(file)
  if err != nil {
    logger.LogError("Error reading settings file")
    return Settings{}, err
  }

  var settings Settings
  err = json.Unmarshal(bytes, &settings)
  if err != nil {
    logger.LogError("Error unmarshaling json to struct")
    return Settings{}, err
  }
  logger.LogDebug("Settings file read successfully")
  return settings, nil
}



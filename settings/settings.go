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
    logger.Logger.Debug("cluster exists in settings", "cluster", clusterName)
    return true
  } else {
    logger.Logger.Error("cluster is not present in settings", "cluster", clusterName)
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
    logger.Logger.Debug("Settings file does not exist")
    return false, nil
  } else if err != nil {
    logger.Logger.Error("Error checking if settings file exists")
    return false, err
  }
  logger.Logger.Debug("Settings file exists")
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
  logger.Logger.Debug("creating an empty settings with defauts", "settings", emptySettings)

  logger.Logger.Debug("Creating new settings file")
  file, err := os.Create(settingsFile)
  if err != nil {
    logger.Logger.Error("Error creating blank settings file")
    return err
  }
  defer file.Close()

  encoder := json.NewEncoder(file)
  encoder.SetIndent("", " ")

  logger.Logger.Debug("Writing default settings to file")
  if err := encoder.Encode(emptySettings); err != nil {
    logger.Logger.Error("Error writing settings defaults to file")
    return err
  }
  logger.Logger.Debug("Default settings file created successfully")
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
    logger.Logger.Error("Error opening settings file")
    return Settings{}, err
  }
  defer file.Close()

  bytes, err := io.ReadAll(file)
  if err != nil {
    logger.Logger.Error("Error reading settings file")
    return Settings{}, err
  }

  var settings Settings
  err = json.Unmarshal(bytes, &settings)
  if err != nil {
    logger.Logger.Error("Error unmarshaling json to struct")
    return Settings{}, err
  }
  logger.Logger.Debug("Settings file read successfully")
  return settings, nil
}



package ansible

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
)

var ansibleRoleDir = "ansible-roles"
var roleCacheFileName = ".role-cache.json"

/*
  RoleCache - holds a map of AnsibleRole type to maintain a
  cache of the current state of ansible roles
*/
type RoleCache struct {
  Roles map[string]settings.AnsibleRole   `json:"roles"`    // the current state of the ansible roles
}

/*
RoleCacheFileExists() 
Checks if the role cache file exists, if the file isn't there 
it will clear out all the role directories and start fresh
*/
func RoleCacheFileExists(appDir string) (bool, error){
  roleCacheFile := filepath.Join(appDir, ansibleRoleDir, roleCacheFileName)

  if _, err := os.Stat(roleCacheFile); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      logger.Logger.Debug("Role cache file does not exist")
      return false, nil
    } else {
      logger.Logger.Error("Error checking for the role cache file")
      return false, err
    }
  } else {
    logger.Logger.Debug("Role cache file exists")
  }
  return true, nil
}

/*
RoleCacheFileDelete()
This will delete the role cache file, this 
is to be used when all roles are cleaned and 
as a reset
*/
func RoleCacheFileDelete(appDir string) (error) {
  roleCacheFile := filepath.Join(appDir, ansibleRoleDir, roleCacheFileName)

  err := os.Remove(roleCacheFile)
  if err != nil {
    logger.Logger.Error("There was an error removing the role cache file")
    return err
  }

  logger.Logger.Debug("Role cache file has been removed")
  return nil 
}

/*
ReadRoleCache()
Read the current ansible role cache file and returns 
a RoleCache object from it 
*/
func ReadRoleCache(appDir string) (RoleCache, error) {
  roleCacheFile := filepath.Join(appDir, ansibleRoleDir, roleCacheFileName)

  logger.Logger.Debug("Opening the role cache file", "file", roleCacheFile)
  file, err := os.Open(roleCacheFile)
  if err != nil {
    logger.Logger.Error("Error opening ansible role cache file")
    return RoleCache{}, err
  }
  defer file.Close()

  logger.Logger.Debug("Reading role cache file")
  bytes, err := io.ReadAll(file)
  if err != nil {
    logger.Logger.Error("Error reading ansible role cache file")
    return RoleCache{}, err
  }

  var roleCache RoleCache
  logger.Logger.Debug("Unmarshaling role cache")
  err = json.Unmarshal(bytes, &roleCache)
  if err != nil {
    logger.Logger.Error("Error unmarshaling json to struct")
    return RoleCache{}, err
  }

  logger.Logger.Debug("cache", roleCache)
  return roleCache, nil
}

/*
WriteRoleCache()
Writes the new ansible role cache to the file
*/
func WriteRoleCache(appDir string, rolecache RoleCache) error {
  roleCacheFile := filepath.Join(appDir, ansibleRoleDir, roleCacheFileName)

  logger.Logger.Debug("Creating or truncating role cache file", "file", roleCacheFile)
  file, err := os.Create(roleCacheFile)
  if err != nil {
    logger.Logger.Error("Error creating ansible role cache file")
    return err
  }
  defer file.Close()

  encoder := json.NewEncoder(file)
  encoder.SetIndent("", " ")

  logger.Logger.Debug("Writing role cache to file")
  if err := encoder.Encode(rolecache); err != nil {
    logger.Logger.Error("Error writing role cache to file")
    return err
  }
  return nil
}

package cmd

import (
	"fmt"
	"os"

	"github.com/dgutierrez1287/local-kube/ansible"
	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/dgutierrez1287/local-kube/util"
	"github.com/spf13/cobra"
)

var roleSyncCmd = &cobra.Command{
  Use: "roles-sync",
  Short: "Syncs ansible roles",
  Long: "Syncs ansible roles that are used by clusters",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(util.TitleText)

    appDir := settings.GetAppDirPath()

    logger.Logger.Info("Syncing ansible roles")

    // check for settings file
    logger.Logger.Debug("Checking if settings file exists")
    settingsExists, err := settings.SettingsFileExists(appDir)
    if err != nil {
      logger.Logger.Error("Error checking if settings file exists", "error", err)
      os.Exit(120)
    }

    if !settingsExists {
      logger.Logger.Info("Could not file settings file, have you initialized local-kube?")
      os.Exit(120)
    }

    // read settings file
    logger.Logger.Info("Reading settings file")
    appSettings, err := settings.ReadSettingsFile(appDir)
    if err != nil {
      logger.Logger.Error("Error reading settings file", "error", err)
    }

    // check for role cache file
    logger.Logger.Debug("Checking if role cache file exists")
    cacheExists, err := ansible.RoleCacheFileExists(appDir)
    if err != nil {
      logger.Logger.Error("Error checking for role cache file", "error", err)
      os.Exit(120)
    }
    
    if !cacheExists {
      roleCache := ansible.RoleCache{
        Roles: make(map[string]settings.AnsibleRole),
      }

      logger.Logger.Info("No current role cache exists, installing roles ")
      
      // Go through roles in settings file and create all the roles
      for roleName, role := range appSettings.ProvisionSettings.AnsibleRoles {
        logger.Logger.Info("Processing role", "roleName", roleName)
        if role.LocationType == "git"{

          logger.Logger.Info("Installing role from git", "roleName", roleName)
          logger.Logger.Debug("roleName", roleName, "role", role)
          err := ansible.CreateGitRole(appDir, roleName, role)
          if err != nil {
            logger.Logger.Error("Error installing role from git", "error", err)
            
            // write out currently installed roles
            logger.Logger.Info("Writing out all currently installed roles, please check config")
            err = ansible.WriteRoleCache(appDir, roleCache)
            if err != nil {
              logger.Logger.Error("Error writing out role cache, manual intervention needed", "error", err)
              os.Exit(200)
            }
            logger.Logger.Info("Partial role cache written")
            os.Exit(120)
          }

          logger.Logger.Info("Adding role to cache file", "roleName", roleName)
          roleCache.Roles[roleName] = appSettings.ProvisionSettings.AnsibleRoles[roleName]

        } else if role.LocationType == "local" {

          logger.Logger.Info("Installing role from the local filesystem", "roleName", roleName)
          logger.Logger.Debug("roleName", roleName, "role", role)
          err := ansible.CreateLocalRole(appDir, roleName, role)
          if err != nil {
            logger.Logger.Error("Error installing role from filesystem", "error", err)

            //write out currently installed roles
            logger.Logger.Info("Writing out currently installed roles, please check the config")
            err = ansible.WriteRoleCache(appDir, roleCache)
            if err != nil {
              logger.Logger.Error("Error writing out the role cache. manual intervention needed", "error", err)
              os.Exit(200)
            }
            logger.Logger.Info("Partial role cache written")
            os.Exit(120)
          }

          logger.Logger.Info("Adding role to cache file", "roleName", roleName)
          roleCache.Roles[roleName] = appSettings.ProvisionSettings.AnsibleRoles[roleName]

        } else {
          // if the role location type isnt supported error out 
          // and roll back (clear all roles)
          logger.Logger.Error("Unsupported location type", "roleName", roleName, "type", role.LocationType)
          logger.Logger.Info("Please check config for the role and rerun")
          os.Exit(120)
        }

        logger.Logger.Info("Writing cache to file")
        err := ansible.WriteRoleCache(appDir, roleCache)
        if err != nil {
          logger.Logger.Error("Error writing role cache", "error", err)
          os.Exit(120)
        }
      }

    } else {
      logger.Logger.Info("Role cache file exists, updating roles")

      currentRoles, err := ansible.ReadRoleCache(appDir)
      if err != nil {
        logger.Logger.Error("Error reading role cache file", "error", err)
        os.Exit(120)
      }

      logger.Logger.Info("Checking for roles to add, update or remove")
      
      rolesToAdd, 
      rolesToUpdateInPlace, 
      rolesToCleanReAdd, 
      rolesToRemove, 
      err := ansible.RoleReconcileLists(currentRoles.Roles, appSettings.ProvisionSettings.AnsibleRoles)
      
      if err != nil {
        logger.Logger.Error("Error getting lists for role actions", "error", err)
        os.Exit(120)
      }

      logger.Logger.Debug("roles to add", "roles", rolesToAdd)
      logger.Logger.Debug("roles to remove", "roles", rolesToRemove)
      logger.Logger.Debug("roles to clear and then readd", "roles", rolesToCleanReAdd)
      logger.Logger.Debug("roles to update in place", "roles", rolesToUpdateInPlace)

      // add roles
      logger.Logger.Debug("Adding new roles", "count", len(rolesToAdd))
      for _, roleName := range rolesToAdd {
        logger.Logger.Info("Adding role", "roleName", roleName)
        
        roleData := appSettings.ProvisionSettings.AnsibleRoles[roleName]

        if roleData.LocationType == "git" {
          logger.Logger.Debug("Role is git based")
          err := ansible.CreateGitRole(appDir, roleName, roleData)

          if err != nil {
            logger.Logger.Error("Error adding git role", "roleName", roleName, "error", err)
            os.Exit(200)
          }

        } else {
          logger.Logger.Debug("Role is local based")
          err := ansible.CreateLocalRole(appDir, roleName, roleData)

          if err != nil {
            logger.Logger.Error("Error adding local role", "roleName", roleName, "error", err)
            os.Exit(200)
          }
        }
        logger.Logger.Debug("Adding role to cache")
        currentRoles.Roles[roleName] = roleData
        err := ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.Logger.Error("Error updating cache", "error", err)
          os.Exit(200)
        }

        logger.Logger.Info("Successfully added role", "roleName", roleName)
      }

      // remove roles
      logger.Logger.Debug("Removing roles", "count", len(rolesToRemove))
      for _, roleName := range rolesToRemove {
        logger.Logger.Info("Removing role", "roleName", roleName)

        err := ansible.ClearRole(appDir, roleName)
        if err != nil {
          logger.Logger.Error("Error removing role, please check and re-run", "error", err)
          os.Exit(200)
        }

        logger.Logger.Debug("removing role from cache")
        delete(currentRoles.Roles, roleName)
        err = ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.Logger.Error("Error updating cache", "error", err)
          os.Exit(200)
        }
        logger.Logger.Info("Successfully removed role", "roleName", roleName)
      }

      // update in place
      logger.Logger.Debug("Updating roles in place", "count", len(rolesToUpdateInPlace))
      for _, roleName := range rolesToUpdateInPlace {
        logger.Logger.Info("Updating role in place", "roleName", roleName)

        currentRole := currentRoles.Roles[roleName]
        newRole := appSettings.ProvisionSettings.AnsibleRoles[roleName]
        
        logger.Logger.Debug("current role data", "role", currentRole)
        logger.Logger.Debug("new role data", "role", newRole)
        
        err = ansible.UpdateGitRole(appDir, roleName, currentRole, newRole)
        if err != nil {
          logger.Logger.Error("Error updating role", "error", err)
          os.Exit(200)
        }

        logger.Logger.Debug("Updating role cache")
        currentRoles.Roles[roleName] = newRole
        err = ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.Logger.Error("Error updating cache", "error", err)
          os.Exit(200)
        }
        logger.Logger.Info("Successfully updated role", "roleName", roleName)
      }

      // clear and re-add
      logger.Logger.Debug("Clearing and re-adding roles", "count", len(rolesToCleanReAdd))
      for _, roleName := range rolesToCleanReAdd {
        logger.Logger.Info("Clearing role", "roleName", roleName)

        err = ansible.ClearRole(appDir, roleName)
        if err != nil {
          logger.Logger.Error("Error clearning role", "error", err)
          os.Exit(200)
        }
        logger.Logger.Info("Successfully cleared role", "roleName", roleName)
        logger.Logger.Info("Re-adding role", "roleName", roleName)

        roleData := appSettings.ProvisionSettings.AnsibleRoles[roleName]
        if roleData.LocationType == "git" {
          logger.Logger.Debug("Role is a git role")

          err := ansible.CreateGitRole(appDir, roleName, roleData)
          if err != nil {
            logger.Logger.Error("Error recreating role", "error", err)
            os.Exit(200)
          }
        } else {
          logger.Logger.Debug("Role is a local role")

          err := ansible.CreateLocalRole(appDir, roleName, roleData)
          if err != nil {
            logger.Logger.Error("Error recreating role", "error", err)
            os.Exit(200)
          }
        }

        logger.Logger.Debug("Updating role cache")
        currentRoles.Roles[roleName] = roleData
        err = ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.Logger.Error("Error updating cache", "error", err)
          os.Exit(200)
        }
        logger.Logger.Info("Successfully cleared and re-added role", "roleName", roleName)
      }
    }
  },
}

func init() {
  RootCmd.AddCommand(roleSyncCmd)
}


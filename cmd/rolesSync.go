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
    if !machineOutput {
      fmt.Println(util.TitleText)
    }

    appDir := settings.GetAppDirPath()

    logger.LogInfo("Syncing ansible roles")

    // check for settings file
    logger.LogDebug("Checking if settings file exists")
    settingsExists, err := settings.SettingsFileExists(appDir)
    if err != nil {
      logger.LogError("Error checking if settings file exists", "error", err)
      os.Exit(120)
    }

    if !settingsExists {
      logger.LogInfo("Could not file settings file, have you initialized local-kube?")
      os.Exit(120)
    }

    // read settings file
    logger.LogInfo("Reading settings file")
    appSettings, err := settings.ReadSettingsFile(appDir)
    if err != nil {
      logger.LogError("Error reading settings file", "error", err)
    }

    // check for role cache file
    logger.LogDebug("Checking if role cache file exists")
    cacheExists, err := ansible.RoleCacheFileExists(appDir)
    if err != nil {
      logger.LogError("Error checking for role cache file", "error", err)
      os.Exit(120)
    }
    
    if !cacheExists {
      roleCache := ansible.RoleCache{
        Roles: make(map[string]settings.AnsibleRole),
      }

      logger.LogInfo("No current role cache exists, installing roles ")
      
      // Go through roles in settings file and create all the roles
      for roleName, role := range appSettings.ProvisionSettings.AnsibleRoles {
        logger.LogInfo("Processing role", "roleName", roleName)
        if role.LocationType == "git"{

          logger.LogInfo("Installing role from git", "roleName", roleName)
          logger.LogDebug("roleName", roleName, "role", role)
          err := ansible.CreateGitRole(appDir, roleName, role)
          if err != nil {
            logger.LogError("Error installing role from git", "error", err)
            
            // write out currently installed roles
            logger.LogInfo("Writing out all currently installed roles, please check config")
            err = ansible.WriteRoleCache(appDir, roleCache)
            if err != nil {
              logger.LogError("Error writing out role cache, manual intervention needed", "error", err)
              os.Exit(200)
            }
            logger.LogInfo("Partial role cache written")
            os.Exit(120)
          }

          logger.LogInfo("Adding role to cache file", "roleName", roleName)
          roleCache.Roles[roleName] = appSettings.ProvisionSettings.AnsibleRoles[roleName]

        } else if role.LocationType == "local" {

          logger.LogInfo("Installing role from the local filesystem", "roleName", roleName)
          logger.LogDebug("roleName", roleName, "role", role)
          err := ansible.CreateLocalRole(appDir, roleName, role)
          if err != nil {
            logger.LogError("Error installing role from filesystem", "error", err)

            //write out currently installed roles
            logger.LogInfo("Writing out currently installed roles, please check the config")
            err = ansible.WriteRoleCache(appDir, roleCache)
            if err != nil {
              logger.LogError("Error writing out the role cache. manual intervention needed", "error", err)
              os.Exit(200)
            }
            logger.LogInfo("Partial role cache written")
            os.Exit(120)
          }

          logger.LogInfo("Adding role to cache file", "roleName", roleName)
          roleCache.Roles[roleName] = appSettings.ProvisionSettings.AnsibleRoles[roleName]

        } else {
          // if the role location type isnt supported error out 
          // and roll back (clear all roles)
          logger.LogError("Unsupported location type", "roleName", roleName, "type", role.LocationType)
          logger.LogInfo("Please check config for the role and rerun")
          os.Exit(120)
        }

        logger.LogInfo("Writing cache to file")
        err := ansible.WriteRoleCache(appDir, roleCache)
        if err != nil {
          logger.LogError("Error writing role cache", "error", err)
          os.Exit(120)
        }
      }

    } else {
      logger.LogInfo("Role cache file exists, updating roles")

      currentRoles, err := ansible.ReadRoleCache(appDir)
      if err != nil {
        logger.LogError("Error reading role cache file", "error", err)
        os.Exit(120)
      }

      logger.LogInfo("Checking for roles to add, update or remove")
      
      rolesToAdd, 
      rolesToUpdateInPlace, 
      rolesToCleanReAdd, 
      rolesToRemove, 
      err := ansible.RoleReconcileLists(currentRoles.Roles, appSettings.ProvisionSettings.AnsibleRoles)
      
      if err != nil {
        logger.LogError("Error getting lists for role actions", "error", err)
        os.Exit(120)
      }

      logger.LogDebug("roles to add", "roles", rolesToAdd)
      logger.LogDebug("roles to remove", "roles", rolesToRemove)
      logger.LogDebug("roles to clear and then readd", "roles", rolesToCleanReAdd)
      logger.LogDebug("roles to update in place", "roles", rolesToUpdateInPlace)

      // add roles
      logger.LogDebug("Adding new roles", "count", len(rolesToAdd))
      for _, roleName := range rolesToAdd {
        logger.LogInfo("Adding role", "roleName", roleName)
        
        roleData := appSettings.ProvisionSettings.AnsibleRoles[roleName]

        if roleData.LocationType == "git" {
          logger.LogDebug("Role is git based")
          err := ansible.CreateGitRole(appDir, roleName, roleData)

          if err != nil {
            logger.LogError("Error adding git role", "roleName", roleName, "error", err)
            os.Exit(200)
          }

        } else {
          logger.LogDebug("Role is local based")
          err := ansible.CreateLocalRole(appDir, roleName, roleData)

          if err != nil {
            logger.LogError("Error adding local role", "roleName", roleName, "error", err)
            os.Exit(200)
          }
        }
        logger.LogDebug("Adding role to cache")
        currentRoles.Roles[roleName] = roleData
        err := ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.LogError("Error updating cache", "error", err)
          os.Exit(200)
        }

        logger.LogInfo("Successfully added role", "roleName", roleName)
      }

      // remove roles
      logger.LogDebug("Removing roles", "count", len(rolesToRemove))
      for _, roleName := range rolesToRemove {
        logger.LogInfo("Removing role", "roleName", roleName)

        err := ansible.ClearRole(appDir, roleName)
        if err != nil {
          logger.LogError("Error removing role, please check and re-run", "error", err)
          os.Exit(200)
        }

        logger.LogDebug("removing role from cache")
        delete(currentRoles.Roles, roleName)
        err = ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.LogError("Error updating cache", "error", err)
          os.Exit(200)
        }
        logger.LogInfo("Successfully removed role", "roleName", roleName)
      }

      // update in place
      logger.LogDebug("Updating roles in place", "count", len(rolesToUpdateInPlace))
      for _, roleName := range rolesToUpdateInPlace {
        logger.LogInfo("Updating role in place", "roleName", roleName)

        currentRole := currentRoles.Roles[roleName]
        newRole := appSettings.ProvisionSettings.AnsibleRoles[roleName]
        
        logger.LogDebug("current role data", "role", currentRole)
        logger.LogDebug("new role data", "role", newRole)
        
        err = ansible.UpdateGitRole(appDir, roleName, currentRole, newRole)
        if err != nil {
          logger.LogError("Error updating role", "error", err)
          os.Exit(200)
        }

        logger.LogDebug("Updating role cache")
        currentRoles.Roles[roleName] = newRole
        err = ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.LogError("Error updating cache", "error", err)
          os.Exit(200)
        }
        logger.LogInfo("Successfully updated role", "roleName", roleName)
      }

      // clear and re-add
      logger.LogDebug("Clearing and re-adding roles", "count", len(rolesToCleanReAdd))
      for _, roleName := range rolesToCleanReAdd {
        logger.LogInfo("Clearing role", "roleName", roleName)

        err = ansible.ClearRole(appDir, roleName)
        if err != nil {
          logger.LogError("Error clearning role", "error", err)
          os.Exit(200)
        }
        logger.LogInfo("Successfully cleared role", "roleName", roleName)
        logger.LogInfo("Re-adding role", "roleName", roleName)

        roleData := appSettings.ProvisionSettings.AnsibleRoles[roleName]
        if roleData.LocationType == "git" {
          logger.LogDebug("Role is a git role")

          err := ansible.CreateGitRole(appDir, roleName, roleData)
          if err != nil {
            logger.LogError("Error recreating role", "error", err)
            os.Exit(200)
          }
        } else {
          logger.LogDebug("Role is a local role")

          err := ansible.CreateLocalRole(appDir, roleName, roleData)
          if err != nil {
            logger.LogError("Error recreating role", "error", err)
            os.Exit(200)
          }
        }

        logger.LogDebug("Updating role cache")
        currentRoles.Roles[roleName] = roleData
        err = ansible.WriteRoleCache(appDir, currentRoles)
        if err != nil {
          logger.LogError("Error updating cache", "error", err)
          os.Exit(200)
        }
        logger.LogInfo("Successfully cleared and re-added role", "roleName", roleName)
      }
    }
  },
}

func init() {
  RootCmd.AddCommand(roleSyncCmd)
}


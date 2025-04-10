package ansible

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgutierrez1287/local-kube/logger"
	"github.com/dgutierrez1287/local-kube/settings"
	"github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/otiai10/copy"
)

/*
ClearRoles()
This will clear all the roles in the ansible-roles directory
this is used to start fresh
*/
func ClearRoles(appDir string) error {
  roleDir := filepath.Join(appDir, ansibleRoleDir)

  entries, err := os.ReadDir(roleDir)
  if err != nil {
    logger.LogError("Error reading the ansible role directory", "dir", roleDir)
    return err
  }

  for _, entry := range entries {
    if entry.IsDir() {
      dirPath := filepath.Join(roleDir, entry.Name())

      logger.LogInfo("Clearing ansible role", "name", entry.Name(), "path", dirPath)
      err := os.RemoveAll(dirPath)
      if err != nil {
        logger.LogError("Error removing role", "name", entry.Name(), "path", dirPath)
        return err
      }
      logger.LogDebug("Ansible role removed", "name", entry.Name(), "path", dirPath)
    }
  }
  return nil
}

/*
ClearRole() 
This will clear out a single role, this is used if a role is deleted,
updating a local role, or changing a role type from git to local or 
local to git
*/
func ClearRole(appDir string, roleName string) error {
  rolePath := filepath.Join(appDir, ansibleRoleDir, roleName)
  
  logger.LogInfo("Clearing ansible role", "name", roleName, "path", rolePath)
  err := os.RemoveAll(rolePath)
  
  if err != nil {
    logger.LogError("Error clearing ansible role", "name", roleName, "path", rolePath)
    return err
  }
  logger.LogDebug("Successfully cleared ansible role")
  return nil
}

/*
RoleReconcileLists()
This will return 3 lists and a error, the list are:
roles to add: roles to be added 
roles to update in place: roles that can be updated in place
roles to clean and readd: roles that have to be deleted and re-added
roles to remove: roles to be removed

These lists will be a reconcilation between the desired
state of the roles and the current state of the roles
*/
func RoleReconcileLists(currentRoles map[string]settings.AnsibleRole, desiredRoles map[string]settings.AnsibleRole) ([]string, []string, []string, []string ,error) {
  var rolesToAdd []string
  var rolesToUpdateInPlace []string
  var rolesToCleanReAdd []string
  var rolesToRemove []string

  // loop through the desired roles to get lists of roles to add 
  // roles to upgrade
  logger.LogDebug("Looping through desired roles to get add/update lists")
  for desiredRoleName, desiredRole := range desiredRoles {
    logger.LogDebug("Processing role", "roleName", desiredRoleName)

    // error check for for bad config
    if desiredRole.LocationType != "git" && desiredRole.LocationType != "local" {
      logger.LogError("Error role config is bad, not supported location type", "roleName", desiredRoleName)
      return []string{}, []string{}, []string{}, []string{}, errors.New("unsupported location type")
    }

    // if the role is not in the current role list, role is to be added
    if _, exists := currentRoles[desiredRoleName]; !exists {
      logger.LogDebug("Role is not in the current role list, to be added", "roleName", desiredRoleName)
      rolesToAdd = append(rolesToAdd, desiredRoleName)
      continue
    }

    // get the current state and desired state for some comparisons
    logger.LogDebug("Getting current and desired state for comparisons")
    currentRole := currentRoles[desiredRoleName]
    logger.LogDebug("Current role state", "name", desiredRoleName, "role", currentRole)
    logger.LogDebug("Desired role state", "name", desiredRoleName, "role", desiredRole)

    // if the role is local it will always require a 
    // clean and recopy
    if desiredRole.LocationType == "local" {
      logger.LogDebug("Role is local type, requires a clean and recopy")
      rolesToCleanReAdd = append(rolesToCleanReAdd, desiredRoleName)
    }

    // if role location type is git further checks are 
    // needed to see if the role can be updated in place
    if desiredRole.LocationType == "git" {
      logger.LogDebug("Role is git type, checking to see if update in place is possible")
      updateable := gitUpdateableRole(currentRole, desiredRole)

      if updateable {
        logger.LogDebug("Role is able to be updated in place")
        rolesToUpdateInPlace = append(rolesToUpdateInPlace, desiredRoleName)
      } else {
        logger.LogDebug("Role is not able to be updated in place")
        rolesToCleanReAdd = append(rolesToCleanReAdd, desiredRoleName)
      }
    }
  }

  // loop through the current roles to get a list of roles to 
  // remove
  logger.LogDebug("Looping through current roles to get a remove list")
  for currentRoleName := range currentRoles {
    logger.LogDebug("Processing role", "roleName", currentRoleName)

    // if the role is not in the desired role list, role to be removed
    if _, exists := desiredRoles[currentRoleName]; !exists {
      logger.LogDebug("Role is not in desired role list, to be removed", "roleName", currentRoleName)
      rolesToRemove = append(rolesToRemove, currentRoleName)
    }
  }
  return rolesToAdd, rolesToUpdateInPlace, rolesToCleanReAdd, rolesToRemove, nil
}

/*
GitUpdateableRole()
determines if a git role can be updated or needs to be cleared
and start over
*/
func gitUpdateableRole(currentRole settings.AnsibleRole, newRole settings.AnsibleRole) bool {

  logger.LogDebug("Checking location type to make sure its still pulled from git")
  if currentRole.LocationType != newRole.LocationType {
    logger.LogDebug("Location type has changed it, it now a local sourced role")
    return false
  }

  logger.LogDebug("checking if location is the same, to make sure the git repo is the same")
  if currentRole.Location != newRole.Location {
    logger.LogDebug("Git repo has changed, role will need to be cleared")
    return false
  }
  
  logger.LogDebug("Git role is updatable in place")
  return true
}

/*
UpdateGitRole()
Updates a git based ansible role in place. This will either do just a pull
if the branch or tag didn't change or fetch the new git reference and 
check out that new branch or tag
*/
func UpdateGitRole(appDir string, roleName string, currentRole settings.AnsibleRole, newRole settings.AnsibleRole) error {
  rolePath := filepath.Join(appDir, ansibleRoleDir, roleName)
  
  plumbingRef, refSpec, err :=  getReferenceName(newRole.RefType, newRole.GitRef)
  if err != nil {
    logger.LogError("Error getting git reference type")
    return err
  }

  logger.LogDebug("Opening the git repo", "repo", rolePath)
  repo, err := git.PlainOpen(rolePath) 
  
  if err != nil {
    logger.LogError("Error opening git repo", "repo", rolePath)
    return err
  }

  // if its the same reference just do a pull to get
  // up to date
  if currentRole.GitRef == newRole.GitRef {
    logger.LogDebug("Branch or tag is the same as previous, updating...")

    worktree, err := repo.Worktree()
    if err != nil {
      logger.LogError("Error getting the worktree for the repo", "repo", rolePath)
      return err
    }

    err = worktree.Pull(&git.PullOptions{
      RemoteName: "origin",
      ReferenceName: plumbingRef,
      SingleBranch: true,
      Progress: os.Stdout,
    })

    if err == git.NoErrAlreadyUpToDate {
      logger.LogInfo("Branch or tag is already up to date", "ref", newRole.GitRef)
      return nil
    } else if err != nil {
      logger.LogError("Error pulling from git")
      return err
    } else {
      logger.LogDebug("Branch or Tag updated")
      return nil
    }
  }

  // fetcb a different reference (tag or branch) and 
  // check out that reference
  logger.LogDebug("Branch or tag is different from previous, fetching new")

  err = repo.Fetch(&git.FetchOptions{
    RemoteName: "origin",
    RefSpecs: []config.RefSpec{refSpec},
    Progress: os.Stdout,
  })

  if err == git.NoErrAlreadyUpToDate {
    logger.LogDebug("Git branch or tag is already up to date")
  } else if err != nil {
    logger.LogError("Error fetching git branch or tag", "ref", newRole.GitRef)
    return err
  } 

  logger.LogDebug("Successfully fetched git branch or tag", "ref", newRole.GitRef)

  worktree, err := repo.Worktree()
  if err != nil {
    logger.LogError("Error checking out worktree")
    return err
  }

  logger.LogDebug("Checking out branch or tag", "ref", newRole.GitRef, "refType", newRole.RefType)

  err = worktree.Checkout(&git.CheckoutOptions{
    Branch: plumbingRef,
  })

  if err != nil {
    logger.LogError("Error checking out branch or tag", "ref", newRole.GitRef)
    return err
  }

  logger.LogDebug("Successfully checked out branch or tag", "ref", newRole.GitRef)
  return nil
}

/*
CreateGitRole()
Creates a new git role that is not present in ansible
roles directory
*/
func CreateGitRole(appDir string, roleName string, role settings.AnsibleRole) error {
  rolePath := filepath.Join(appDir, ansibleRoleDir, roleName)
  url := role.Location
  ref := role.GitRef

  logger.LogDebug("Getting full git reference")
  plumbingRef, _, err := getReferenceName(role.RefType, ref) 
  if err != nil {
    logger.LogError("Error getting the git reference name")
    return err
  }
  logger.LogDebug("Full git reference", "ref", plumbingRef)

  logger.LogDebug("Cloning git repo for role", "role", roleName, "refType", role.RefType, "ref", ref)
  _, err = git.PlainClone(rolePath, false, &git.CloneOptions{
    URL: url,
    ReferenceName: plumbingRef,
    SingleBranch: true,
    Depth: 1,
  })

  if err != nil {
    logger.LogError("Error cloning the git repo for role", "role", roleName, "refType", role.RefType, "ref", ref)
    return err
  }
  logger.LogDebug("Successfully checked out role", "role", roleName, "refType", role.RefType, "ref", ref)
  return nil 
}

/*
CreateLocalRole()
Creates a new local role that is not present in ansible
roles directory
*/
func CreateLocalRole(appDir string,roleName string, role settings.AnsibleRole) error {
  rolePath := filepath.Join(appDir, ansibleRoleDir, roleName)

  logger.LogDebug("Copying new role to roles directory", "role", roleName, "rolePath", rolePath, "srcPath", role.Location)
  
  err := copy.Copy(role.Location, rolePath)
  if err != nil {
    logger.LogError("Error copying role", "role", roleName, "rolePath", rolePath, "srcPath", role.Location)
    return err
  }

  logger.LogDebug("role copied successfully")
  return nil
}

/*
getReferenceName()
gets the reference name for the reference type
*/
func getReferenceName(refType string, ref string) (plumbing.ReferenceName, config.RefSpec, error) {

  var refSpec config.RefSpec
  var plumbingRef plumbing.ReferenceName

  switch refType {
  case "branch":
    logger.LogDebug("Using git branch for reference", "ref", ref)

    refSpec = config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/remotes/origin/%s", ref, ref))
    plumbingRef = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", ref))

    return plumbingRef, refSpec, nil
  case "tag":
    logger.LogDebug("Using git tag for reference", "ref", ref)

    refSpec = config.RefSpec(fmt.Sprintf("+refs/tags/%s:refs/tags/%s", ref, ref))
    plumbingRef = plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", ref))

    return plumbingRef, refSpec, nil
  default:
    logger.LogError("Error ref is not a supported type", "refType", refType, "ref", ref)
    return plumbing.ReferenceName(""), config.RefSpec(""), errors.New("unknown ref type")
  }
}


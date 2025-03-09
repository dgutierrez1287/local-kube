package settings

/*
  ProvisionSettings - Settings for machine provisioning
*/
type ProvisionSettings struct {
  AnsibleVersion string                 `json:"ansibleVersion"`       // The version of ansible to use
  AnsibleRoles map[string]AnsibleRole   `json:"ansibleRoles"`         // map of roles to use for ansible
  AnsibleCollections []string           `json:"ansibleCollections"`   // List of additional ansible collections to install 
}

/*
  AnsibleRole - An ansible role to pull to include during 
  machine provisioning
*/
type AnsibleRole struct {
  LocationType string         `json:"locationType"`     // Role location type can be git or local
  Location string             `json:"location"`         // The local location or git repo for the role
  GitBranch string            `json:"gitBranch"`        // The git branch or tag to use if the location is git
}

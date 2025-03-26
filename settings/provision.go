package settings

/*
  ProvisionSettings - Settings for machine provisioning
*/
type ProvisionSettings struct {
  AnsibleVersion string                 `json:"ansibleVersion"`       // The version of ansible to use
  AnsibleRoles map[string]AnsibleRole   `json:"ansibleRoles"`         // map of roles to download that can be used for clusters
}

/*
  AnsibleRole - An ansible role to pull to include during 
  machine provisioning
*/
type AnsibleRole struct {
  LocationType string         `json:"locationType"`     // Role location type can be git or local
  Location string             `json:"location"`         // The local location or git repo for the role
  RefType string              `json:"gitRefType"`       // (branch or tag) the git refence to use
  GitRef string               `json:"gitRef"`           // The git branch or tag to use 
}


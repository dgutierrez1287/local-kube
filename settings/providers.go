package settings

/*
 Provider - A provider for vagrant to use to create a
 cluster
*/
type Provider struct {
  /*
  Common settings - These are settings common to all providers
  */ 
  ProviderType string   `json:"providerType"`    // The type of the provider

  /*
  Semi-Common settings - These are settings that are used by
  several different providers
  */
  BoxName string        `json:"boxName"`         // vagrant box name for the provider

  /*
  VmWare Fusion/Workstation - These are settings for vmware
  fusion/workstation provider only
  */
  VmNet string          `json:"vmNet"`          // The vmnet for the cluster to use
  
}

package settings

import (
	"errors"

	"github.com/dgutierrez1287/local-kube/logger"
)

/*
  ClusterFeatures - feature settings for a cluster that is
  using autoconfigure to use the default role and simplied
  configuration
*/ 
type ClusterFeatures struct {
  // kube version
  KubeVersion string                `json:"kubeVersion,omitempty"`                // The kube version to use 

  // CNI controller           
  CniController string              `json:"cniController,omitempty"`              // The Cni Controller
  CniControllerVersion string       `json:"cniControllerVersion,omitempty"`       // The Cni Controller Version
  ManagedCniController bool         `json:"managedCniController,omitempty"`       // If Cni Controller should be installed

  //Cilium specific 
  CiliumCliVersion string          `json:"cillumCliVersion,omitempty"`           // the cillum cli version

  // Ingress Controller
  IngressController string          `json:"ingressController,omitempty"`          // The ingress controller
  
  // storage Controller
  StorageController string          `json:"storageController,omitempty"`          // The storage controller
  StorageControllerVersion string   `json:"storageControllerVersion,omitempty"`   // The Version of the storage controller
  ManagedStorageController bool     `json:"managedStorageController,omitempty"`   // If the storage controller should be installed

  // kubeVip
  KubeVipEnable bool                `json:"kubeVipEnable,omitempty"`              // Enable KubeVip
  KubeVipVersion string             `json:"kubeVipVersion,omitempty"`             // KubeVip Version

  // other settings
  DisableDefaultMetrics bool        `json:"disableDefaultMetrics,omitempty"`      // Disable default cluster metrics
}

var featuresDefaults = ClusterFeatures {
  KubeVersion: "1.31.4",
  CniController: "flannel",
  IngressController: "native-traefik",
  StorageController: "local-storage",
}

var CiliumDefaultVersion = "1.16.4"
var CiliumCliDefaultVersion = "0.16.22"
var CalicoDefaultVersion = "3.25.0"
var KubevipDefaultVersion = "0.5.0"
var LonghornDefaultVersion = "1.8.0"

func (features *ClusterFeatures) SetDefaults(clusterType string, vip string) error {

  // error checking
  if clusterType == "ha" && !features.KubeVipEnable {
    logger.LogError("Error you should enable KubeVip, for ha clusters")
    return errors.New("ha cluster, but kubevip not enabled")
  }

  if features.KubeVipEnable && vip == "" {
    logger.LogError("Error the vip cannot be empty with kubevip enabled")
    return errors.New("kubevip enabled but no vip provided")
  }

  // Kube Version defaults
  if features.KubeVersion == "" {
    features.KubeVersion = featuresDefaults.KubeVersion
    logger.LogDebug("No Kubenetes version supplied, setting default", "version", featuresDefaults.KubeVersion)
  }

  //KubeVip defaults
  if features.KubeVipEnable && features.KubeVipVersion == "" {
    logger.LogDebug("Kubevip enabled but no version supplied, using default", "version", KubevipDefaultVersion)
    features.KubeVipVersion = KubevipDefaultVersion
  }

  // Cni Controller defaults
  if features.CniController == "" {
    features.CniController = featuresDefaults.CniController
    logger.LogDebug("No cni controller supplied, setting default", "cni", featuresDefaults.CniController)

  } else {
    if features.CniController == "cilium" {
      logger.LogDebug("Cni controller supplied", "controller", features.CniController, "managed", features.ManagedCniController)

      if features.CniControllerVersion == "" {
        logger.LogDebug("Cni Controller version not set, using default", "controller", features.CniController, "version", CiliumDefaultVersion)
        features.CniControllerVersion = CiliumDefaultVersion
      }     

      if features.CiliumCliVersion == "" {
        logger.LogDebug("Cni is Cilium and cli version is not set, using default", "cliVersion", CiliumCliDefaultVersion)
        features.CiliumCliVersion = CiliumCliDefaultVersion
      }

    } else if features.CniController == "calico" {
      logger.LogDebug("Cni controller supplied", "controller", features.CniController, "managed", features.ManagedCniController)

      if features.CniControllerVersion == "" {
        logger.LogDebug("Cni Controller version not set, using default", "controller", features.CniController, "version", CalicoDefaultVersion)
        features.CniControllerVersion = CalicoDefaultVersion
      }

    } else {
      logger.LogError("Error cni controller is not supported", "controller", features.CniController)
      return errors.New("cni controller is not supported")
    }
  }

  // Ingress Controller defaults
  if features.IngressController == "" {
    features.IngressController = featuresDefaults.IngressController
    logger.LogDebug("No ingress controller supplied, setting default", "ingress", featuresDefaults.IngressController)
  }

  // Storage Controller defaults
  if features.StorageController == "" {
    features.StorageController = featuresDefaults.StorageController
    logger.LogDebug("No storage controller supplied, setting default", "storage", featuresDefaults.StorageController)
  } else {
    logger.LogDebug("Storage controller supplied", "controller", features.StorageController, "managed", features.ManagedStorageController)

    if features.StorageControllerVersion == "" && features.StorageController == "longhorn"{
      logger.LogDebug("No storage controller version supplied settings default")
      features.StorageControllerVersion = LonghornDefaultVersion
    }
  }
  return nil
}

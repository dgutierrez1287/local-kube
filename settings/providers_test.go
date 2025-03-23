package settings

import (
	"testing"
  "encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestProviderInitialization(t *testing.T) {
	// Test case where a provider is initialized with all values set
	provider := Provider{
		ProviderType: "VMware",
		BoxName:      "ubuntu-20.04",
		VmNet:        "vmnet0",
	}

	// Validate the fields of the provider
	assert.Equal(t, "VMware", provider.ProviderType)
	assert.Equal(t, "ubuntu-20.04", provider.BoxName)
	assert.Equal(t, "vmnet0", provider.VmNet)
}

func TestProvider_DefaultValues(t *testing.T) {
	// Test case where a provider is initialized with zero values
	provider := Provider{}

	// Validate the default values of the fields
	assert.Equal(t, "", provider.ProviderType)
	assert.Equal(t, "", provider.BoxName)
	assert.Equal(t, "", provider.VmNet)
}

func TestProvider_JsonSerialization(t *testing.T) {
	// Test case for JSON serialization and deserialization
	provider := Provider{
		ProviderType: "VMware",
		BoxName:      "ubuntu-20.04",
		VmNet:        "vmnet0",
	}

	// Serialize to JSON
	data, err := json.Marshal(provider)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedProvider Provider
	err = json.Unmarshal(data, &deserializedProvider)
	assert.NoError(t, err)

	// Assert the values after unmarshaling
	assert.Equal(t, provider.ProviderType, deserializedProvider.ProviderType)
	assert.Equal(t, provider.BoxName, deserializedProvider.BoxName)
	assert.Equal(t, provider.VmNet, deserializedProvider.VmNet)
}

func TestProvider_EmptyJsonSerialization(t *testing.T) {
	// Test case for serializing an empty provider
	provider := Provider{}

	// Serialize to JSON
	data, err := json.Marshal(provider)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedProvider Provider
	err = json.Unmarshal(data, &deserializedProvider)
	assert.NoError(t, err)

	// Assert the values after unmarshaling (they should be empty)
	assert.Equal(t, provider.ProviderType, deserializedProvider.ProviderType)
	assert.Equal(t, provider.BoxName, deserializedProvider.BoxName)
	assert.Equal(t, provider.VmNet, deserializedProvider.VmNet)
}

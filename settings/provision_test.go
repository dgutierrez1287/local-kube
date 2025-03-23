package settings

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvisionSettingsInitialization(t *testing.T) {
	// Test case where ProvisionSettings is initialized with values
	provisionSettings := ProvisionSettings{
		AnsibleVersion: "2.9.1",
		AnsibleRoles: map[string]AnsibleRole{
			"role1": {
				LocationType: "git",
				Location:     "https://github.com/example/role1",
				RefType:      "branch",
				GitRef:       "main",
			},
		},
		AnsibleCollections: []string{"collection1", "collection2"},
	}

	// Validate the fields of the ProvisionSettings struct
	assert.Equal(t, "2.9.1", provisionSettings.AnsibleVersion)
	assert.Equal(t, 1, len(provisionSettings.AnsibleRoles))
	assert.Equal(t, "git", provisionSettings.AnsibleRoles["role1"].LocationType)
	assert.Equal(t, "https://github.com/example/role1", provisionSettings.AnsibleRoles["role1"].Location)
	assert.Equal(t, "branch", provisionSettings.AnsibleRoles["role1"].RefType)
	assert.Equal(t, "main", provisionSettings.AnsibleRoles["role1"].GitRef)
	assert.Equal(t, 2, len(provisionSettings.AnsibleCollections))
}

func TestProvisionSettingsDefaultValues(t *testing.T) {
	// Test case where ProvisionSettings is initialized with default (empty) values
	provisionSettings := ProvisionSettings{}

	// Validate the default values of the fields
	assert.Equal(t, "", provisionSettings.AnsibleVersion)
	assert.Empty(t, provisionSettings.AnsibleRoles)
	assert.Empty(t, provisionSettings.AnsibleCollections)
}

func TestProvisionSettingsJsonSerialization(t *testing.T) {
	// Test case for JSON serialization and deserialization of ProvisionSettings
	provisionSettings := ProvisionSettings{
		AnsibleVersion: "2.9.1",
		AnsibleRoles: map[string]AnsibleRole{
			"role1": {
				LocationType: "git",
				Location:     "https://github.com/example/role1",
				RefType:      "branch",
				GitRef:       "main",
			},
		},
		AnsibleCollections: []string{"collection1", "collection2"},
	}

	// Serialize to JSON
	data, err := json.Marshal(provisionSettings)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedProvisionSettings ProvisionSettings
	err = json.Unmarshal(data, &deserializedProvisionSettings)
	assert.NoError(t, err)

	// Assert the values after unmarshaling
	assert.Equal(t, provisionSettings.AnsibleVersion, deserializedProvisionSettings.AnsibleVersion)
	assert.Equal(t, provisionSettings.AnsibleRoles, deserializedProvisionSettings.AnsibleRoles)
	assert.Equal(t, provisionSettings.AnsibleCollections, deserializedProvisionSettings.AnsibleCollections)
}

func TestProvisionSettingsEmptyJsonSerialization(t *testing.T) {
	// Test case for serializing an empty ProvisionSettings struct
	provisionSettings := ProvisionSettings{}

	// Serialize to JSON
	data, err := json.Marshal(provisionSettings)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedProvisionSettings ProvisionSettings
	err = json.Unmarshal(data, &deserializedProvisionSettings)
	assert.NoError(t, err)

	// Assert the values after unmarshaling (they should be empty)
	assert.Equal(t, provisionSettings.AnsibleVersion, deserializedProvisionSettings.AnsibleVersion)
	assert.Empty(t, deserializedProvisionSettings.AnsibleRoles)
	assert.Empty(t, deserializedProvisionSettings.AnsibleCollections)
}

func TestAnsibleRoleInitialization(t *testing.T) {
	// Test case where an AnsibleRole is initialized with values
	ansibleRole := AnsibleRole{
		LocationType: "git",
		Location:     "https://github.com/example/role1",
		RefType:      "branch",
		GitRef:       "main",
	}

	// Validate the fields of the AnsibleRole struct
	assert.Equal(t, "git", ansibleRole.LocationType)
	assert.Equal(t, "https://github.com/example/role1", ansibleRole.Location)
	assert.Equal(t, "branch", ansibleRole.RefType)
	assert.Equal(t, "main", ansibleRole.GitRef)
}

func TestAnsibleRoleDefaultValues(t *testing.T) {
	// Test case where AnsibleRole is initialized with default (empty) values
	ansibleRole := AnsibleRole{}

	// Validate the default values of the fields
	assert.Equal(t, "", ansibleRole.LocationType)
	assert.Equal(t, "", ansibleRole.Location)
	assert.Equal(t, "", ansibleRole.RefType)
	assert.Equal(t, "", ansibleRole.GitRef)
}

func TestAnsibleRoleJsonSerialization(t *testing.T) {
	// Test case for JSON serialization and deserialization of AnsibleRole
	ansibleRole := AnsibleRole{
		LocationType: "git",
		Location:     "https://github.com/example/role1",
		RefType:      "branch",
		GitRef:       "main",
	}

	// Serialize to JSON
	data, err := json.Marshal(ansibleRole)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedAnsibleRole AnsibleRole
	err = json.Unmarshal(data, &deserializedAnsibleRole)
	assert.NoError(t, err)

	// Assert the values after unmarshaling
	assert.Equal(t, ansibleRole.LocationType, deserializedAnsibleRole.LocationType)
	assert.Equal(t, ansibleRole.Location, deserializedAnsibleRole.Location)
	assert.Equal(t, ansibleRole.RefType, deserializedAnsibleRole.RefType)
	assert.Equal(t, ansibleRole.GitRef, deserializedAnsibleRole.GitRef)
}

func TestAnsibleRoleEmptyJsonSerialization(t *testing.T) {
	// Test case for serializing an empty AnsibleRole struct
	ansibleRole := AnsibleRole{}

	// Serialize to JSON
	data, err := json.Marshal(ansibleRole)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedAnsibleRole AnsibleRole
	err = json.Unmarshal(data, &deserializedAnsibleRole)
	assert.NoError(t, err)

	// Assert the values after unmarshaling (they should be empty)
	assert.Equal(t, ansibleRole.LocationType, deserializedAnsibleRole.LocationType)
	assert.Equal(t, ansibleRole.Location, deserializedAnsibleRole.Location)
	assert.Equal(t, ansibleRole.RefType, deserializedAnsibleRole.RefType)
	assert.Equal(t, ansibleRole.GitRef, deserializedAnsibleRole.GitRef)
}

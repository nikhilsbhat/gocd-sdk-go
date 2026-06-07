package gocd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectKeys_ShouldConvertAllExtensionSettingKeys(t *testing.T) {
	newSetting := func(key string) *PluginSettingAttribute {
		return &PluginSettingAttribute{
			Configurations: []*PluginConfiguration{
				{Key: key},
			},
		}
	}

	plugin := &Plugin{
		Extensions: []PluginAttributes{
			{
				AuthConfigSettings:          newSetting("AuthConfigID"),
				ArtifactConfigSettings:      newSetting("ArtifactStoreID"),
				ElasticAgentProfileSettings: newSetting("ElasticProfileID"),
				FetchArtifactSettings:       newSetting("FetchArtifactID"),
				ClusterProfileSettings:      newSetting("ClusterProfileID"),
				PluginSettings:              newSetting("PluginSettingID"),
				PackageSettings:             newSetting("PackageID"),
				RepositorySettings:          newSetting("RepositoryID"),
				ScmSettings:                 newSetting("ScmID"),
				StoreConfigSettings:         newSetting("StoreConfigID"),
				SecretConfigSettings:        newSetting("SecretConfigID"),
				RoleSettings:                newSetting("RoleID"),
				TaskSettings:                newSetting("TaskID"),
			},
		},
	}

	correctKeys(plugin)

	extension := plugin.Extensions[0]
	assert.Equal(t, "auth_config_id", extension.AuthConfigSettings.Configurations[0].Key)
	assert.Equal(t, "artifact_store_id", extension.ArtifactConfigSettings.Configurations[0].Key)
	assert.Equal(t, "elastic_profile_id", extension.ElasticAgentProfileSettings.Configurations[0].Key)
	assert.Equal(t, "fetch_artifact_id", extension.FetchArtifactSettings.Configurations[0].Key)
	assert.Equal(t, "cluster_profile_id", extension.ClusterProfileSettings.Configurations[0].Key)
	assert.Equal(t, "plugin_setting_id", extension.PluginSettings.Configurations[0].Key)
	assert.Equal(t, "package_id", extension.PackageSettings.Configurations[0].Key)
	assert.Equal(t, "repository_id", extension.RepositorySettings.Configurations[0].Key)
	assert.Equal(t, "scm_id", extension.ScmSettings.Configurations[0].Key)
	assert.Equal(t, "store_config_id", extension.StoreConfigSettings.Configurations[0].Key)
	assert.Equal(t, "secret_config_id", extension.SecretConfigSettings.Configurations[0].Key)
	assert.Equal(t, "role_id", extension.RoleSettings.Configurations[0].Key)
	assert.Equal(t, "task_id", extension.TaskSettings.Configurations[0].Key)
}

package gocd

import "encoding/xml"

const (
	defaultRetryCount    = 5
	defaultRetryWaitTime = 5
)

// AgentsConfig holds information of all agent of GoCD.
type AgentsConfig struct {
	Config Agents `json:"_embedded,omitempty"`
}

// Agents holds information of all agent of GoCD.
type Agents struct {
	Config []Agent `json:"agents,omitempty"`
}

// Agent holds information of a particular agent.
type Agent struct {
	IPAddress          string        `json:"ip_address,omitempty"`
	Name               string        `json:"hostname,omitempty"`
	ID                 string        `json:"uuid,omitempty"`
	Version            string        `json:"agent_version,omitempty"`
	CurrentState       string        `json:"agent_state,omitempty"`
	OS                 string        `json:"operating_system,omitempty"`
	ConfigState        string        `json:"agent_config_state,omitempty"`
	Sandbox            string        `json:"sandbox,omitempty"`
	DiskSpaceAvailable interface{}   `json:"free_space,omitempty"`
	Resources          []string      `json:"resources,omitempty"`
	Environments       []interface{} `json:"environments,omitempty"`
	Operations         struct {
		Resources    AddRemoves `json:"resources,omitempty"`
		Environments AddRemoves `json:"environments,omitempty"`
	} `json:"operations,omitempty"`
	UUIDS []string `json:"uuids,omitempty"`
}

type AddRemoves struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// ServerVersion holds version information GoCd server.
type ServerVersion struct {
	Version     string `json:"version,omitempty"`
	GitSha      string `json:"git_sha,omitempty"`
	FullVersion string `json:"full_version,omitempty"`
	CommitURL   string `json:"commit_url,omitempty"`
}

// ServerHealth holds information of GoCD server health.
type ServerHealth struct {
	Level   string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
	Time    string `json:"time,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

// ConfigRepoConfig holds information of all config-repos present in GoCD.
type ConfigRepoConfig struct {
	ConfigRepos ConfigRepos `json:"_embedded,omitempty"`
}

// ConfigRepos holds information of all config-repos present in GoCD.
type ConfigRepos struct {
	ConfigRepos []ConfigRepo `json:"config_repos,omitempty"`
}

// ConfigRepo holds information of the specified config-repo.
type ConfigRepo struct {
	PluginID      string                   `json:"plugin_id"`
	ID            string                   `json:"config_repos,omitempty"`
	Material      Material                 `json:"material,omitempty"`
	Configuration []map[string]interface{} `json:"configuration,omitempty"`
	Rules         []map[string]interface{} `json:"rules,omitempty"`
	ETAG          string
}

// PipelineGroupsConfig holds information on the various pipeline groups present in GoCD.
type PipelineGroupsConfig struct {
	PipelineGroups PipelineGroups `json:"_embedded,omitempty"`
}

// PipelineGroups holds information on the various pipeline groups present in GoCD.
type PipelineGroups struct {
	PipelineGroups []PipelineGroup `json:"groups,omitempty"`
}

// PipelineGroup holds information of a specific pipeline group instance.
type PipelineGroup struct {
	Name          string                 `json:"name,omitempty"`
	PipelineCount int                    `json:"pipeline_count,omitempty"`
	Pipelines     []Pipeline             `json:"pipelines,omitempty"`
	Authorization map[string]interface{} `json:"authorization,omitempty"`
	ETAG          string
}

// SystemAdmins holds information of the system admins present.
type SystemAdmins struct {
	Roles []string `json:"roles,omitempty"`
	Users []string `json:"users,omitempty"`
	ETAG  string
}

// BackupConfig holds information of the backup configured.
type BackupConfig struct {
	EmailOnSuccess   bool   `json:"email_on_success,omitempty"`
	EmailOnFailure   bool   `json:"email_on_failure,omitempty"`
	Schedule         string `json:"schedule,omitempty"`
	PostBackupScript string `json:"post_backup_script,omitempty"`
}

// PipelineSize holds information of the pipeline size.
type PipelineSize struct {
	Size float64
	Type string
}

// PipelinesInfo holds information of list of pipelines.
type PipelinesInfo struct {
	XMLName xml.Name `xml:"pipelines"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Pipeline []struct {
		Href string `xml:"href,attr"`
	} `xml:"pipeline"`
}

// PipelineState holds information of the latest state of pipeline.
type PipelineState struct {
	Name        string `json:"name,omitempty"`
	Paused      bool   `json:"paused,omitempty"`
	Locked      bool   `json:"locked,omitempty"`
	Schedulable bool   `json:"schedulable,omitempty"`
	PausedBy    string `json:"paused_by,omitempty"`
	PausedCause string `json:"paused_cause,omitempty"`
}

// Pipelines holds information of the pipelines present in GoCD.
type Pipelines struct {
	Pipelines []Pipeline `json:"pipelines,omitempty"`
}

// Pipeline holds information of a specific pipeline instance.
type Pipeline struct {
	Name string `json:"name,omitempty"`
}

// EnvironmentInfo holds information of all environments present in GoCD.
type EnvironmentInfo struct {
	Environments Environments `json:"_embedded,omitempty"`
}

// Environments holds information of all environments present in GoCD.
type Environments struct {
	Environments []Environment `json:"environments,omitempty"`
}

// Environment holds information of a specific environment present in GoCD.
type Environment struct {
	Name      string     `json:"name,omitempty"`
	Pipelines []Pipeline `json:"pipelines,omitempty"`
	EnvVars   []struct {
		Name           string `json:"name,omitempty"`
		Value          string `json:"value,omitempty"`
		EncryptedValue string `json:"encrypted_value,omitempty"`
		Secure         bool   `json:"secure,omitempty"`
	} `json:"environment_variables,omitempty"`
	ETAG string
}

// PatchEnvironment holds information that is handy while patching GoCD environment.
type PatchEnvironment struct {
	Name      string `json:"name"`
	Pipelines struct {
		Add    []string `json:"add,omitempty"`
		Remove []string `json:"remove,omitempty"`
	} `json:"pipelines,omitempty"`
	EnvVars struct {
		Add []struct {
			Name  string `json:"name,omitempty"`
			Value string `json:"value,omitempty"`
		} `json:"add,omitempty"`
		Remove []string `json:"remove,omitempty"`
	} `json:"environment_variables,omitempty"`
}

// VersionInfo holds version information of GoCD server.
type VersionInfo struct {
	Version     string `json:"version,omitempty"`
	FullVersion string `json:"full_version,omitempty"`
	GitSHA      string `json:"git_sha,omitempty"`
}

// AgentJobHistory holds information of pipeline run history of all GoCD agents.
type AgentJobHistory struct {
	Jobs       []JobRunHistory `json:"jobs,omitempty"`
	Pagination Pagination      `json:"pagination"`
}

// JobRunHistory holds information of pipeline run history of a specific GoCD agent.
type JobRunHistory struct {
	Name            string `json:"pipeline_name,omitempty"`
	JobName         string `json:"job_name,omitempty"`
	StageName       string `json:"stage_name,omitempty"`
	StageCounter    int64  `json:"stage_counter,string,omitempty"`
	PipelineCounter int64  `json:"pipeline_counter,omitempty"`
	Result          string `json:"result,omitempty"`
}

// Pagination holds information which is helpful in paginating the results of job run history.
type Pagination struct {
	PageSize int64 `json:"page_size,omitempty"`
	Offset   int64 `json:"offset,omitempty"`
	Total    int64 `json:"total,omitempty"`
}

// Maintenance holds latest information available in server about maintenance mode.
type Maintenance struct {
	MaintenanceInfo struct {
		Enabled  bool `json:"is_maintenance_mode,omitempty"`
		Metadata struct {
			UpdatedBy string `json:"updated_by,omitempty"`
			UpdatedOn string `json:"updated_on,omitempty"`
		} `json:"metadata,omitempty"`
	} `json:"_embedded,omitempty"`
}

// Encrypted holds the encrypted value of the passed plain text.
type Encrypted struct {
	EncryptedValue string `json:"encrypted_value,omitempty"`
}

// ArtifactInfo holds the latest information of the artifacts.
type ArtifactInfo struct {
	ArtifactsDir  string `json:"artifacts_dir,omitempty"`
	PurgeSettings struct {
		PurgeStartDiskSpace float64 `json:"purge_start_disk_space,omitempty"`
		PurgeUptoDiskSpace  float64 `json:"purge_upto_disk_space,omitempty"`
	} `json:"purge_settings,omitempty"`
	ETAG string
}

// Schedule holds config of the pipeline that needs to be scheduled.
type Schedule struct {
	EnvVars        []map[string]interface{} `json:"environment_variables,omitempty"`
	Materials      []map[string]interface{} `json:"materials,omitempty"`
	UpdateMaterial bool                     `json:"update_materials_before_scheduling,omitempty"`
}

// AuthConfigs holds information of multiple authorization configurations.
type AuthConfigs struct {
	Config struct {
		AuthConfigs []CommonConfig `json:"auth_configs"`
	} `json:"_embedded,omitempty"`
}

// PluginConfiguration holds information of the various plugin properties.
type PluginConfiguration struct {
	Key            string `json:"key,omitempty"`
	Value          string `json:"value,omitempty"`
	EncryptedValue string `json:"encrypted_value,omitempty"`
}

// SiteURLConfig holds information of the site url of GoCD.
type SiteURLConfig struct {
	SiteURL       string `json:"site_url,omitempty"`
	SecureSiteURL string `json:"secure_site_url,omitempty"`
}

// MailServerConfig holds information required for GoCD mail-server configuration.
type MailServerConfig struct {
	Hostname          string `json:"hostname,omitempty"`
	Port              int64  `json:"port,omitempty"`
	UserName          string `json:"username,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	TLS               bool   `json:"tls,omitempty"`
	SenderEmail       string `json:"sender_email,omitempty"`
	AdminEmail        string `json:"admin_email,omitempty"`
}

// PluginSettings holds information of plugin settings of GoCD.
type PluginSettings struct {
	ID            string                `json:"plugin_id,omitempty"`
	Configuration []PluginConfiguration `json:"configuration,omitempty"`
	ETAG          string
}

// PipelineObject holds information required to comment/get/history of pipeline or instance of the same.
type PipelineObject struct {
	Name    string
	Counter int
	Message string
}

// PipelineHistory holds information of the pipeline history that also helps in paginating the responses.
type PipelineHistory struct {
	Links     map[string]interface{}   `json:"_links,omitempty"`
	Pipelines []map[string]interface{} `json:"pipelines,omitempty"`
}

// ArtifactStoresConfigs holds information of the specified artifact-stores/cluster-profiles/agent-profiles.
type ArtifactStoresConfigs struct {
	ArtifactStoresConfigs ArtifactStoresConfig `json:"_embedded,omitempty"`
}

// ArtifactStoresConfig holds information of all config-repos present in GoCD.
type ArtifactStoresConfig struct {
	CommonConfigs []CommonConfig `json:"artifact_stores,omitempty"`
	ETAG          string
}

// ProfilesConfigs holds information of the specified artifact-stores/cluster-profiles/agent-profiles.
type ProfilesConfigs struct {
	ProfilesConfigs ProfilesConfig `json:"_embedded,omitempty"`
}

// ProfilesConfig holds information of all config-repos present in GoCD.
type ProfilesConfig struct {
	CommonConfigs []CommonConfig `json:"profiles,omitempty"`
	ETAG          string
}

// SecretsConfigs holds information of all secret configs present in GoCD.
type SecretsConfigs struct {
	SecretsConfigs SecretsConfig `json:"_embedded,omitempty"`
}

// SecretsConfig holds information of a specified secret config present in GoCD.
type SecretsConfig struct {
	CommonConfigs []CommonConfig `json:"secret_configs,omitempty"`
	ETAG          string
}

// CommonConfig holds information of the specified artifact store.
type CommonConfig struct {
	ID                  string                `json:"id,omitempty"`
	Name                string                `json:"name,omitempty"`
	PluginID            string                `json:"plugin_id,omitempty"`
	Description         string                `json:"description,omitempty"`
	ClusterProfileID    string                `json:"cluster_profile_id,omitempty"`
	AllowOnlyKnownUsers bool                  `json:"allow_only_known_users_to_login,omitempty"`
	Properties          []PluginConfiguration `json:"properties,omitempty"`
	Rules               []map[string]string   `json:"rules,omitempty"`
	ETAG                string                `json:"etag,omitempty"`
}

// PackageRepositories holds information of all package repositories present in GoCD.
type PackageRepositories struct {
	Repositories struct {
		PackageRepositories []PackageRepository `json:"package_repositories,omitempty"`
	} `json:"_embedded"`
}

// PackageRepository holds information of the specified package repository.
type PackageRepository struct {
	ID             string                `json:"repo_id,omitempty"`
	Name           string                `json:"name,omitempty"`
	PluginMetaData map[string]string     `json:"plugin_metadata,omitempty"`
	Configuration  []PluginConfiguration `json:"configuration,omitempty"`
	Packages       struct {
		Packages []CommonConfig `json:"packages,omitempty"`
	} `json:"_embedded,omitempty"`
	ETAG string
}

// Packages holds information of all packages present in GoCD.
type Packages struct {
	Packages struct {
		Packages []Package `json:"packages,omitempty"`
	} `json:"_embedded"`
}

// Package holds information of the specified packages of the package repository.
type Package struct {
	CommonConfig
	AutoUpdate    bool                  `json:"auto_update,omitempty"`
	PackageRepos  CommonConfig          `json:"package_repo,omitempty"`
	Configuration []PluginConfiguration `json:"configuration,omitempty"`
	ETAG          string                `json:"etag,omitempty"`
}

// Materials holds information of all material type present in GoCD.
type Materials struct {
	Materials struct {
		Materials []Material `json:"materials,omitempty"`
	} `json:"_embedded,omitempty"`
}

// Material holds information of a particular material type present in GoCD.
type Material struct {
	Type        string    `json:"type,omitempty"`
	Fingerprint string    `json:"fingerprint,omitempty"`
	Attributes  Attribute `json:"attributes,omitempty"`
}

// Attribute holds information of material type present in GoCD.
type Attribute struct {
	URL                 string `json:"url,omitempty"`
	Username            string `json:"username,omitempty"`
	Password            string `json:"password,omitempty"`
	EncryptedPassword   string `json:"encrypted_password,omitempty"`
	Branch              string `json:"branch,omitempty"`
	AutoUpdate          bool   `json:"auto_update,omitempty"`
	CheckExternals      bool   `json:"check_externals,omitempty"`
	UseTickets          bool   `json:"use_tickets,omitempty"`
	View                string `json:"view,omitempty"`
	Port                string `json:"port,omitempty"`
	ProjectPath         string `json:"project_path,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Ref                 string `json:"ref,omitempty"`
	Name                string `json:"name,omitempty"`
	Stage               string `json:"stage,omitempty"`
	Pipeline            string `json:"pipeline,omitempty"`
	IgnoreForScheduling bool   `json:"ignore_for_scheduling,omitempty"`
	Destination         string `json:"destination,omitempty"`
	InvertFilter        bool   `json:"invert_filter,omitempty"`
	Filter              struct {
		Ignore []string `json:"ignore,omitempty"`
	} `json:"filter,omitempty"`
}

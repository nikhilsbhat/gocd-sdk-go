package gocd

import "encoding/xml"

const (
	defaultRetryCount    = 5
	defaultRetryWaitTime = 5
)

// AgentsConfig holds information of all agent of GoCD.
type AgentsConfig struct {
	Config Agents `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// Agents holds information of all agent of GoCD.
type Agents struct {
	Config []Agent `json:"agents,omitempty" yaml:"agents,omitempty"`
}

// Agent holds information of a particular agent.
type Agent struct {
	ID                 string      `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Name               string      `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	ElasticAgentID     string      `json:"elastic_agent_id,omitempty" yaml:"elastic_agent_id,omitempty"`
	ElasticPluginID    string      `json:"elastic_plugin_id,omitempty" yaml:"elastic_plugin_id,omitempty"`
	IPAddress          string      `json:"ip_address,omitempty" yaml:"ip_address,omitempty"`
	Sandbox            string      `json:"sandbox,omitempty" yaml:"sandbox,omitempty"`
	OS                 string      `json:"operating_system,omitempty" yaml:"operating_system,omitempty"`
	DiskSpaceAvailable interface{} `json:"free_space,omitempty" yaml:"free_space,omitempty"`
	ConfigState        string      `json:"agent_config_state,omitempty" yaml:"agent_config_state,omitempty"`
	CurrentState       string      `json:"agent_state,omitempty" yaml:"agent_state,omitempty"`
	Version            string      `json:"agent_version,omitempty" yaml:"agent_version,omitempty"`
	Resources          []string    `json:"resources,omitempty" yaml:"resources,omitempty"`
	Environments       any         `json:"environments,omitempty" yaml:"environments,omitempty"`
	BuildState         string      `json:"build_state,omitempty" yaml:"build_state,omitempty"`
	BuildDetails       BuildInfo   `json:"build_details,omitempty" yaml:"build_details,omitempty"`
	Operations         Operations  `json:"operations,omitempty" yaml:"operations,omitempty"`
	UUIDS              []string    `json:"uuids,omitempty" yaml:"uuids,omitempty"`
}

type BuildInfo struct {
	Pipeline string `json:"pipeline_name,omitempty" yaml:"pipeline_name,omitempty"`
	Stage    string `json:"stage_name,omitempty" yaml:"stage_name,omitempty"`
	Job      string `json:"job_name,omitempty" yaml:"job_name,omitempty"`
}

// Operations holds information of the operations to be performed on GoCD agent.
type Operations struct {
	Resources    AddRemoves `json:"resources,omitempty" yaml:"resources,omitempty"`
	Environments AddRemoves `json:"environments,omitempty" yaml:"environments,omitempty"`
}

type AddRemoves struct {
	Add    []string `json:"add,omitempty" yaml:"add,omitempty"`
	Remove []string `json:"remove,omitempty" yaml:"remove,omitempty"`
}

// ServerVersion holds version information GoCd server.
type ServerVersion struct {
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	GitSha      string `json:"git_sha,omitempty" yaml:"git_sha,omitempty"`
	FullVersion string `json:"full_version,omitempty" yaml:"full_version,omitempty"`
	CommitURL   string `json:"commit_url,omitempty" yaml:"commit_url,omitempty"`
}

// ServerHealth holds information of GoCD server health.
type ServerHealth struct {
	Level   string `json:"level,omitempty" yaml:"level,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Time    string `json:"time,omitempty" yaml:"time,omitempty"`
	Detail  string `json:"detail,omitempty" yaml:"detail,omitempty"`
}

// ConfigRepoConfig holds information of all config-repos present in GoCD.
type ConfigRepoConfig struct {
	ConfigRepos ConfigRepos `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// ConfigRepos holds information of all config-repos present in GoCD.
type ConfigRepos struct {
	ConfigRepos []ConfigRepo `json:"config_repos,omitempty" yaml:"config_repos,omitempty"`
}

// ConfigRepo holds information of the specified config-repo.
type ConfigRepo struct {
	PluginID            string                `json:"plugin_id,omitempty" yaml:"plugin_id,omitempty"`
	ID                  string                `json:"id,omitempty" yaml:"id,omitempty"`
	Material            Material              `json:"material,omitempty" yaml:"material,omitempty"`
	Configuration       []PluginConfiguration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
	Rules               []map[string]string   `json:"rules,omitempty" yaml:"rules,omitempty"`
	ETAG                string                `json:"etag,omitempty" yaml:"etag,omitempty"`
	Environments        []Environment         `json:"environments,omitempty" yaml:"environments,omitempty"`
	Groups              []PipelineGroup       `json:"groups,omitempty" yaml:"groups,omitempty"`
	ConfigRepoParseInfo ConfigRepoParseInfo   `json:"parse_info,omitempty" yaml:"parse_info,omitempty"`
}

type ConfigRepoParseInfo struct {
	LatestParsedModification map[string]interface{} `json:"latest_parsed_modification,omitempty" yaml:"latest_parsed_modification,omitempty"`
	Error                    string                 `json:"error,omitempty" yaml:"error,omitempty"`
	GoodModification         map[string]interface{} `json:"good_modification,omitempty" yaml:"good_modification,omitempty"`
}

// PipelineGroupsConfig holds information on the various pipeline groups present in GoCD.
type PipelineGroupsConfig struct {
	PipelineGroups PipelineGroups `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// PipelineGroups holds information on the various pipeline groups present in GoCD.
type PipelineGroups struct {
	PipelineGroups []PipelineGroup `json:"groups,omitempty" yaml:"groups,omitempty"`
}

// PipelineGroup holds information of a specific pipeline group instance.
type PipelineGroup struct {
	Name          string                           `json:"name,omitempty" yaml:"name,omitempty"`
	PipelineCount int                              `json:"pipeline_count,omitempty" yaml:"pipeline_count,omitempty"`
	Pipelines     []Pipeline                       `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
	Authorization PipelineGroupAuthorizationConfig `json:"authorization,omitempty" yaml:"authorization,omitempty"`
	ETAG          string                           `json:"etag,omitempty" yaml:"etag,omitempty"`
}

type PipelineGroupAuthorizationConfig struct {
	View    AuthorizationConfig `json:"view,omitempty" yaml:"view,omitempty"`
	Admins  AuthorizationConfig `json:"admins,omitempty" yaml:"admins,omitempty"`
	Operate AuthorizationConfig `json:"operate,omitempty" yaml:"operate,omitempty"`
}

// AuthorizationConfig of a pipeline group.
type AuthorizationConfig struct {
	Roles []string `json:"roles,omitempty" yaml:"roles,omitempty"`
	Users []string `json:"users,omitempty" yaml:"users,omitempty"`
}

// PipelineExport holds information of the pipeline that is exported to a specific config repo format.
type PipelineExport struct {
	PluginID         string `json:"plugin_id,omitempty" yaml:"plugin_id,omitempty"`
	PipelineFileName string `json:"pipeline_file_name,omitempty" yaml:"pipeline_file_name,omitempty"`
	PipelineContent  string `json:"pipeline_content,omitempty" yaml:"pipeline_content,omitempty"`
	ETAG             string `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// SystemAdmins holds information of the system admins present.
type SystemAdmins struct {
	Roles []string `json:"roles,omitempty" yaml:"roles,omitempty"`
	Users []string `json:"users,omitempty" yaml:"users,omitempty"`
	ETAG  string   `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// BackupConfig holds information of the backup configured.
type BackupConfig struct {
	EmailOnSuccess   bool   `json:"email_on_success,omitempty" yaml:"email_on_success,omitempty"`
	EmailOnFailure   bool   `json:"email_on_failure,omitempty" yaml:"email_on_failure,omitempty"`
	Schedule         string `json:"schedule,omitempty" yaml:"schedule,omitempty"`
	PostBackupScript string `json:"post_backup_script,omitempty" yaml:"post_backup_script,omitempty"`
}

// BackupStats holds information about the backup that was taken.
type BackupStats struct {
	Time           string `json:"time,omitempty" yaml:"time,omitempty"`
	Path           string `json:"path,omitempty" yaml:"path,omitempty"`
	Status         string `json:"status,omitempty" yaml:"status,omitempty"`
	ProgressStatus string `json:"progress_status,omitempty" yaml:"progress_status,omitempty"`
	Message        string `json:"message,omitempty" yaml:"message,omitempty"`
}

// PipelineConfig holds configuration information of a specific pipeline.
type PipelineConfig struct {
	Group                string                         `json:"group,omitempty" yaml:"group,omitempty"`
	LabelTemplate        string                         `json:"label_template,omitempty" yaml:"label_template,omitempty"`
	LockBehavior         string                         `json:"lock_behavior,omitempty" yaml:"lock_behavior,omitempty"`
	Name                 string                         `json:"name,omitempty" yaml:"name,omitempty"`
	Template             string                         `json:"template,omitempty" yaml:"template,omitempty"`
	Origin               PipelineOrigin                 `json:"origin" yaml:"origin"`
	Parameters           []PipelineEnvironmentVariables `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	EnvironmentVariables []PipelineEnvironmentVariables `json:"environment_variables,omitempty" yaml:"environment_variables,omitempty"`
	Materials            []Material                     `json:"materials,omitempty" yaml:"materials,omitempty"`
	Stages               []PipelineStageConfig          `json:"stages,omitempty" yaml:"stages,omitempty"`
	TrackingTool         PipelineTracingToolConfig      `json:"tracking_tool" yaml:"tracking_tool"`
	Timer                PipelineTimerConfig            `json:"timer" yaml:"timer"`
	CreateOptions        PipelineCreateOptions          `json:"create_options" yaml:"create_options"`
	Config               map[string]interface{}         `json:"config,omitempty" yaml:"config,omitempty"`
	ETAG                 string                         `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// PipelineOrigin holds the information on the source of pipeline present on GoCD.
type PipelineOrigin struct {
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	ID   string `json:"id,omitempty" yaml:"id,omitempty"`
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

// PipelineRunHistory holds information of pipeline run history.
type PipelineRunHistory struct {
	Name          string             `json:"name,omitempty" yaml:"name,omitempty"`
	Counter       int                `json:"counter,omitempty" yaml:"counter,omitempty"`
	ScheduledDate float64            `json:"scheduled_date,omitempty" yaml:"scheduled_date,omitempty"`
	BuildCause    PipelineBuildCause `json:"build_cause,omitempty" yaml:"build_cause,omitempty"`
}

type PipelineBuildCause struct {
	Message       string `json:"trigger_message,omitempty" yaml:"message,omitempty"`
	Approver      string `json:"approver,omitempty" yaml:"approver,omitempty"`
	TriggerForced bool   `json:"trigger_forced,omitempty" yaml:"trigger_forced,omitempty"`
}

// PipelineSchedules holds information of pipeline schedules.
type PipelineSchedules struct {
	Name   string                    `json:"pipelineName,omitempty" yaml:"pipelineName,omitempty"`
	Count  int                       `json:"count,omitempty" yaml:"count,omitempty"`
	Groups []PipelineSchedulesGroups `json:"groups,omitempty" yaml:"groups,omitempty"`
}

// PipelineSchedulesGroups holds information of pipeline schedules group.
type PipelineSchedulesGroups struct {
	History []PipelineSchedulesHistory `json:"history,omitempty" yaml:"history,omitempty"`
}

// PipelineSchedulesHistory holds information of pipeline schedule history.
type PipelineSchedulesHistory struct {
	Label              string `json:"label,omitempty" yaml:"label,omitempty"`
	ScheduledDate      string `json:"scheduled_date,omitempty" yaml:"scheduled_date,omitempty"`
	ScheduledTimestamp int64  `json:"scheduled_timestamp,omitempty" yaml:"scheduled_timestamp,omitempty"`
	ModificationDate   string `json:"modification_date,omitempty" yaml:"modification_date,omitempty"`
	BuildCause         string `json:"buildCauseBy,omitempty" yaml:"buildCauseBy,omitempty"`
}

// ScheduledJobs holds information of Scheduled Jobs.
type ScheduledJobs struct {
	Job []Job `xml:"job"`
}

// Job holds information of Scheduled Job.
type Job struct {
	Name         string `xml:"name,attr"`
	ID           string `xml:"id,attr"`
	BuildLocator string `xml:"buildLocator"`
	Environment  string `xml:"environment"`
}

// PipelineState holds information of the latest state of pipeline.
type PipelineState struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Paused      bool   `json:"paused,omitempty" yaml:"paused,omitempty"`
	Locked      bool   `json:"locked,omitempty" yaml:"locked,omitempty"`
	Schedulable bool   `json:"schedulable,omitempty" yaml:"schedulable,omitempty"`
	PausedBy    string `json:"paused_by,omitempty" yaml:"paused_by,omitempty"`
	PausedCause string `json:"paused_cause,omitempty" yaml:"paused_cause,omitempty"`
}

// Pipelines holds information of the pipelines present in GoCD.
type Pipelines struct {
	Pipelines []Pipeline `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
}

// Pipeline holds information of a specific pipeline instance.
type Pipeline struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

// EnvironmentInfo holds information of all environments present in GoCD.
type EnvironmentInfo struct {
	Environments Environments `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// Environments holds information of all environments present in GoCD.
type Environments struct {
	Environments []Environment `json:"environments,omitempty" yaml:"environments,omitempty"`
}

// Environment holds information of a specific environment present in GoCD.
type Environment struct {
	Name      string     `json:"name,omitempty" yaml:"name,omitempty"`
	Pipelines []Pipeline `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
	EnvVars   []EnvVars  `json:"environment_variables,omitempty" yaml:"environment_variables,omitempty"`
	ETAG      string     `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// EnvVars holds information of environment variables present in GoCD.
type EnvVars struct {
	Name           string `json:"name,omitempty" yaml:"name,omitempty"`
	Value          string `json:"value,omitempty" yaml:"value,omitempty"`
	EncryptedValue string `json:"encrypted_value,omitempty" yaml:"encrypted_value,omitempty"`
	Secure         bool   `json:"secure,omitempty" yaml:"secure,omitempty"`
}

// PatchEnvironment holds information that is handy while patching GoCD environment.
type PatchEnvironment struct {
	Name      string `json:"name" yaml:"name"`
	Pipelines struct {
		Add    []string `json:"add,omitempty" yaml:"add,omitempty"`
		Remove []string `json:"remove,omitempty" yaml:"remove,omitempty"`
	} `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
	EnvVars struct {
		Add []struct {
			Name  string `json:"name,omitempty" yaml:"name,omitempty"`
			Value string `json:"value,omitempty" yaml:"value,omitempty"`
		} `json:"add,omitempty" yaml:"add,omitempty"`
		Remove []string `json:"remove,omitempty" yaml:"remove,omitempty"`
	} `json:"environment_variables,omitempty" yaml:"environment_variables,omitempty"`
}

// VersionInfo holds version information of GoCD server.
type VersionInfo struct {
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	FullVersion string `json:"full_version,omitempty" yaml:"full_version,omitempty"`
	GitSHA      string `json:"git_sha,omitempty" yaml:"git_sha,omitempty"`
}

// AgentJobHistory holds information of pipeline run history of all GoCD agents.
type AgentJobHistory struct {
	Jobs       []JobRunHistory `json:"jobs,omitempty" yaml:"jobs,omitempty"`
	Pagination Pagination      `json:"pagination" yaml:"pagination"`
}

// JobRunHistory holds information of pipeline run history of a specific GoCD agent.
type JobRunHistory struct {
	Name            string `json:"pipeline_name,omitempty" yaml:"pipeline_name,omitempty"`
	JobName         string `json:"job_name,omitempty" yaml:"job_name,omitempty"`
	StageName       string `json:"stage_name,omitempty" yaml:"stage_name,omitempty"`
	StageCounter    int64  `json:"stage_counter,string,omitempty" yaml:"stage_counter,string,omitempty"`
	PipelineCounter int64  `json:"pipeline_counter,omitempty" yaml:"pipeline_counter,omitempty"`
	Result          string `json:"result,omitempty" yaml:"result,omitempty"`
}

// Pagination holds information which is helpful in paginating the results of job run history.
type Pagination struct {
	PageSize int64 `json:"page_size,omitempty" yaml:"page_size,omitempty"`
	Offset   int64 `json:"offset,omitempty" yaml:"offset,omitempty"`
	Total    int64 `json:"total,omitempty" yaml:"total,omitempty"`
}

// Maintenance holds latest information available in server about maintenance mode.
type Maintenance struct {
	MaintenanceInfo struct {
		Enabled  bool `json:"is_maintenance_mode,omitempty" yaml:"is_maintenance_mode,omitempty"`
		Metadata struct {
			UpdatedBy string `json:"updated_by,omitempty" yaml:"updated_by,omitempty"`
			UpdatedOn string `json:"updated_on,omitempty" yaml:"updated_on,omitempty"`
		} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	} `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// Encrypted holds the encrypted value of the passed plain text.
type Encrypted struct {
	EncryptedValue string `json:"encrypted_value,omitempty" yaml:"encrypted_value,omitempty"`
}

// ArtifactInfo holds the latest information of the artifacts.
type ArtifactInfo struct {
	ArtifactsDir  string        `json:"artifacts_dir,omitempty" yaml:"artifacts_dir,omitempty"`
	PurgeSettings PurgeSettings `json:"purge_settings,omitempty" yaml:"purge_settings,omitempty"`
	ETAG          string        `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// PurgeSettings holds the latest information on the purge setting configuration of GoCD.
type PurgeSettings struct {
	PurgeStartDiskSpace float64 `json:"purge_start_disk_space,omitempty" yaml:"purge_start_disk_space,omitempty"`
	PurgeUptoDiskSpace  float64 `json:"purge_upto_disk_space,omitempty" yaml:"purge_upto_disk_space,omitempty"`
}

// Schedule holds config of the pipeline that needs to be scheduled.
type Schedule struct {
	EnvVars        []map[string]interface{} `json:"environment_variables,omitempty" yaml:"environment_variables,omitempty"`
	Materials      []map[string]interface{} `json:"materials,omitempty" yaml:"materials,omitempty"`
	UpdateMaterial bool                     `json:"update_materials_before_scheduling,omitempty" yaml:"update_materials_before_scheduling,omitempty"`
}

// AuthConfigs holds information of multiple authorization configurations.
type AuthConfigs struct {
	Config struct {
		AuthConfigs []CommonConfig `json:"auth_configs" yaml:"auth_configs"`
	} `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// PluginConfiguration holds information of the various plugin properties.
type PluginConfiguration struct {
	Key            string                 `json:"key,omitempty" yaml:"key,omitempty"`
	Value          string                 `json:"value,omitempty" yaml:"value,omitempty"`
	EncryptedValue string                 `json:"encrypted_value,omitempty" yaml:"encrypted_value,omitempty"`
	IsSecure       bool                   `json:"is_secure,omitempty" yaml:"is_secure,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// SiteURLConfig holds information of the site url of GoCD.
type SiteURLConfig struct {
	SiteURL       string `json:"site_url,omitempty" yaml:"site_url,omitempty"`
	SecureSiteURL string `json:"secure_site_url,omitempty" yaml:"secure_site_url,omitempty"`
}

// MailServerConfig holds information required for GoCD mail-server configuration.
type MailServerConfig struct {
	Hostname          string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Port              int64  `json:"port,omitempty" yaml:"port,omitempty"`
	UserName          string `json:"username,omitempty" yaml:"username,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty" yaml:"encrypted_password,omitempty"`
	TLS               bool   `json:"tls,omitempty" yaml:"tls,omitempty"`
	SenderEmail       string `json:"sender_email,omitempty" yaml:"sender_email,omitempty"`
	AdminEmail        string `json:"admin_email,omitempty" yaml:"admin_email,omitempty"`
}

// PluginSettings holds information of plugin settings of GoCD.
type PluginSettings struct {
	ID            string                `json:"plugin_id,omitempty" yaml:"plugin_id,omitempty"`
	Configuration []PluginConfiguration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
	ETAG          string                `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// PipelineObject holds information required to comment/get/history of pipeline or instance of the same.
type PipelineObject struct {
	Name    string
	Counter int
	Message string
}

// PipelineHistory holds information of the pipeline history that also helps in paginating the responses.
type PipelineHistory struct {
	Links     map[string]interface{}   `json:"_links,omitempty" yaml:"_links,omitempty"`
	Pipelines []map[string]interface{} `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
}

// ArtifactStoresConfigs holds information of the specified artifact-stores/cluster-profiles/agent-profiles.
type ArtifactStoresConfigs struct {
	ArtifactStoresConfigs ArtifactStoresConfig `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// ArtifactStoresConfig holds information of all config-repos present in GoCD.
type ArtifactStoresConfig struct {
	CommonConfigs []CommonConfig `json:"artifact_stores,omitempty" yaml:"artifact_stores,omitempty"`
	ETAG          string         `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// ProfilesConfigs holds information of the specified artifact-stores/cluster-profiles/agent-profiles.
type ProfilesConfigs struct {
	ProfilesConfigs ProfilesConfig `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// ProfilesConfig holds information of all config-repos present in GoCD.
type ProfilesConfig struct {
	CommonConfigs         []CommonConfig `json:"profiles,omitempty" yaml:"profiles,omitempty"`
	ClusterProfilesConfig []CommonConfig `json:"cluster_profiles,omitempty" yaml:"cluster_profiles,omitempty"`
	ETAG                  string         `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// SecretsConfigs holds information of all secret configs present in GoCD.
type SecretsConfigs struct {
	SecretsConfigs SecretsConfig `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// SecretsConfig holds information of a specified secret config present in GoCD.
type SecretsConfig struct {
	CommonConfigs []CommonConfig `json:"secret_configs,omitempty" yaml:"secret_configs,omitempty"`
	ETAG          string         `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// CommonConfig holds information of the specified artifact store.
type CommonConfig struct {
	ID                  string                `json:"id,omitempty" yaml:"id,omitempty"`
	Name                string                `json:"name,omitempty" yaml:"name,omitempty"`
	PluginID            string                `json:"plugin_id,omitempty" yaml:"plugin_id,omitempty"`
	Description         string                `json:"description,omitempty" yaml:"description,omitempty"`
	ClusterProfileID    string                `json:"cluster_profile_id,omitempty" yaml:"cluster_profile_id,omitempty"`
	AllowOnlyKnownUsers bool                  `json:"allow_only_known_users_to_login,omitempty" yaml:"allow_only_known_users_to_login,omitempty"`
	Properties          []PluginConfiguration `json:"properties,omitempty" yaml:"properties,omitempty"`
	Rules               []map[string]string   `json:"rules,omitempty" yaml:"rules,omitempty"`
	ETAG                string                `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// PackageRepositories holds information of all package repositories present in GoCD.
type PackageRepositories struct {
	Repositories struct {
		PackageRepositories []PackageRepository `json:"package_repositories,omitempty" yaml:"package_repositories,omitempty"`
	} `json:"_embedded" yaml:"_embedded"`
}

// PackageRepository holds information of the specified package repository.
type PackageRepository struct {
	ID             string                `json:"repo_id,omitempty" yaml:"repo_id,omitempty"`
	Name           string                `json:"name,omitempty" yaml:"name,omitempty"`
	PluginMetaData map[string]string     `json:"plugin_metadata,omitempty" yaml:"plugin_metadata,omitempty"`
	Configuration  []PluginConfiguration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
	Packages       struct {
		Packages []CommonConfig `json:"packages,omitempty" yaml:"packages,omitempty"`
	} `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
	ETAG string `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// Packages holds information of all packages present in GoCD.
type Packages struct {
	Packages struct {
		Packages []Package `json:"packages,omitempty" yaml:"packages,omitempty"`
	} `json:"_embedded" yaml:"_embedded"`
}

// Package holds information of the specified packages of the package repository.
type Package struct {
	CommonConfig
	AutoUpdate    bool                  `json:"auto_update,omitempty" yaml:"auto_update,omitempty"`
	PackageRepos  CommonConfig          `json:"package_repo,omitempty" yaml:"package_repo,omitempty"`
	Configuration []PluginConfiguration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
	ETAG          string                `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// Materials holds information of all material type present in GoCD.
type Materials struct {
	Materials []Material `json:"materials,omitempty" yaml:"materials,omitempty"`
}

// MaterialUsage holds information of pipelines using a specific material.
type MaterialUsage struct {
	Usages []string `json:"usages,omitempty" yaml:"usages,omitempty"`
}

// Material holds information of a particular material type present in GoCD.
type Material struct {
	Type                     string              `json:"type,omitempty" yaml:"type,omitempty"`
	Fingerprint              string              `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	Attributes               Attribute           `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	RepoURL                  string              `json:"repository_url,omitempty" yaml:"repository_url,omitempty"`
	Config                   MaterialConfig      `json:"config,omitempty" yaml:"config,omitempty"`
	CanTriggerUpdate         bool                `json:"can_trigger_update,omitempty" yaml:"can_trigger_update,omitempty"`
	MaterialUpdateInProgress bool                `json:"material_update_in_progress,omitempty" yaml:"material_update_in_progress,omitempty"`
	Messages                 []map[string]string `json:"messages,omitempty" yaml:"messages,omitempty"`
}

// MaterialConfig holds information materials defined under Material.
type MaterialConfig struct {
	Type        string    `json:"type,omitempty" yaml:"type,omitempty"`
	Fingerprint string    `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	Attributes  Attribute `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// Attribute holds information of material type present in GoCD.
type Attribute struct {
	URL                 string            `json:"url,omitempty" yaml:"url,omitempty"`
	Username            string            `json:"username,omitempty" yaml:"username,omitempty"`
	Password            string            `json:"password,omitempty" yaml:"password,omitempty"`
	EncryptedPassword   string            `json:"encrypted_password,omitempty" yaml:"encrypted_password,omitempty"`
	Branch              string            `json:"branch,omitempty" yaml:"branch,omitempty"`
	AutoUpdate          bool              `json:"auto_update,omitempty" yaml:"auto_update,omitempty"`
	CheckExternals      bool              `json:"check_externals,omitempty" yaml:"check_externals,omitempty"`
	UseTickets          bool              `json:"use_tickets,omitempty" yaml:"use_tickets,omitempty"`
	View                string            `json:"view,omitempty" yaml:"view,omitempty"`
	Port                string            `json:"port,omitempty" yaml:"port,omitempty"`
	ProjectPath         string            `json:"project_path,omitempty" yaml:"project_path,omitempty"`
	Domain              string            `json:"domain,omitempty" yaml:"domain,omitempty"`
	Ref                 string            `json:"ref,omitempty" yaml:"ref,omitempty"`
	Name                string            `json:"name,omitempty" yaml:"name,omitempty"`
	Stage               string            `json:"stage,omitempty" yaml:"stage,omitempty"`
	Pipeline            string            `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
	IgnoreForScheduling bool              `json:"ignore_for_scheduling,omitempty" yaml:"ignore_for_scheduling,omitempty"`
	Destination         string            `json:"destination,omitempty" yaml:"destination,omitempty"`
	InvertFilter        bool              `json:"invert_filter,omitempty" yaml:"invert_filter,omitempty"`
	SubmoduleFolder     string            `json:"submodule_folder,omitempty" yaml:"submodule_folder,omitempty"`
	ShallowClone        bool              `json:"shallow_clone,omitempty" yaml:"shallow_clone,omitempty"`
	Origin              map[string]string `json:"origin,omitempty" yaml:"origin,omitempty"`
	Filter              struct {
		Ignore []string `json:"ignore,omitempty" yaml:"ignore,omitempty"`
	} `json:"filter,omitempty" yaml:"filter,omitempty"`
}

// RolesConfigs holds information of all role configs present in GoCd.
type RolesConfigs struct {
	RolesConfigs RolesConfig `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// RolesConfig holds information of all role configs present in GoCd.
type RolesConfig struct {
	Role []Role `json:"roles,omitempty" yaml:"roles,omitempty"`
	ETAG string `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// Role holds information of a specific role in GoCd.
type Role struct {
	Name         string                `json:"name,omitempty" yaml:"name,omitempty"`
	Type         string                `json:"type,omitempty" yaml:"type,omitempty"`
	Attributes   RoleAttribute         `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Policy       []map[string]string   `json:"policy,omitempty" yaml:"policy,omitempty"`
	AuthConfigID string                `json:"auth_config_id,omitempty" yaml:"auth_config_id,omitempty"`
	Properties   []PluginConfiguration `json:"properties,omitempty" yaml:"properties,omitempty"`
	ETAG         string                `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// RoleAttribute holds information of a specific attribute of a role in GoCd.
type RoleAttribute struct {
	Users        []string              `json:"users,omitempty" yaml:"users,omitempty"`
	AuthConfigID string                `json:"auth_config_id,omitempty" yaml:"auth_config_id,omitempty"`
	Properties   []PluginConfiguration `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// PluginsInfos holds information of all plugins present in GoCd.
type PluginsInfos struct {
	PluginsInfos PluginsInfo `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// PluginsInfo holds information of all plugins present in GoCd.
type PluginsInfo struct {
	Plugins []Plugin `json:"plugin_info,omitempty" yaml:"plugin_info,omitempty"`
	ETAG    string   `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// Plugin holds information of a specific plugins present in GoCd.
type Plugin struct {
	ID                 string                 `json:"id,omitempty" yaml:"id,omitempty"`
	Status             PluginStatus           `json:"status,omitempty" yaml:"status,omitempty"`
	PluginFileLocation string                 `json:"plugin_file_location,omitempty" yaml:"plugin_file_location,omitempty"`
	BundledPlugin      bool                   `json:"bundled_plugin,omitempty" yaml:"bundled_plugin,omitempty"`
	About              map[string]interface{} `json:"about,omitempty" yaml:"about,omitempty"`
	Extensions         []PluginAttributes     `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	ETAG               string                 `json:"etag,omitempty" yaml:"etag,omitempty"`
}

// PluginStatus holds the status information of the plugin.
type PluginStatus struct {
	State string `json:"state,omitempty" yaml:"state,omitempty"`
}

// PluginAttributes holds information all available attributes of a given plugin.
type PluginAttributes struct {
	Type                        string                  `json:"type,omitempty" yaml:"type,omitempty"`
	DisplayName                 string                  `json:"display_name,omitempty" yaml:"display_name,omitempty"`
	AuthConfigSettings          *PluginSettingAttribute `json:"auth_config_settings,omitempty" yaml:"auth_config_settings,omitempty"`
	ArtifactConfigSettings      *PluginSettingAttribute `json:"artifact_config_settings,omitempty" yaml:"artifact_config_settings,omitempty"`
	ElasticAgentProfileSettings *PluginSettingAttribute `json:"elastic_agent_profile_settings,omitempty" yaml:"elastic_agent_profile_settings,omitempty"`
	FetchArtifactSettings       *PluginSettingAttribute `json:"fetch_artifact_settings,omitempty" yaml:"fetch_artifact_settings,omitempty"`
	ClusterProfileSettings      *PluginSettingAttribute `json:"cluster_profile_settings,omitempty" yaml:"cluster_profile_settings,omitempty"`
	PluginSettings              *PluginSettingAttribute `json:"plugin_settings,omitempty" yaml:"plugin_settings,omitempty"`
	PackageSettings             *PluginSettingAttribute `json:"package_settings,omitempty" yaml:"package_settings,omitempty"`
	RepositorySettings          *PluginSettingAttribute `json:"repository_settings,omitempty" yaml:"repository_settings,omitempty"`
	ScmSettings                 *PluginSettingAttribute `json:"scm_settings,omitempty" yaml:"scm_settings,omitempty"`
	StoreConfigSettings         *PluginSettingAttribute `json:"store_config_settings,omitempty" yaml:"store_config_settings,omitempty"`
	SecretConfigSettings        *PluginSettingAttribute `json:"secret_config_settings,omitempty" yaml:"secret_config_settings,omitempty"`
	RoleSettings                *PluginSettingAttribute `json:"role_settings,omitempty" yaml:"role_settings,omitempty"`
	TaskSettings                *PluginSettingAttribute `json:"task_settings,omitempty" yaml:"task_settings,omitempty"`
}

// PluginSettingAttribute holds plugin configuration information.
type PluginSettingAttribute struct {
	Configurations []PluginConfiguration `json:"configurations,omitempty" yaml:"configurations,omitempty"`
}

// Users holds information of all users present in GoCD.
type Users struct {
	GoCDUsers struct {
		Users []User `json:"users,omitempty" yaml:"users,omitempty"`
	} `json:"_embedded,omitempty" yaml:"_embedded,omitempty"`
}

// User holds information of the users present in GoCD.
// This is golang implementation of GoCD's user API https://api.gocd.org/current/#the-user-object.
type User struct {
	Name         string     `json:"display_name,omitempty" yaml:"display_name,omitempty"`
	LoginName    string     `json:"login_name,omitempty" yaml:"login_name,omitempty"`
	Enabled      bool       `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	EmailID      string     `json:"email,omitempty" yaml:"email,omitempty"`
	EmailMe      bool       `json:"email_me,omitempty" yaml:"email_me,omitempty"`
	Admin        bool       `json:"admin,omitempty" yaml:"admin,omitempty"`
	CheckInAlias []string   `json:"checkin_aliases,omitempty" yaml:"checkin_aliases,omitempty"`
	Roles        []UserRole `json:"roles,omitempty" yaml:"roles,omitempty"`
}

// UserRole holds information of the user role present in GoCD.
// This is golang implementation of GoCD's role API https://api.gocd.org/current/#the-user-role-object
type UserRole struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

// ElasticProfileUsage holds information on elastic agent profile being used by a pipeline/stage/job.
// This is golang implementation of GoCD's internal API 'api/internal/elastic/profiles/<elastic-profile-name>/usages'.
type ElasticProfileUsage struct {
	PipelineName         string `json:"pipeline_name,omitempty" yaml:"pipeline_name,omitempty"`
	StageName            string `json:"stage_name,omitempty" yaml:"stage_name,omitempty"`
	JobName              string `json:"job_name,omitempty" yaml:"job_name,omitempty"`
	PipelineConfigOrigin string `json:"pipeline_config_origin,omitempty" yaml:"pipeline_config_origin,omitempty"`
}

// Stage holds information of GoCd stage.
type Stage struct {
	Pipeline         string   `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
	PipelineInstance string   `json:"pipeline_instance,omitempty" yaml:"pipeline_instance,omitempty"`
	StageCounter     string   `json:"stage_counter,omitempty" yaml:"stage_counter,omitempty"`
	Name             string   `json:"name,omitempty" yaml:"name,omitempty"`
	Jobs             []string `json:"jobs,omitempty" yaml:"jobs,omitempty"`
}

type VSM struct {
	Pipeline string           `json:"current_pipeline,omitempty" yaml:"current_pipeline,omitempty"`
	Level    []PipelineLevels `json:"levels,omitempty" yaml:"levels,omitempty"`
}

type PipelineLevels struct {
	Nodes []PipelineNode `json:"nodes,omitempty" yaml:"nodes,omitempty"`
}

type PipelineNode struct {
	Parents    []string `json:"parents,omitempty" yaml:"parents,omitempty"`
	Dependents []string `json:"dependents,omitempty" yaml:"dependents,omitempty"`
	Name       string   `json:"name,omitempty" yaml:"name,omitempty"`
	ID         string   `json:"ID,omitempty" yaml:"ID,omitempty"`
}

// Permission holds information of permissions that the invoked user has in GoCD.
type Permission struct {
	Environment         EntityPermissions `json:"environment,omitempty" yaml:"environment,omitempty"`
	ConfigRepo          EntityPermissions `json:"config_repo,omitempty" yaml:"config_repo,omitempty"`
	ClusterProfile      EntityPermissions `json:"cluster_profile,omitempty" yaml:"cluster_profile,omitempty"`
	ElasticAgentProfile EntityPermissions `json:"elastic_agent_profile,omitempty" yaml:"elastic_agent_profile,omitempty"`
}

// EntityPermissions is the list of the permission level one has over various entities in GoCD.
type EntityPermissions struct {
	View       []string `json:"view,omitempty" yaml:"view,omitempty"`
	Administer []string `json:"administer,omitempty" yaml:"administer,omitempty"`
}

// PipelineTemplateConfig holds pipeline template config.
type PipelineTemplateConfig struct {
	Name   string                `json:"name,omitempty" yaml:"name,omitempty"`
	Stages []PipelineStageConfig `json:"stages,omitempty" yaml:"stages,omitempty"`
}

// PipelineStageConfig holds information of pipeline stage.
type PipelineStageConfig struct {
	Name                  string                         `json:"name,omitempty" yaml:"name,omitempty"`
	FetchMaterials        bool                           `json:"fetch_materials,omitempty" yaml:"fetch_materials,omitempty"`
	CleanWorkingDirectory bool                           `json:"clean_working_directory,omitempty" yaml:"clean_working_directory,omitempty"`
	NeverCleanupArtifacts bool                           `json:"never_cleanup_artifacts,omitempty" yaml:"never_cleanup_artifacts,omitempty"`
	Approval              PipelineApprovalConfig         `json:"approval,omitempty" yaml:"approval,omitempty"`
	EnvironmentVariables  []PipelineEnvironmentVariables `json:"environment_variables,omitempty" yaml:"environment_variables,omitempty"`
	Job                   []PipelineJobConfig            `json:"job,omitempty" yaml:"job,omitempty"`
}

// PipelineJobConfig holds information of pipeline job.
type PipelineJobConfig struct {
	Name                 string                         `json:"name,omitempty" yaml:"name,omitempty"`
	RunInstanceCount     interface{}                    `json:"run_instance_count,omitempty" yaml:"run_instance_count,omitempty"`
	ElasticProfileID     string                         `json:"elastic_profile_id,omitempty" yaml:"elastic_profile_id,omitempty"`
	Timeout              int                            `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	EnvironmentVariables []PipelineEnvironmentVariables `json:"environment_variables,omitempty" yaml:"environment_variables,omitempty"`
	Resources            []string                       `json:"resources,omitempty" yaml:"resources,omitempty"`
	Tasks                []PipelineTaskConfig           `json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Tabs                 []PipelineTab                  `json:"tabs,omitempty" yaml:"tabs,omitempty"`
	Artifacts            []PipelineArtifact             `json:"artifacts,omitempty" yaml:"artifacts,omitempty"`
}

// PipelineEnvironmentVariables holds information of all Environment variable defined in pipeline.
type PipelineEnvironmentVariables struct {
	Name           string `json:"name,omitempty" yaml:"name,omitempty"`
	Value          string `json:"value,omitempty" yaml:"value,omitempty"`
	Secure         bool   `json:"secure,omitempty" yaml:"secure,omitempty"`
	EncryptedValue string `json:"encrypted_value,omitempty" yaml:"encrypted_value,omitempty"`
}

// PipelineApprovalConfig holds information on pipeline approval config.
type PipelineApprovalConfig struct {
	Type               string              `json:"type,omitempty" yaml:"type,omitempty"`
	AllowOnlyOnSuccess bool                `json:"allow_only_on_success,omitempty" yaml:"allow_only_on_success,omitempty"`
	Authorization      AuthorizationConfig `json:"authorization,omitempty" yaml:"authorization,omitempty"`
}

type PipelineTab struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

// PipelineArtifact holds information on artifacts defined in pipeline.
type PipelineArtifact struct {
	Type          string                `json:"type,omitempty" yaml:"type,omitempty"`
	ArtifactID    string                `json:"artifact_id,omitempty" yaml:"artifact_id,omitempty"`
	StoreID       string                `json:"store_id,omitempty" yaml:"store_id,omitempty"`
	Configuration []PluginConfiguration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
}

// PipelineTaskConfig holds information of the tasks defined in jobs of GoCD pipeline.
type PipelineTaskConfig struct {
	Type       string              `json:"type,omitempty" yaml:"type,omitempty"`
	Attributes TaskAttributeConfig `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// TaskAttributeConfig holds information of tasks.
type TaskAttributeConfig struct {
	IsSourceAFile       bool                  `json:"is_source_a_file,omitempty" yaml:"is_source_a_file,omitempty"`
	ArtifactID          string                `json:"artifact_id,omitempty" yaml:"artifact_id,omitempty"`
	ArtifactOrigin      string                `json:"artifact_origin,omitempty" yaml:"artifact_origin,omitempty"`
	Command             string                `json:"command,omitempty" yaml:"command,omitempty"`
	WorkingDirectory    string                `json:"working_directory,omitempty" yaml:"working_directory,omitempty"`
	Pipeline            string                `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
	Stage               string                `json:"stage,omitempty" yaml:"stage,omitempty"`
	Job                 string                `json:"job,omitempty" yaml:"job,omitempty"`
	Source              string                `json:"source,omitempty" yaml:"source,omitempty"`
	Destination         string                `json:"destination,omitempty" yaml:"destination,omitempty"`
	Configuration       []PluginConfiguration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
	RunIf               []string              `json:"run_if,omitempty" yaml:"run_if,omitempty"`
	Arguments           []string              `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	PluginConfiguration struct {
		ID      string `json:"id,omitempty" yaml:"id,omitempty"`
		Version string `json:"version,omitempty" yaml:"version,omitempty"`
	} `json:"plugin_configuration,omitempty" yaml:"plugin_configuration,omitempty"`
	OnCancel struct {
		RunIf            []string `json:"run_if,omitempty" yaml:"run_if,omitempty"`
		Command          string   `json:"command,omitempty" yaml:"command,omitempty"`
		Arguments        []string `json:"arguments,omitempty" yaml:"arguments,omitempty"`
		WorkingDirectory string   `json:"working_directory,omitempty" yaml:"working_directory,omitempty"`
	} `json:"on_cancel,omitempty" yaml:"on_cancel,omitempty"`
}

// PipelineCreateOptions holds information of Pipeline Create Options.
type PipelineCreateOptions struct {
	PausePipeline bool   `json:"pause_pipeline,omitempty" yaml:"pause_pipeline,omitempty"`
	PauseReason   string `json:"pause_reason,omitempty" yaml:"pause_reason,omitempty"`
}

// PipelineTimerConfig holds information of GoCD pipeline timers, when to trigger pipeline.
type PipelineTimerConfig struct {
	Spec          string `json:"spec,omitempty" yaml:"spec,omitempty"`
	OnlyOnChanges bool   `json:"only_on_changes,omitempty" yaml:"only_on_changes,omitempty"`
}

// PipelineTracingToolConfig holds information of tracing tool for issue tracking.
type PipelineTracingToolConfig struct {
	Type       string `json:"type,omitempty" yaml:"type,omitempty"`
	Attributes struct {
		URLPattern string `json:"url_pattern,omitempty" yaml:"url_pattern,omitempty"`
		Regex      string `json:"regex,omitempty" yaml:"regex,omitempty"`
	} `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

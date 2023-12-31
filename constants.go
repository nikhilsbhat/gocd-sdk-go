package gocd

const (
	AgentsEndpoint              = "/api/agents"
	VersionEndpoint             = "/api/version"
	ServerHealthEndpoint        = "/api/server_health_messages"
	ConfigReposEndpoint         = "/api/admin/config_repos"
	ConfigReposInternalEndpoint = "/api/internal/config_repos"
	SystemAdminEndpoint         = "/api/admin/security/system_admins"
	BackupConfigEndpoint        = "/api/config/backup"
	BackupStatsEndpoint         = "/api/backups"
	PipelineGroupEndpoint       = "/api/admin/pipeline_groups"
	EnvironmentEndpoint         = "/api/admin/environments"
	JobRunHistoryEndpoint       = "/api/agents/%s/job_run_history"
	LastXPipelineScheduledDates = "/pipelineHistory.json?pipelineName=%s"
	MaintenanceEndpoint         = "/api/admin/maintenance_mode"
	APIFeedPipelineEndpoint     = "/api/feed/pipelines.xml"
	APIJobFeedEndpoint          = "/api/feed/jobs/scheduled.xml"
	JobsAPIEndpoint             = "/api/jobs"
	StageEndpoint               = "/api/stages"
	PipelineStatus              = "/api/pipelines/%s/status"
	EncryptEndpoint             = "/api/admin/encrypt"
	ArtifactInfoEndpoint        = "/api/admin/config/server/artifact_config"
	PipelinesEndpoint           = "/api/pipelines"
	PipelineConfigEndpoint      = "/api/admin/pipelines"
	PipelineExportEndpoint      = "/api/admin/export/pipelines"
	HealthEndpoint              = "/api/v1/health"
	DefaultTimeoutEndpoint      = "/api/admin/config/server/default_job_timeout"
	MailServerConfigEndpoint    = "/api/config/mailserver"
	PluginSettingsEndpoint      = "/api/admin/plugin_settings"
	AuthConfigEndpoint          = "/api/admin/security/auth_configs"
	ClusterProfileEndpoint      = "/api/admin/elastic/cluster_profiles"
	AgentProfileEndpoint        = "/api/elastic/profiles"
	ArtifactStoreEndpoint       = "/api/admin/artifact_stores"
	SiteURLEndpoint             = "/api/admin/config/server/site_urls"
	SecretsConfigEndpoint       = "/api/admin/secret_configs" //nolint:gosec
	PackageRepositoriesEndpoint = "/api/admin/repositories"
	PackagesEndpoint            = "/api/admin/packages"
	MaterialEndpoint            = "/api/internal/materials"
	MaterialUsageEndpoint       = "/api/internal/materials/%s/usages"
	MaterialNotifyEndpoint      = "/api/admin/materials/%s/notify"
	MaterialTriggerUpdate       = "/api/internal/materials/%s/trigger_update"
	RolesEndpoint               = "/api/admin/security/roles"
	PluginInfoEndpoint          = "/api/admin/plugin_info"
	UsersEndpoint               = "/api/users"
	AdminOperationStateEndpoint = "/api/admin/operations/state"
	ElasticProfileUsageEndpoint = "/api/internal/elastic/profiles/%s/usages"
	PreflightCheckEndpoint      = "/api/admin/config_repo_ops/preflight"
	CurrentUserEndpoint         = "/api/current_user"
	PermissionsEndpoint         = "/api/auth/permissions"
	VSMEndpoint                 = "/pipelines/value_stream_map"
	HeaderVersionZero           = "application/vnd.go.cd+json"
	HeaderVersionOne            = "application/vnd.go.cd.v1+json"
	HeaderVersionTwo            = "application/vnd.go.cd.v2+json"
	HeaderVersionThree          = "application/vnd.go.cd.v3+json"
	HeaderVersionFour           = "application/vnd.go.cd.v4+json"
	HeaderVersionSeven          = "application/vnd.go.cd.v7+json"
	HeaderVersionEleven         = "application/vnd.go.cd.v11+json"
)

const (
	goCdAPILoggerName = "gocd-sdk-go"
	ContentJSON       = "application/json"
	HeaderConfirm     = "X-GoCD-Confirm"
	PipelinePrefix    = "/go/api/feed/pipelines/"
	PipelineSuffix    = "/stages.xml"
	LocationHeader    = "Location"
)

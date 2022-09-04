package gocd

const (
	GoCdAgentsEndpoint          = "/api/agents"
	GoCdVersionEndpoint         = "/api/version"
	GoCdServerHealthEndpoint    = "/api/server_health_messages"
	GoCdConfigReposEndpoint     = "/api/admin/config_repos"
	GoCdSystemAdminEndpoint     = "/api/admin/security/system_admins"
	GoCdBackupConfigEndpoint    = "/api/config/backup"
	GoCdPipelineGroupEndpoint   = "/api/admin/pipeline_groups"
	GoCdEnvironmentEndpoint     = "/api/admin/environments"
	GoCdJobRunHistoryEndpoint   = "/api/agents/%s/job_run_history"
	GoCdAPIFeedPipelineEndpoint = "/api/feed/pipelines.xml"
	GoCdPipelineStatus          = "/api/pipelines/%s/status"
	GoCdHeaderVersionOne        = "application/vnd.go.cd.v1+json"
	GoCdHeaderVersionTwo        = "application/vnd.go.cd.v2+json"
	GoCdHeaderVersionThree      = "application/vnd.go.cd.v3+json"
	GoCdHeaderVersionFour       = "application/vnd.go.cd.v4+json"
	GoCdHeaderVersionSeven      = "application/vnd.go.cd.v7+json"
)

const (
	goCdAPILoggerName = "gocd-sdk-go"
)

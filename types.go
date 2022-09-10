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
	PluginID string `json:"plugin_id"`
	ID       string `json:"config_repos,omitempty"`
	Material struct {
		Type       string `json:"type,omitempty"`
		Attributes struct {
			URL               string `json:"url,omitempty"`
			Username          string `json:"username,omitempty"`
			Password          string `json:"password,omitempty"`
			EncryptedPassword string `json:"encrypted_password,omitempty"`
			Branch            string `json:"branch,omitempty"`
			AutoUpdate        bool   `json:"auto_update,omitempty"`
		}
	}
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
	Name          string     `json:"name,omitempty"`
	PipelineCount int        `json:"pipeline_count,omitempty"`
	Pipelines     []Pipeline `json:"pipelines,omitempty"`
}

// SystemAdmins holds information of the system admins present.
type SystemAdmins struct {
	Roles []string `json:"roles,omitempty"`
	Users []string `json:"users,omitempty"`
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

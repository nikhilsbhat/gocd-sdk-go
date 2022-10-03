package gocd

import (
	"crypto/tls"
	"crypto/x509"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// client holds resty.Client which could be used for interacting with GoCD and other information.
type client struct {
	httpClient *resty.Client
	logger     *log.Logger
}

// GoCd implements methods to get various information from GoCD.
type GoCd interface {
	GetAgents() ([]Agent, error)
	GetAgentJobRunHistory(agent string) (AgentJobHistory, error)
	UpdateAgent(id string, agent Agent) error
	UpdateAgentBulk(agent Agent) error
	DeleteAgent(id string) (string, error)
	DeleteAgentBulk(agent Agent) (string, error)
	AgentKillTask(agent Agent) error
	GetHealthMessages() ([]ServerHealth, error)
	GetConfigRepos() ([]ConfigRepo, error)
	GetConfigRepo(repo string) (ConfigRepo, error)
	CreateConfigRepo(repoObj ConfigRepo) error
	UpdateConfigRepo(repo ConfigRepo, etag string) (string, error)
	DeleteConfigRepo(repo string) error
	ConfigRepoStatus(repo string) (map[string]bool, error)
	ConfigRepoTriggerUpdate(name string) (map[string]string, error)
	EnableMaintenanceMode() error
	DisableMaintenanceMode() error
	GetMaintenanceModeInfo() (Maintenance, error)
	GetSystemAdmins() (SystemAdmins, error)
	UpdateSystemAdmins(data SystemAdmins) (SystemAdmins, error)
	CreatePipelineGroup(group PipelineGroup) error
	GetPipelineGroups() ([]PipelineGroup, error)
	GetPipelineGroup(name string) (PipelineGroup, error)
	DeletePipelineGroup(name string) error
	UpdatePipelineGroup(group PipelineGroup) (PipelineGroup, error)
	GetEnvironments() ([]Environment, error)
	GetEnvironment(name string) (Environment, error)
	CreateEnvironment(environment Environment) error
	UpdateEnvironment(environment Environment) (Environment, error)
	PatchEnvironment(environment any) (Environment, error)
	DeleteEnvironment(name string) error
	GetVersionInfo() (VersionInfo, error)
	CreateOrUpdateBackup(backup BackupConfig) error
	GetBackupInfo() (BackupConfig, error)
	GetPipelines() (PipelinesInfo, error)
	GetPipelineState(pipeline string) (PipelineState, error)
	PipelinePause(name string, message any) error
	PipelineUnPause(name string) error
	PipelineUnlock(name string) error
	SchedulePipeline(name string, schedule Schedule) error
	SetRetryCount(count int)
	SetRetryWaitTime(count int)
	EncryptText(value string) (Encrypted, error)
	GetArtifactConfig() (ArtifactInfo, error)
	UpdateArtifactConfig(ArtifactInfo) (ArtifactInfo, error)
	GetAuthConfigs() ([]AuthConfig, error)
	GetAuthConfig(name string) (AuthConfig, error)
	CreateAuthConfig(config AuthConfig) (AuthConfig, error)
	UpdateAuthConfig(config AuthConfig) (AuthConfig, error)
	DeleteAuthConfig(name string) error
	GetSiteURL() (SiteURLConfig, error)
	CreateOrUpdateSiteURL(SiteURLConfig) (SiteURLConfig, error)
}

// NewClient returns new instance of httpClient when invoked.
func NewClient(baseURL, userName, passWord, logLevel string,
	caContent []byte,
) GoCd {
	logger := log.New()
	logger.SetLevel(GetLoglevel(logLevel))
	logger.WithField(goCdAPILoggerName, true)
	logger.SetFormatter(&log.JSONFormatter{})

	newClient := resty.New()
	newClient.SetLogger(logger)
	newClient.SetRetryCount(defaultRetryCount)
	newClient.SetRetryWaitTime(defaultRetryWaitTime * time.Second)
	if logLevel == "debug" {
		newClient.SetDebug(true)
	}
	newClient.SetBaseURL(baseURL)
	newClient.SetBasicAuth(userName, passWord)
	if len(caContent) != 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caContent)
		newClient.SetTLSClientConfig(&tls.Config{RootCAs: certPool}) //nolint:gosec
	} else {
		newClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}

	return &client{
		httpClient: newClient,
		logger:     logger,
	}
}

// GetLoglevel sets the loglevel to the kind of log asked for.
func GetLoglevel(level string) log.Level {
	switch strings.ToLower(level) {
	case log.WarnLevel.String():
		return log.WarnLevel
	case log.DebugLevel.String():
		return log.DebugLevel
	case log.TraceLevel.String():
		return log.TraceLevel
	case log.FatalLevel.String():
		return log.FatalLevel
	case log.ErrorLevel.String():
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}

// SetRetryCount sets retry count for the go-resty client.
func (conf *client) SetRetryCount(count int) {
	conf.httpClient.SetRetryCount(count)
}

// SetRetryWaitTime sets retry wait time for the go-resty client.
func (conf *client) SetRetryWaitTime(count int) {
	conf.httpClient.SetRetryWaitTime(time.Duration(count) * time.Second)
}

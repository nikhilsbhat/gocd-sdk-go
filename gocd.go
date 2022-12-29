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
	GetBackupConfig() (BackupConfig, error)
	CreateOrUpdateBackupConfig(backup BackupConfig) error
	DeleteBackupConfig() error
	GetPipelines() (PipelinesInfo, error)
	GetPipelineState(pipeline string) (PipelineState, error)
	PipelinePause(name string, message any) error
	PipelineUnPause(name string) error
	PipelineUnlock(name string) error
	SchedulePipeline(name string, schedule Schedule) error
	GetPipelineInstance(pipeline PipelineObject) (map[string]interface{}, error)
	GetPipelineHistory(name string, size, after int) ([]map[string]interface{}, error)
	CommentOnPipeline(comment PipelineObject) error
	EncryptText(value string) (Encrypted, error)
	GetArtifactConfig() (ArtifactInfo, error)
	UpdateArtifactConfig(ArtifactInfo) (ArtifactInfo, error)
	GetAuthConfigs() ([]CommonConfig, error)
	GetAuthConfig(name string) (CommonConfig, error)
	CreateAuthConfig(config CommonConfig) (CommonConfig, error)
	UpdateAuthConfig(config CommonConfig) (CommonConfig, error)
	DeleteAuthConfig(name string) error
	GetSiteURL() (SiteURLConfig, error)
	CreateOrUpdateSiteURL(SiteURLConfig) (SiteURLConfig, error)
	GetMailServerConfig() (MailServerConfig, error)
	CreateOrUpdateMailServerConfig(mailConfig MailServerConfig) (MailServerConfig, error)
	DeleteMailServerConfig() error
	GetDefaultJobTimeout() (map[string]string, error)
	UpdateDefaultJobTimeout(timeoutMinutes int) error
	GetPluginSettings(name string) (PluginSettings, error)
	CreatePluginSettings(settings PluginSettings) (PluginSettings, error)
	UpdatePluginSettings(settings PluginSettings) (PluginSettings, error)
	GetClusterProfiles() (ProfilesConfig, error)
	GetClusterProfile(name string) (CommonConfig, error)
	CreateClusterProfile(config CommonConfig) (CommonConfig, error)
	UpdateClusterProfile(config CommonConfig) (CommonConfig, error)
	DeleteClusterProfile(name string) error
	GetElasticAgentProfiles() (ProfilesConfig, error)
	GetElasticAgentProfile(name string) (CommonConfig, error)
	CreateElasticAgentProfile(config CommonConfig) (CommonConfig, error)
	UpdateElasticAgentProfile(config CommonConfig) (CommonConfig, error)
	DeleteElasticAgentProfile(name string) error
	GetArtifactStores() (ArtifactStoresConfig, error)
	GetArtifactStore(name string) (CommonConfig, error)
	CreateArtifactStore(config CommonConfig) (CommonConfig, error)
	UpdateArtifactStore(config CommonConfig) (CommonConfig, error)
	DeleteArtifactStore(name string) error
	GetSecretConfigs() (SecretsConfig, error)
	GetSecretConfig(id string) (CommonConfig, error)
	CreateSecretConfig(config CommonConfig) (CommonConfig, error)
	UpdateSecretConfig(config CommonConfig) (CommonConfig, error)
	DeleteSecretConfig(id string) error
	GetPackageRepositories() ([]PackageRepository, error)
	GetPackageRepository(id string) (PackageRepository, error)
	CreatePackageRepository(config PackageRepository) (PackageRepository, error)
	UpdatePackageRepository(config PackageRepository) (PackageRepository, error)
	DeletePackageRepository(id string) error
	GetPackages() ([]Package, error)
	GetPackage(id string) (Package, error)
	CreatePackage(config Package) (Package, error)
	UpdatePackage(config Package) (Package, error)
	DeletePackage(id string) error
	GetMaterials() ([]Material, error)
	SetRetryCount(count int)
	SetRetryWaitTime(count int)
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

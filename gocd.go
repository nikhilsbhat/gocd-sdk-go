package gocd

import (
	"crypto/tls"
	"crypto/x509"
	"reflect"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
	goCdLogger "github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/plugin"
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
	GetAgent(agentID string) (Agent, error)
	GetAgentJobRunHistory(agent string) (AgentJobHistory, error)
	UpdateAgent(agent Agent) error
	UpdateAgentBulk(agent Agent) error
	DeleteAgent(id string) (string, error)
	DeleteAgentBulk(agent Agent) (string, error)
	AgentKillTask(agent Agent) error
	GetServerHealthMessages() ([]ServerHealth, error)
	GetServerHealth() (map[string]string, error)
	GetConfigRepos() ([]ConfigRepo, error)
	GetConfigReposInternal() ([]ConfigRepo, error)
	GetConfigRepo(repo string) (ConfigRepo, error)
	GetConfigRepoDefinitions(repo string) (ConfigRepo, error)
	CreateConfigRepo(repoObj ConfigRepo) error
	UpdateConfigRepo(repo ConfigRepo) (string, error)
	DeleteConfigRepo(repo string) error
	ConfigRepoStatus(repo string) (map[string]bool, error)
	ConfigRepoTriggerUpdate(name string) (map[string]string, error)
	ConfigRepoPreflightCheck(pipelines map[string]string, pluginID string, repoID string) (bool, error)
	SetPipelineFiles(pipelines []PipelineFiles) map[string]string
	GetPipelineFiles(path string, pipelines []string, patterns ...string) ([]PipelineFiles, error)
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
	GetPipelineRunHistory(pipeline, pageSize string, delay time.Duration) ([]PipelineRunHistory, error)
	GetLimitedPipelineRunHistory(pipeline, pageSize, after string) ([]PipelineRunHistory, error)
	GetPipelineSchedules(pipeline, start, perPage string) (PipelineSchedules, error)
	ValidatePipelineSyntax(pluginCfg plugin.Plugin, pipelines []string, fetchVersionFromServer bool) (bool, error)
	GetEnvironments() ([]Environment, error)
	GetEnvironment(name string) (Environment, error)
	GetEnvironmentsMerged(names []string) ([]Environment, error)
	CreateEnvironment(environment Environment) error
	UpdateEnvironment(environment Environment) (Environment, error)
	PatchEnvironment(environment any) (Environment, error)
	DeleteEnvironment(name string) error
	GetVersionInfo() (VersionInfo, error)
	GetBackupConfig() (BackupConfig, error)
	CreateOrUpdateBackupConfig(backup BackupConfig) error
	DeleteBackupConfig() error
	GetBackup(ID string) (BackupStats, error)
	ScheduleBackup() (map[string]string, error)
	GetPipelines() (PipelinesInfo, error)
	GetPipelineState(pipeline string) (PipelineState, error)
	PipelinePause(name string, message any) error
	PipelineUnPause(name string) error
	PipelineUnlock(name string) error
	SchedulePipeline(name string, schedule Schedule) error
	GetPipelineInstance(pipeline PipelineObject) (map[string]interface{}, error)
	CommentOnPipeline(comment PipelineObject) error
	GetPipelineConfig(name string) (PipelineConfig, error)
	UpdatePipelineConfig(config PipelineConfig) (PipelineConfig, error)
	CreatePipeline(config PipelineConfig) (PipelineConfig, error)
	DeletePipeline(name string) error
	ExportPipelineToConfigRepoFormat(pipelineName, pluginID string) (PipelineExport, error)
	GetScheduledJobs() (ScheduledJobs, error)
	RunFailedJobs(stage Stage) (string, error)
	RunJobs(stage Stage) (string, error)
	RunStage(stage Stage) (string, error)
	CancelStage(stage Stage) (string, error)
	ExtractTemplatePipeline(pipeline, template string) (PipelineConfig, error)
	EncryptText(value string) (Encrypted, error)
	DecryptText(value, cipherKey string) (string, error)
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
	GetElasticAgentProfileUsage(profileID string) ([]ElasticProfileUsage, error)
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
	GetMaterialUsage(materialID string) ([]string, error)
	NotifyMaterial(material Material) (string, error)
	MaterialTriggerUpdate(id string) (map[string]string, error)
	GetRoles() (RolesConfig, error)
	GetRole(name string) (Role, error)
	GetRolesByType(roleType string) (RolesConfig, error)
	CreateRole(config Role) (Role, error)
	UpdateRole(config Role) (Role, error)
	DeleteRole(name string) error
	GetPluginsInfo() (PluginsInfo, error)
	GetPluginInfo(name string) (Plugin, error)
	GetUsers() ([]User, error)
	GetUser(user string) (User, error)
	GetCurrentUser() (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	UpdateCurrentUser(user User) (User, error)
	DeleteUser(user string) error
	BulkDeleteUsers(users map[string]interface{}) error
	BulkEnableDisableUsers(users map[string]interface{}) error
	GetPipelineVSM(pipeline, instance string) (VSM, error)
	GetPermissions(query map[string]string) (Permission, error)
	GetCCTray() ([]Project, error)
	SetRetryCount(count int)
	SetRetryWaitTime(count int)
}

// Auth holds information of authorisations configurations used for GoCd.
type Auth struct {
	UserName    string `json:"user_name,omitempty" yaml:"user_name,omitempty"`
	Password    string `json:"password,omitempty" yaml:"password,omitempty"`
	BearerToken string `json:"bearer_token,omitempty" yaml:"bearer_token,omitempty"`
	NoAuth      bool   `json:"no_auth,omitempty" yaml:"no_auth,omitempty"`
}

// NewClient returns new instance of httpClient when invoked.
func NewClient(baseURL string, auth Auth, logLevel string, caContent []byte) GoCd {
	logger := log.New()
	logger.SetLevel(goCdLogger.GetLoglevel(logLevel))
	logger.WithField(goCdAPILoggerName, true)
	logger.SetFormatter(&log.JSONFormatter{})

	newClient := resty.New()
	newClient.SetLogger(logger)
	newClient.SetRetryCount(defaultRetryCount)
	newClient.SetRetryWaitTime(defaultRetryWaitTime * time.Second)

	if logLevel == "debug" {
		newClient.SetDebug(true)
	}

	auth.setAuth(newClient)

	// setting authorization
	newClient.SetBaseURL(baseURL)

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

func (auth *Auth) setAuth(newClient *resty.Client) {
	if auth.NoAuth {
		// skip do nothing
	} else if len(auth.BearerToken) != 0 {
		newClient.SetAuthToken(auth.BearerToken)
	} else {
		newClient.SetBasicAuth(auth.UserName, auth.Password)
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

func GetGoCDMethodNames() []string {
	t := reflect.TypeOf((*GoCd)(nil)).Elem()

	var methodNames []string

	for i := 0; i < t.NumMethod(); i++ {
		methodNames = append(methodNames, t.Method(i).Name)
	}

	sort.Strings(methodNames)

	return methodNames
}

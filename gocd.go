package gocd

import (
	"crypto/tls"
	"crypto/x509"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

//go:generate $PWD/scripts/mocks $PWD/internal/mocks

// client holds resty.Client which could be used for interacting with GoCD and other information.
type client struct {
	httpClient *resty.Client
	logger     *log.Logger
}

// GoCd implements methods to get various information from GoCD.
type GoCd interface {
	GetAgentsInfo() ([]Agent, error)
	GetAgentJobRunHistory(agents []string) ([]AgentJobHistory, error)
	GetHealthInfo() ([]ServerHealth, error)
	GetConfigRepoInfo() ([]ConfigRepo, error)
	GetAdminsInfo() (SystemAdmins, error)
	GetPipelineGroupInfo() ([]PipelineGroup, error)
	GetEnvironmentInfo() ([]Environment, error)
	GetVersionInfo() (VersionInfo, error)
	GetBackupInfo() (BackupConfig, error)
	GetPipelines() (PipelinesInfo, error)
	GetPipelineState(pipelines []string) ([]PipelineState, error)
}

// NewClient returns new instance of httpClient when invoked.
func NewClient(baseURL, userName, passWord, loglevel string,
	caContent []byte, logLevel string,
) GoCd {
	newClient := resty.New()
	newClient.SetRetryCount(defaultRetryCount)
	newClient.SetRetryWaitTime(defaultRetryWaitTime * time.Second)
	if loglevel == "debug" {
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

	logger := log.New()
	logger.SetLevel(getLoglevel(logLevel))
	logger.WithField(goCdAPILoggerName, true)
	logger.SetFormatter(&log.JSONFormatter{})

	return &client{
		httpClient: newClient,
		logger:     logger,
	}
}

func getLoglevel(level string) log.Level {
	switch strings.ToLower(level) {
	case log.WarnLevel.String():
		return log.WarnLevel
	case log.DebugLevel.String():
		return log.DebugLevel
	case log.TraceLevel.String():
		return log.TraceLevel
	default:
		return log.InfoLevel
	}
}

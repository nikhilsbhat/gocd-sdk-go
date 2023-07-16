package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
	goCdLogger "github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var (
	yamlPluginURLTemplate   = "https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/%s/yaml-config-plugin-%s.jar"
	jsonPluginURLTemplate   = "https://github.com/tomzo/gocd-json-config-plugin/releases/download/%s/json-config-plugin-%s.jar"
	groovyPluginURLTemplate = "https://github.com/gocd-contrib/gocd-groovy-dsl-config-plugin/releases/download/v%s/gocd-groovy-dsl-config-plugin-%s.jar"
	githubAPIBaseURL        = "https://api.github.com/repos/%s"
	yamlPluginAPIURL        = fmt.Sprintf(githubAPIBaseURL, "tomzo/gocd-yaml-config-plugin/tags")
	jsonPluginAPIURL        = fmt.Sprintf(githubAPIBaseURL, "tomzo/gocd-json-config-plugin/tags")
	groovyPluginAPIURL      = fmt.Sprintf(githubAPIBaseURL, "gocd-contrib/gocd-groovy-dsl-config-plugin/tags")
)

type Plugin interface {
	ValidatePlugin(pipelines []string) (bool, error)
	Download() (string, error)
	SetType(pipelines []string) error
	GetType() string
	SetVersion(version string)
	GetVersion() string
}

type Config struct {
	Version      string `json:"version,omitempty" yaml:"version,omitempty" mapstructure:"version"`
	Path         string `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path"`
	URL          string `json:"url,omitempty" yaml:"url,omitempty" mapstructure:"url"`
	log          *log.Logger
	PipelineType string
}

type GithubTags struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

func (cfg *Config) ValidatePlugin(pipelines []string) (bool, error) {
	if missingPipelines, ok := cfg.exists(pipelines); !ok {
		return false, &errors.PipelineValidationError{
			Message: fmt.Sprintf("failed to validate pipelines, following pipelines are not found '%s'", strings.Join(missingPipelines, ",")),
		}
	}

	cmdArgs := append([]string{"-jar", cfg.Path, "syntax"}, pipelines...)

	cmd := exec.Command("java", cmdArgs...)

	cfg.log.Debugf("command that would be executed to validate syntax is '%s'", cmd.String())

	out, err := cmd.CombinedOutput()
	if err != nil {
		cfg.log.Debugf("invoking '%s' errored with non ok exit code: '%v'", cfg.Path, err)

		return false, &errors.PipelineValidationError{
			Message: fmt.Sprintf("validating pipeline failed with '%s'", string(out)),
		}
	}

	cfg.log.Debugf("validating pipeline against plugin returned '%s'", string(out))

	return true, nil
}

func (cfg *Config) exists(pipelines []string) ([]string, bool) {
	missingPipelines := make([]string, 0)

	for _, pipeline := range pipelines {
		if _, err := os.Stat(pipeline); os.IsNotExist(err) {
			cfg.log.Errorf("pipeline '%s' does not exits", pipeline)
			missingPipelines = append(missingPipelines, pipeline)
		}
	}

	if len(missingPipelines) != 0 {
		return missingPipelines, false
	}

	return nil, true
}

func (cfg *Config) SetType(pipelines []string) error {
	var fileType string

	for _, pipeline := range pipelines {
		if len(fileType) != 0 {
			if fileType != strings.TrimPrefix(filepath.Ext(pipeline), ".") {
				return &errors.PipelineValidationError{
					Message: "cannot club multiple pipeline file types for validation, should be one of yaml|json|groovy",
				}
			}
		}
		fileType = strings.TrimPrefix(filepath.Ext(pipeline), ".")
	}

	cfg.PipelineType = fileType

	return nil
}

func (cfg *Config) GetType() string {
	return cfg.PipelineType
}

func (cfg *Config) SetVersion(version string) {
	cfg.Version = version
}

func (cfg *Config) GetVersion() string {
	return cfg.Version
}

func (cfg *Config) GetLatestRelease(pluginURL string) (string, error) {
	var tags []GithubTags

	cfg.log.Debugf("fetching latest version information using '%s'", pluginURL)

	httpClient := resty.New()
	resp, err := httpClient.R().Get(pluginURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &tags); err != nil {
		return "", &errors.MarshalError{Err: err}
	}

	return tags[0].Name, nil
}

func (cfg *Config) setURL() error {
	if len(cfg.URL) != 0 {
		return nil
	}

	cfg.log.Debugf("plugin download url is not passed, setting it to default (github release) value")

	switch cfg.PipelineType {
	case "yaml":
		if len(cfg.Version) == 0 {
			version, err := cfg.GetLatestRelease(yamlPluginAPIURL)
			if err != nil {
				return err
			}
			cfg.Version = version
		}
		cfg.URL = fmt.Sprintf(yamlPluginURLTemplate, cfg.Version, cfg.Version)
	case "json":
		if len(cfg.Version) == 0 {
			version, err := cfg.GetLatestRelease(jsonPluginAPIURL)
			if err != nil {
				return err
			}
			cfg.Version = version
		}
		cfg.URL = fmt.Sprintf(jsonPluginURLTemplate, cfg.Version, cfg.Version)
	case "groovy":
		if len(cfg.Version) == 0 {
			version, err := cfg.GetLatestRelease(groovyPluginAPIURL)
			if err != nil {
				return err
			}
			cfg.Version = version
		}
		cfg.URL = fmt.Sprintf(groovyPluginURLTemplate, cfg.Version, cfg.Version)
	default:
		return &errors.PipelineValidationError{
			Message: fmt.Sprintf("unknown filetype '%s', supported are yaml|json|groovy", cfg.PipelineType),
		}
	}

	return nil
}

func (cfg *Config) Download() (string, error) {
	if len(cfg.Path) != 0 {
		cfg.log.Debugf("local path to plugin is set to '%s', skipping downloading plugin", cfg.Path)

		return cfg.Path, nil
	}

	if err := cfg.setURL(); err != nil {
		return "", err
	}

	cfg.log.Debugf("plugin download url is set to '%s'", cfg.URL)

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return "", err
	}

	pluginName := path.Base(parsedURL.Path)
	pluginLocalPath := filepath.Join(filepath.Join(home, ".gocd", "plugins"), pluginName)

	if _, err = os.Stat(pluginLocalPath); err == nil {
		cfg.log.Debugf("plugin jar already present under '%s', skipping plugin download", pluginLocalPath)
		cfg.Path = pluginLocalPath

		return pluginLocalPath, nil
	}

	cfg.log.Debugf("downloading plugin under '%s'", pluginLocalPath)

	httpClient := resty.New()

	resp, err := httpClient.R().
		SetOutput(pluginLocalPath).
		Get(cfg.URL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", &errors.PipelineValidationError{
			Message: fmt.Sprintf("downloading plugin returned non OK response code '%d' with BODY: '%s'", resp.StatusCode(), resp.Body()),
		}
	}

	cfg.log.Debugf("plugin '%s' downloaded successfully under '%s'", pluginName, pluginLocalPath)

	cfg.Path = pluginLocalPath

	return pluginLocalPath, nil
}

func NewPluginConfig(version, path, url, loglevel string) Plugin {
	logger := log.New()
	logger.SetLevel(goCdLogger.GetLoglevel(loglevel))
	logger.WithField("pipeline-validator", true)
	logger.SetFormatter(&log.JSONFormatter{})

	return &Config{
		log:     logger,
		Version: version,
		Path:    path,
		URL:     url,
	}
}

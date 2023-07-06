package plugin_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type YAMLPluginTestSuite struct {
	suite.Suite
	config   plugin.Plugin
	homePath string
}

func (suite *YAMLPluginTestSuite) SetupTest() {
	cfg := plugin.NewPluginConfig("0.13.0", "", "", "debug")

	homePath, err := os.UserHomeDir()
	suite.NoError(err)

	suite.config = cfg
	suite.homePath = homePath
}

func (suite *YAMLPluginTestSuite) TearDownSuite() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.yaml")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	pluginPath, err := suite.config.Download()
	suite.NoError(err)

	err = os.RemoveAll(pluginPath)
	suite.NoError(err)
}

func (suite *YAMLPluginTestSuite) TestValidatePluginTests_ShouldFailValidatingPipelineDueToMissMatchOfMultiplePipelineTypes() {
	yamlPipelinePath := filepath.Join(suite.homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.yaml")
	jsonPipelinePath := filepath.Join(suite.homePath, "gocd-sdk-go/internal/fixtures/agent_run_history.json")

	pipelines := []string{yamlPipelinePath, jsonPipelinePath}
	err := suite.config.SetType(pipelines)
	suite.EqualError(err, "cannot club multiple pipeline file types for validation, should be one of yaml|json|groovy")
}

func (suite *YAMLPluginTestSuite) TestValidatePluginTests_ShouldFailValidatingYAMLPipelineDueToWrongPath() {
	pipelinePath := filepath.Join(suite.homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.yaml")
	expectedError := fmt.Sprintf("failed to validate pipelines, following pipelines are not found '%s'", pipelinePath)

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.EqualError(err, expectedError)
	suite.Equal(false, actual)
}

func (suite *YAMLPluginTestSuite) TestValidatePluginTests_ShouldSuccessfullyValidateTheYAMLPipelineFile() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.yaml")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.NoError(err)
	suite.Equal(true, actual)
}

func (suite *YAMLPluginTestSuite) TestValidatePluginTests_ValidationShouldFailDueToErrorsInYAMLPipelineFile() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline-defect.gocd.yaml")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.Error(err)
	suite.Equal(false, actual)
}

func TestPipelineValidateTestSuite(t *testing.T) {
	suite.Run(t, new(YAMLPluginTestSuite))
}

type JOSNPluginTestSuite struct {
	suite.Suite
	config   plugin.Plugin
	homePath string
}

func (suite *JOSNPluginTestSuite) SetupTest() {
	cfg := plugin.NewPluginConfig("0.6.0", "", "", "debug")

	homePath, err := os.UserHomeDir()
	suite.NoError(err)

	suite.config = cfg
	suite.homePath = homePath
}

func (suite *JOSNPluginTestSuite) TearDownSuite() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.json")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	pluginPath, err := suite.config.Download()
	suite.NoError(err)

	err = os.RemoveAll(pluginPath)
	suite.NoError(err)
}

func (suite *JOSNPluginTestSuite) TestValidatePluginTests_ShouldFailValidatingJSONPipelineDueToWrongPath() {
	pipelinePath := filepath.Join(suite.homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.json")
	expectedError := fmt.Sprintf("failed to validate pipelines, following pipelines are not found '%s'", pipelinePath)

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.EqualError(err, expectedError)
	suite.Equal(false, actual)
}

func (suite *JOSNPluginTestSuite) TestValidatePluginTests_ShouldSuccessfullyValidateTheJSONPipelineFile() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.json")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.NoError(err)
	suite.Equal(true, actual)
}

func (suite *JOSNPluginTestSuite) TestValidatePluginTests_ValidationShouldFailDueToErrorsInJSONPipelineFile() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline-defect.gocd.json")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.Error(err)
	suite.Equal(false, actual)
}

func TestJSONPipelineValidateTestSuite(t *testing.T) {
	suite.Run(t, new(JOSNPluginTestSuite))
}

type GroovyPluginTestSuite struct {
	suite.Suite
	config   plugin.Plugin
	homePath string
}

func (suite *GroovyPluginTestSuite) SetupTest() {
	cfg := plugin.NewPluginConfig("2.1.3-512", "", "", "debug")

	homePath, err := os.UserHomeDir()
	suite.NoError(err)

	suite.config = cfg
	suite.homePath = homePath
}

func (suite *GroovyPluginTestSuite) TearDownSuite() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.groovy")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	pluginPath, err := suite.config.Download()
	suite.NoError(err)

	err = os.RemoveAll(pluginPath)
	suite.NoError(err)
}

func (suite *GroovyPluginTestSuite) TestValidatePluginTests_ShouldFailValidatingGroovyPipelineDueToWrongPath() {
	pipelinePath := filepath.Join(suite.homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.groovy")
	expectedError := fmt.Sprintf("failed to validate pipelines, following pipelines are not found '%s'", pipelinePath)

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.EqualError(err, expectedError)
	suite.Equal(false, actual)
}

func (suite *GroovyPluginTestSuite) TestValidatePluginTests_ShouldSuccessfullyValidateTheGroovyPipelineFile() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.groovy")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.NoError(err)
	suite.Equal(true, actual)
}

func (suite *GroovyPluginTestSuite) TestValidatePluginTests_ValidationShouldFailDueToErrorsInGroovyPipelineFile() {
	pipelinePath := filepath.Join(suite.homePath, "my-opensource/gocd-sdk-go/internal/fixtures/sample-pipeline-defect.gocd.groovy")

	err := suite.config.SetType([]string{pipelinePath})
	suite.NoError(err)

	_, err = suite.config.Download()
	suite.NoError(err)

	actual, err := suite.config.ValidatePlugin([]string{pipelinePath})
	suite.Error(err)
	suite.Equal(false, actual)
}

func TestGroovyPipelineValidateTestSuite(t *testing.T) {
	suite.Run(t, new(GroovyPluginTestSuite))
}

func TestConfig_Download(t *testing.T) {
	t.Run("should be able to download the plugin successfully", func(t *testing.T) {
		cfg := plugin.NewPluginConfig("0.13.0", "", "", "debug")

		homePath, err := os.UserHomeDir()
		assert.NoError(t, err)

		pipelinePath := filepath.Join(homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.yaml")

		err = cfg.SetType([]string{pipelinePath})
		assert.NoError(t, err)

		expectedPluginPath := filepath.Join(homePath, ".gocd/plugins/yaml-config-plugin-0.13.0.jar")

		pluginPath, err := cfg.Download()
		assert.NoError(t, err)
		assert.Equal(t, expectedPluginPath, pluginPath)
	})

	t.Run("should error out due to unsupported plugin type", func(t *testing.T) {
		cfg := plugin.NewPluginConfig("0.13.0", "", "", "debug")

		homePath, err := os.UserHomeDir()
		assert.NoError(t, err)

		pipelinePath := filepath.Join(homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.toml")

		err = cfg.SetType([]string{pipelinePath})
		assert.NoError(t, err)

		pluginPath, err := cfg.Download()
		assert.EqualError(t, err, "unknown filetype 'toml', supported are yaml|json|groovy")
		assert.Equal(t, "", pluginPath)
	})

	t.Run("should error out due to wrong url set", func(t *testing.T) {
		cfg := plugin.NewPluginConfig("0.13.0", "", "://github.com/gocd-contrib/gocd-groovy-dsl-config-plugin", "debug")

		homePath, err := os.UserHomeDir()
		assert.NoError(t, err)

		pipelinePath := filepath.Join(homePath, "gocd-sdk-go/internal/fixtures/sample-pipeline.gocd.yaml")

		err = cfg.SetType([]string{pipelinePath})
		assert.NoError(t, err)

		pluginPath, err := cfg.Download()
		assert.EqualError(t, err, "parse \"://github.com/gocd-contrib/gocd-groovy-dsl-config-plugin\": missing protocol scheme")
		assert.Equal(t, "", pluginPath)
	})
}

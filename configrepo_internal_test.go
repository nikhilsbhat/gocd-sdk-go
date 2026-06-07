package gocd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_client_GetPipelineFiles_ShouldReturnAbsErrorForExplicitPipeline(t *testing.T) {
	expectedErr := errors.New("absolute path failed")
	originalFilepathAbs := filepathAbs
	filepathAbs = func(string) (string, error) {
		return "", expectedErr
	}
	t.Cleanup(func() {
		filepathAbs = originalFilepathAbs
	})

	client := NewClient("http://localhost:8156/go", Auth{NoAuth: true}, "debug", nil).(*client)

	actual, err := client.GetPipelineFiles("", []string{"internal/fixtures/mail_server_config.json"})
	require.ErrorIs(t, err, expectedErr)
	assert.Nil(t, actual)
}

func Test_client_GetPipelineFiles_ShouldSkipMatchedFileWhenAbsFailsDuringWalk(t *testing.T) {
	expectedErr := errors.New("absolute path failed")
	originalFilepathAbs := filepathAbs
	filepathAbs = func(path string) (string, error) {
		return path, expectedErr
	}
	t.Cleanup(func() {
		filepathAbs = originalFilepathAbs
	})

	client := NewClient("http://localhost:8156/go", Auth{NoAuth: true}, "debug", nil).(*client)

	actual, err := client.GetPipelineFiles("internal/fixtures", nil, "mail_server_config.json")
	require.NoError(t, err)
	assert.Equal(t, []PipelineFiles{
		{
			Name: "mail_server_config.json",
			Path: "internal/fixtures/mail_server_config.json",
		},
	}, actual)
}

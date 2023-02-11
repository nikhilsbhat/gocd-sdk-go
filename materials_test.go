package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/materials.json
var attributesJSON string

func Test_client_GetMaterials(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}
	t.Run("should be able to fetch all available materials present in GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(attributesJSON), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.Material{
			{
				Type:        "svn",
				Fingerprint: "6b0cd6b9181434866c555f3c3bb780e950389a456764c876d804c848efbad554",
				Attributes: gocd.Attribute{
					URL:                 "https://github.com/gocd/gocd",
					Username:            "admin",
					EncryptedPassword:   "AES:lTFKwi5NvrlqmQ3LOQ5UQA==:Ebggz5N27w54NrhSXKIbng==",
					Destination:         "test",
					AutoUpdate:          false,
					CheckExternals:      false,
					UseTickets:          false,
					IgnoreForScheduling: false,
					InvertFilter:        false,
				},
			},
			{
				Type:        "hg",
				Fingerprint: "f6b61bc6b33e524c2a94c5be4a4661e5f1d2b74ca089418de20de10b282e39e9",
				Attributes: gocd.Attribute{
					URL:                 "ssh://hg@bitbucket.org/example/gocd-hg",
					Destination:         "test",
					AutoUpdate:          false,
					CheckExternals:      false,
					UseTickets:          false,
					IgnoreForScheduling: false,
					InvertFilter:        false,
				},
			},
			{
				Type:        "dependency",
				Fingerprint: "21f7963a8eca6cced5b4f20b57b4e9cb8edaed29d683c83621a77dfb499c553d",
				Attributes: gocd.Attribute{
					Pipeline:            "up42",
					Stage:               "up42_stage",
					Name:                "up42_material",
					AutoUpdate:          true,
					IgnoreForScheduling: false,
				},
			},
		}

		actual, err := client.GetMaterials()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(attributesJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(attributesJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("attributesJSON"), http.StatusOK, correctArtifactHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "call made to get all available materials errored with: "+
			"Get \"http://localhost:8156/go/api/config/materials\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/encryption.json
var encryptionJSON string

func Test_client_EncryptText(t *testing.T) {
	correctEncryptionHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}
	t.Run("should be able to encrypt the value passed successfully", func(t *testing.T) {
		server := mockServer([]byte(encryptionJSON), http.StatusOK, correctEncryptionHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.Encrypted{EncryptedValue: "aSdiFgRRZ6A="}

		actual, err := client.EncryptText("value_to_encrypt")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while encrypting a value as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("encryptionJSON"), http.StatusBadGateway, correctEncryptionHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.EncryptText("value_to_encrypt")
		assert.EqualError(t, err, "body: encryptionJSON httpcode: 502")
		assert.Equal(t, gocd.Encrypted{}, actual)
	})

	t.Run("should error out while encrypting a value server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"encrypting a value"}`), http.StatusOK, correctEncryptionHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.EncryptText("value_to_encrypt")
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.Encrypted{}, actual)
	})

	t.Run("should error out while encrypting a value as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.EncryptText("value_to_encrypt")
		assert.EqualError(t, err, "call made to encrypt a value errored with: "+
			"Post \"http://localhost:8156/go/api/admin/encrypt\": dial tcp 127.0.0.1:8156: connect: connection refused")
		assert.Equal(t, gocd.Encrypted{}, actual)
	})
}

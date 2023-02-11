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
			"Post \"http://localhost:8156/go/api/admin/encrypt\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Encrypted{}, actual)
	})
}

func Test_client_DecryptText(t *testing.T) {
	cipher := "ab533bc2b64169f487412301afa6f5f6"
	t.Run("should be able to decrypt the secret successfully", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		response, err := client.DecryptText("AES:wSOqnltxM6Rp9j0Tb8uWpw==:4zVLtLx9msGleK+pLOOUHg==", cipher)
		assert.NoError(t, err)
		assert.Equal(t, "badger", response)
	})

	t.Run("should error out while decrypting secret due to wrong cipher passed", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		response, err := client.DecryptText("AES:wSOqnltxM6Rp9j0Tb8uWpw==:4zVLtLx9msGleK+pLOOUHg==", "kencehcf84nnkcxjrfjx48")
		assert.EqualError(t, err, "encoding/hex: invalid byte: U+006B 'k'")
		assert.Equal(t, "", response)
	})

	t.Run("should error out while decrypting secret due to malformed encrypted value", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		response, err := client.DecryptText("AES:wSOqnltxM6Rp9j0Tb8uWpw==:hjdsdjxwerj474x3+pLOOUHg==", "kencehcf84nnkcxjrfjx48")
		assert.EqualError(t, err, "illegal base64 data at input byte 24")
		assert.Equal(t, "", response)
	})

	t.Run("should error out while decrypting secret due to malformed encoded IV", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		response, err := client.DecryptText("AES:wefxe343348xnwh43x4ux==:4zVLtLx9msGleK+pLOOUHg==", "kencehcf84nnkcxjrfjx48")
		assert.EqualError(t, err, "illegal base64 data at input byte 21")
		assert.Equal(t, "", response)
	})

	t.Run("should error out while decrypting secret as no secret or cipher is passed", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		response, err := client.DecryptText("", "")
		assert.EqualError(t, err, "value or cipher key cannot be empty")
		assert.Equal(t, "", response)
	})

	t.Run("should be able to decrypt the secret successfully", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		response, err := client.DecryptText("AES:wSOqnltxM6Rp9j0Tb8uWpw==:4zVLtLx9msGleK+pLOOUHg==", "cb533bc2b64169f487412301afa6f5f")
		assert.EqualError(t, err, "encoding/hex: odd length hex string")
		assert.Equal(t, "", response)
	})
}

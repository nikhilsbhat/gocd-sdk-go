package gocd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/copier"
)

var errValueCipher = errors.New("value or cipher key cannot be empty")

func (conf *client) EncryptText(value string) (Encrypted, error) {
	var encryptedValue Encrypted
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return encryptedValue, err
	}

	valueObj := map[string]string{"value": value}
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(valueObj).
		Post(EncryptEndpoint)
	if err != nil {
		return encryptedValue, fmt.Errorf("call made to encrypt a value errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return encryptedValue, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &encryptedValue); err != nil {
		return encryptedValue, ResponseReadError(err.Error())
	}

	return encryptedValue, nil
}

func (conf *client) DecryptText(value, cipherKey string) (string, error) {
	var decryptedValue string
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return decryptedValue, err
	}

	if len(value) == 0 || len(cipherKey) == 0 {
		return "", errValueCipher
	}
	// AES encrypted value should be split to get encoded IV and data from it.
	// Sample AES encrypted value: 'AES:wSOqnltxM6Rp9j0Tb8uWpw==:4zVLtLx9msGleK+pLOOUHg=='. Upon splitting we would get three elements
	// By ignoring first 'AES' we are left with two elements,
	// first 'wSOqnltxM6Rp9j0Tb8uWpw==' would be encodedIV and the second/last element '4zVLtLx9msGleK+pLOOUHg==' would be encoded data.
	// Both IV and data are base64 encoded, to decrypt we should base64 decode it first.
	dataSplit := strings.Split(value, ":")

	decodedIV, err := base64.StdEncoding.DecodeString(dataSplit[1])
	if err != nil {
		return "", err
	}
	decodedData, err := base64.StdEncoding.DecodeString(dataSplit[2])
	if err != nil {
		return "", err
	}

	// cipher block to be obtained from the cipher key.
	block, err := cipherBlock(cipherKey)
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, decodedIV)
	decrypted := make([]byte, len(decodedData))
	ecb.CryptBlocks(decrypted, decodedData)
	pk5trimmedValue := string(pkcs5Trimming(decrypted))

	return pk5trimmedValue, nil
}

func cipherBlock(key string) (cipher.Block, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]

	return encrypt[:len(encrypt)-int(padding)]
}

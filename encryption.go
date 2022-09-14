package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

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

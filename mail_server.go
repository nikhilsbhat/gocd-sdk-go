package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

	"github.com/jinzhu/copier"
)

func (conf *client) GetMailServerConfig() (MailServerConfig, error) {
	var mailConfig MailServerConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return mailConfig, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(MailServerConfigEndpoint)
	if err != nil {
		return mailConfig, &errors.APIError{Err: err, Message: "get mail server config"}
	}

	if resp.StatusCode() != http.StatusOK {
		return mailConfig, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &mailConfig); err != nil {
		return mailConfig, &errors.MarshalError{Err: err}
	}

	return mailConfig, nil
}

func (conf *client) CreateOrUpdateMailServerConfig(mailCfg MailServerConfig) (MailServerConfig, error) {
	var mailConfig MailServerConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return mailConfig, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(mailCfg).
		Get(MailServerConfigEndpoint)
	if err != nil {
		return mailConfig, &errors.APIError{Err: err, Message: "create or update mail server config"}
	}

	if resp.StatusCode() != http.StatusOK {
		return mailConfig, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &mailConfig); err != nil {
		return mailConfig, &errors.MarshalError{Err: err}
	}

	return mailConfig, nil
}

func (conf *client) DeleteMailServerConfig() error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Delete(MailServerConfigEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: "delete mail server config"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

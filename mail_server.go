package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		return mailConfig, fmt.Errorf("call made to get mail server config errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return mailConfig, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &mailConfig); err != nil {
		return mailConfig, ResponseReadError(err.Error())
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
		return mailConfig, fmt.Errorf("call made to create or update mail server config errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return mailConfig, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &mailConfig); err != nil {
		return mailConfig, ResponseReadError(err.Error())
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
		return fmt.Errorf("call made to delete mail server config errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

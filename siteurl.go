package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetSiteURL fetches the site url config configured from GoCD.
func (conf *client) GetSiteURL() (SiteURLConfig, error) {
	var site SiteURLConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return site, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(SiteURLEndpoint)
	if err != nil {
		return site, fmt.Errorf("call made to get site url errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return site, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &site); err != nil {
		return site, err
	}

	return site, nil
}

// CreateOrUpdateSiteURL creates/updates the site url configured in GoCD.
func (conf *client) CreateOrUpdateSiteURL(config SiteURLConfig) (SiteURLConfig, error) {
	var site SiteURLConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return site, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(SiteURLEndpoint)
	if err != nil {
		return site, fmt.Errorf("call made to create/update site url errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return site, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &site); err != nil {
		return site, err
	}

	return site, nil
}
